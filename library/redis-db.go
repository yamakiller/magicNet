package library

import (
	"fmt"
	"time"

	"github.com/yamakiller/magicNet/util"

	"github.com/gomodule/redigo/redis"
)

//RedisDB desc
//@struct RedisDB desc: redis db opertioner
//@member (*redis.Pool) a redis connection pool
//@member (string) redis address
//@member (int)  redis db code
type RedisDB struct {
	c    *redis.Pool
	host string
	db   int
}

//Init desc
//@method Init desc: initialization Redis DB
//@param (string) redis host address
//@param (int) redis db code
//@param (int) redis connection max idle of number
//@param (int) redis connection max active of number
//@param (int) redis connection idle time (unit/sec)
//@return (error) if connecting fail return error ,success return nil
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

//Do desc
//@method Do desc: execute redis command
//@param (string) command name
//@param (...interface{}) command params
//@return (interface{}) execute result
//@return (error) if execute fail return error, execute success return nil
func (slf *RedisDB) Do(commandName string, args ...interface{}) (interface{}, error) {
	c := slf.c.Get()
	defer c.Close()
	return c.Do(commandName, args...)
}

//Close desc
//@method Close desc: close redis db operation
func (slf *RedisDB) Close() {
	slf.c.Close()
}
