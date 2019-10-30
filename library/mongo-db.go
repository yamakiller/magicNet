package library

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/yamakiller/magicNet/util"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

type mongoClient struct {
	c       *mongo.Client
	db      *mongo.Database
	timeSec int
}

// Connect : Connection mongo db service
func (slf *mongoClient) connect(host []string,
	uri string,
	dbName string,
	poolSize uint16,
	sckTimeSec int,
	timeSec int,
	hbSec int,
	idleSec int) error {
	slf.timeSec = timeSec
	opt := options.ClientOptions{}
	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(timeSec)*time.Second)
	defer cancel()
	client, err := mongo.Connect(ctx,
		opt.ApplyURI(uri),
		opt.SetHosts(host),
		opt.SetHeartbeatInterval(time.Duration(hbSec)*time.Second),
		opt.SetMaxPoolSize(uint64(poolSize)),
		opt.SetMaxConnIdleTime(time.Duration(idleSec)*time.Second),
		opt.SetSocketTimeout(time.Duration(sckTimeSec)*time.Second))
	if err != nil {
		return err
	}

	slf.db = client.Database(dbName)
	if slf.db == nil {
		client.Disconnect(ctx)
		return fmt.Errorf("mongoDB Database %s does not exist", dbName)
	}

	slf.c = client

	return nil
}

func (slf *mongoClient) close() {
	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(slf.timeSec)*time.Second)
	defer cancel()
	slf.c.Disconnect(ctx)
	slf.c = nil
	slf.db = nil
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

// Close : 关闭
func (mgb *MongoDB) Close() {
	for {
		mgb.mx.Lock()
		if mgb.size == 0 {
			mgb.mx.Unlock()
			break
		}

		n := len(mgb.cs)
		if n == 0 {
			mgb.mx.Unlock()
			time.Sleep(time.Millisecond * 5)
			continue
		}

		for _, v := range mgb.cs {
			v.close()
		}

		mgb.cs = mgb.cs[n:]
		mgb.size -= n
	}
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

	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(mgb.TimeSec)*time.Second)
	defer cancel()

	if err := client.c.Ping(ctx, readpref.Primary()); err != nil {
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

	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(mgb.TimeSec)*time.Second)
	defer cancel()

	r, rerr := client.db.Collection(name).InsertOne(ctx, document)
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

	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(mgb.TimeSec)*time.Second)
	defer cancel()

	r, rerr := client.db.Collection(name).InsertMany(ctx, document)
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

	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(mgb.TimeSec)*time.Second)
	defer cancel()

	r := client.db.Collection(name).FindOne(ctx, filter)
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

	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(mgb.TimeSec)*time.Second)
	defer cancel()

	r, rerr := client.db.Collection(name).Find(ctx, filter)
	if rerr != nil {
		return nil, rerr
	}

	defer r.Close(ctx)
	ary := make([]interface{}, 0, 4)
	for r.Next(ctx) {
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

	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(mgb.TimeSec)*time.Second)
	defer cancel()

	r, rerr := client.db.Collection(name).UpdateOne(ctx, filter, update)
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

	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(mgb.TimeSec)*time.Second)
	defer cancel()

	r, rerr := client.db.Collection(name).UpdateMany(ctx, filter, update)
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

	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(mgb.TimeSec)*time.Second)
	defer cancel()

	r, rerr := client.db.Collection(name).ReplaceOne(ctx, filter, replacement)

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

	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(mgb.TimeSec)*time.Second)
	defer cancel()

	r, rerr := client.db.Collection(name).DeleteOne(ctx, filter)
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

	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(mgb.TimeSec)*time.Second)
	defer cancel()

	r, rerr := client.db.Collection(name).DeleteMany(ctx, filter)
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

	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(mgb.TimeSec)*time.Second)
	defer cancel()

	r := client.db.Collection(name).FindOneAndDelete(ctx, filter)

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

	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(mgb.TimeSec)*time.Second)
	defer cancel()

	r := client.db.Collection(name).FindOneAndUpdate(ctx, filter, update)
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

	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(mgb.TimeSec)*time.Second)
	defer cancel()

	r := client.db.Collection(name).FindOneAndReplace(ctx, filter, replacement)
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

	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(mgb.TimeSec)*time.Second)
	defer cancel()

	r, rerr := client.db.Collection(name).Distinct(ctx, fieldName, filter)
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

	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(mgb.TimeSec)*time.Second)
	defer cancel()

	return client.db.Collection(name).Drop(ctx)
}

// CountDocuments : 获取文档总数
func (mgb *MongoDB) CountDocuments(name string, filter interface{}) (int64, error) {
	client, err := mgb.getClient()
	if err != nil {
		return 0, err
	}
	defer mgb.freeClient(client)

	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(mgb.TimeSec)*time.Second)
	defer cancel()

	return client.db.Collection(name).CountDocuments(ctx, filter)
}
