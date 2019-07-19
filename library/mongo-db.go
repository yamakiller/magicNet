package library

import (
	"context"
	"fmt"
	"magicNet/engine/util"
	"sync"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

type mongoClient struct {
	c      *mongo.Client
	db     *mongo.Database
	ctx    context.Context
	cancel context.CancelFunc
}

// Connect : 连接服务
func (mgd *mongoClient) connect(host []string,
	uri string,
	dbName string,
	poolSize uint16,
	sckTimeSec int,
	timeSec int,
	hbSec int,
	idleSec int) error {
	opt := options.ClientOptions{}
	mgd.ctx, mgd.cancel = context.WithTimeout(context.Background(), time.Duration(timeSec)*time.Second)
	client, err := mongo.Connect(mgd.ctx,
		opt.ApplyURI(uri),
		opt.SetHosts(host),
		opt.SetHeartbeatInterval(time.Duration(hbSec)*time.Second),
		opt.SetMaxPoolSize(poolSize),
		opt.SetMaxConnIdleTime(time.Duration(idleSec)*time.Second),
		opt.SetSocketTimeout(time.Duration(sckTimeSec)*time.Second))
	if err != nil {
		mgd.cancel()
		mgd.cancel = nil
		return err
	}

	mgd.db = client.Database(dbName)
	if mgd.db == nil {
		client.Disconnect(mgd.ctx)
		mgd.cancel()
		mgd.cancel = nil
		return fmt.Errorf("mongoDB Database %s does not exist", dbName)
	}

	mgd.c = client

	return nil
}

// Close : 关闭
func (mgd *mongoClient) close() {
	defer mgd.cancel()
	mgd.c.Disconnect(mgd.ctx)
	mgd.c = nil
}

// MongoDB :
type MongoDB struct {
	cs   []*mongoClient
	size int
	mx   sync.Mutex

	URI           string
	HOSTS         []string
	DBName        string
	PoolSize      uint16
	SocketTimeSec int
	TimeSec       int
	HbtSec        int
	IdleSec       int
	MinClient     int
	MaxClient     int
}

// Init : x
func (mgb *MongoDB) Init() error {
	mgb.mx.Lock()
	defer mgb.mx.Unlock()
	for i := 0; i < mgb.MinClient; i++ {
		mgc := &mongoClient{}
		err := mgc.connect(mgb.HOSTS,
			mgb.URI,
			mgb.DBName,
			mgb.PoolSize,
			mgb.SocketTimeSec,
			mgb.TimeSec,
			mgb.HbtSec,
			mgb.IdleSec)
		util.AssertEmpty(err != nil, err.Error())
		mgb.cs = append(mgb.cs, mgc)
		mgb.size++
	}

	return nil
}

func (mgb *MongoDB) freeClient(free *mongoClient) {
	mgb.mx.Lock()
	defer mgb.mx.Unlock()
	mgb.cs = append(mgb.cs, free)
}

func (mgb *MongoDB) getClient() (*mongoClient, error) {
	mgb.mx.Lock()
	defer mgb.mx.Unlock()
	if len(mgb.cs) == 0 {
		if mgb.size < mgb.MaxClient {
			mgc := &mongoClient{}
			err := mgc.connect(mgb.HOSTS,
				mgb.URI,
				mgb.DBName,
				mgb.PoolSize,
				mgb.SocketTimeSec,
				mgb.TimeSec,
				mgb.HbtSec,
				mgb.IdleSec)
			if err != nil {
				return nil, err
			}
			mgb.size++
			return mgc, nil
		}
		return nil, fmt.Errorf("mongoDB dbpooling is fulled")
	}

	client := mgb.cs[0]
	mgb.cs = mgb.cs[1:]

	if err := client.c.Ping(client.ctx, readpref.Primary()); err != nil {
		mgb.size--
		client.close()
		return nil, err
	}

	return client, nil
}

// InsertOne : 插入一条数据
func (mgb *MongoDB) InsertOne(name string, document interface{}) (interface{}, error) {
	client, err := mgb.getClient()
	if err != nil {
		return nil, err
	}
	defer mgb.freeClient(client)

	r, rerr := client.db.Collection(name).InsertOne(client.ctx, document)
	if rerr != nil {
		return nil, rerr
	}

	return r.InsertedID, nil
}

// InsertMany : 插入多条数据
func (mgb *MongoDB) InsertMany(name string, document []interface{}) ([]interface{}, error) {
	client, err := mgb.getClient()
	if err != nil {
		return nil, err
	}
	defer mgb.freeClient(client)

	r, rerr := client.db.Collection(name).InsertMany(client.ctx, document)
	if rerr != nil {
		return nil, rerr
	}
	return r.InsertedIDs, nil
}

// FindOne : 查询一条数据
func (mgb *MongoDB) FindOne(name string, filter interface{}, out interface{}) error {
	client, err := mgb.getClient()
	if err != nil {
		return err
	}
	defer mgb.freeClient(client)

	r := client.db.Collection(name).FindOne(client.ctx, filter)
	if derr := r.Decode(out); derr != nil {
		return derr
	}

	return nil
}

// Find : 查询多条数据
func (mgb *MongoDB) Find(name string, filter interface{}, decode interface{}) ([]interface{}, error) {
	client, err := mgb.getClient()
	if err != nil {
		return nil, err
	}
	defer mgb.freeClient(client)

	r, rerr := client.db.Collection(name).Find(client.ctx, filter)
	if rerr != nil {
		return nil, rerr
	}

	defer r.Close(client.ctx)
	ary := make([]interface{}, 0, 4)
	for r.Next(client.ctx) {
		if derr := r.Decode(&decode); derr != nil {
			return nil, derr
		}

		ary = append(ary, decode)
	}

	return ary, nil
}

//UpdateOne : xx
func (mgb *MongoDB) UpdateOne(name string, filter interface{}, update interface{}) (int64, int64, int64, interface{}, error) {
	client, err := mgb.getClient()
	if err != nil {
		return 0, 0, 0, nil, err
	}
	defer mgb.freeClient(client)
	r, rerr := client.db.Collection(name).UpdateOne(client.ctx, filter, update)
	if rerr != nil {
		return 0, 0, 0, nil, rerr
	}

	return r.MatchedCount, r.ModifiedCount, r.UpsertedCount, r.UpsertedID, nil
}

//UpdateMany : x
func (mgb *MongoDB) UpdateMany(name string, filter interface{}, update interface{}) (int64, int64, int64, interface{}, error) {
	client, err := mgb.getClient()
	if err != nil {
		return 0, 0, 0, nil, err
	}
	defer mgb.freeClient(client)
	r, rerr := client.db.Collection(name).UpdateMany(client.ctx, filter, update)
	if rerr != nil {
		return 0, 0, 0, nil, rerr
	}

	return r.MatchedCount, r.ModifiedCount, r.UpsertedCount, r.UpsertedID, nil
}

// ReplaceOne : 替换一条数据
func (mgb *MongoDB) ReplaceOne(name string, filter interface{}, replacement interface{}) (int64, int64, int64, interface{}, error) {
	client, err := mgb.getClient()
	if err != nil {
		return 0, 0, 0, nil, err
	}
	defer mgb.freeClient(client)

	r, rerr := client.db.Collection(name).ReplaceOne(client.ctx, filter, replacement)

	if rerr != nil {
		return 0, 0, 0, nil, rerr
	}

	return r.MatchedCount, r.ModifiedCount, r.UpsertedCount, r.UpsertedID, nil
}

//DeleteOne : 删除一条数据
func (mgb *MongoDB) DeleteOne(name string, filter interface{}) (int64, error) {
	client, err := mgb.getClient()
	if err != nil {
		return 0, err
	}
	defer mgb.freeClient(client)

	r, rerr := client.db.Collection(name).DeleteOne(client.ctx, filter)
	if rerr != nil {
		return 0, rerr
	}

	return r.DeletedCount, nil
}

//DeleteMany : 删除多条数据
func (mgb *MongoDB) DeleteMany(name string, filter interface{}) (int64, error) {
	client, err := mgb.getClient()
	if err != nil {
		return 0, err
	}
	defer mgb.freeClient(client)

	r, rerr := client.db.Collection(name).DeleteMany(client.ctx, filter)
	if rerr != nil {
		return 0, rerr
	}

	return r.DeletedCount, nil
}

// FindOneAndDelete :
func (mgb *MongoDB) FindOneAndDelete(name string, filter interface{}, out interface{}) error {
	client, err := mgb.getClient()
	if err != nil {
		return err
	}
	defer mgb.freeClient(client)

	r := client.db.Collection(name).FindOneAndDelete(client.ctx, filter)

	if derr := r.Decode(out); derr != nil {
		return derr
	}

	return nil
}

// FindOneAndUpdate :
func (mgb *MongoDB) FindOneAndUpdate(name string, filter interface{}, update interface{}, out interface{}) error {
	client, err := mgb.getClient()
	if err != nil {
		return err
	}
	defer mgb.freeClient(client)

	r := client.db.Collection(name).FindOneAndUpdate(client.ctx, filter, update)
	if derr := r.Decode(out); derr != nil {
		return derr
	}

	return nil
}

// FindOneAndReplace :
func (mgb *MongoDB) FindOneAndReplace(name string, filter interface{}, replacement interface{}, out interface{}) error {
	client, err := mgb.getClient()
	if err != nil {
		return err
	}
	defer mgb.freeClient(client)

	r := client.db.Collection(name).FindOneAndReplace(client.ctx, filter, replacement)
	if derr := r.Decode(out); derr != nil {
		return derr
	}

	return nil
}

// Distinct : 在指定字段中，查找
func (mgb *MongoDB) Distinct(name string, fieldName string, filter interface{}) ([]interface{}, error) {
	client, err := mgb.getClient()
	if err != nil {
		return nil, err
	}
	defer mgb.freeClient(client)

	r, rerr := client.db.Collection(name).Distinct(client.ctx, fieldName, filter)
	if rerr != nil {
		return nil, rerr
	}

	return r, nil
}

// Drop : 删除集
func (mgb *MongoDB) Drop(name string) error {
	client, err := mgb.getClient()
	if err != nil {
		return err
	}

	return client.db.Collection(name).Drop(client.ctx)
}

// CountDocuments : 获取文档总数
func (mgb *MongoDB) CountDocuments(name string, filter interface{}) (int64, error) {
	client, err := mgb.getClient()
	if err != nil {
		return 0, err
	}
	defer mgb.freeClient(client)

	return client.db.Collection(name).CountDocuments(client.ctx, filter)
}
