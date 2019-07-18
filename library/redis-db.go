package library

import (
	"fmt"
	"magicNet/engine/util"
	"time"

	"github.com/gomodule/redigo/redis"
)

// RedisDB : xx
type RedisDB struct {
	c    *redis.Pool
	host string
	db   int
}

// Init : 初始化DB
func (rdb *RedisDB) Init(host string, db int, maxIdle int, maxActive int, idleSec int) error {
	rdb.host = host
	rdb.db = db

	rdb.c = &redis.Pool{
		MaxIdle:     maxIdle,
		MaxActive:   maxActive,
		IdleTimeout: time.Duration(idleSec) * time.Second,
		Dial: func() (redis.Conn, error) {
			c, err := redis.Dial("tcp", host)
			if err != nil {
				return nil, err
			}
			_, derr := c.Do("SELECT", db)
			util.Assert(derr == nil, fmt.Sprintln("redis select db error:", derr))
			return c, nil
		},
	}

	return nil
}

// Do : 执行Redis
func (rdb *RedisDB) Do(commandName string, args ...interface{}) (interface{}, error) {
	c := rdb.c.Get()
	defer c.Close()
	return c.Do(commandName, args)
}

// Close : 关闭
func (rdb *RedisDB) Close() {
	rdb.c.Close()
}
