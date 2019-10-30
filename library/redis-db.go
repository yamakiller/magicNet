package library

import (
	"fmt"
	"time"

	"github.com/yamakiller/magicNet/util"

	"github.com/gomodule/redigo/redis"
)

// RedisDB : xx
type RedisDB struct {
	c    *redis.Pool
	host string
	db   int
}

// Init : 初始化DB
func (slf *RedisDB) Init(host string, db int, maxIdle int, maxActive int, idleSec int) error {
	slf.host = host
	slf.db = db

	slf.c = &redis.Pool{
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
func (slf *RedisDB) Do(commandName string, args ...interface{}) (interface{}, error) {
	c := slf.c.Get()
	defer c.Close()
	return c.Do(commandName, args...)
}

// Close : 关闭
func (slf *RedisDB) Close() {
	slf.c.Close()
}
