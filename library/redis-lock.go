package library

import (
	"fmt"
	"log"

	"github.com/gomodule/redigo/redis"
)

//RedisMethodDo Execution method
type RedisMethodDo func(commandName string, args ...interface{}) (interface{}, error)

//RedisMutex Distributed sync lock
type RedisMutex struct {
	doFun    RedisMethodDo
	resource string
	token    string
	timeout  int
}

func (slf *RedisMutex) tryLock() (ok bool, err error) {
	_, err = slf.doFun("SET", slf.key(), slf.token, "EX", int(slf.timeout), "NX")
	if err == redis.ErrNil {
		return false, nil
	}

	if err != nil {
		return false, err
	}
	return true, nil
}

func (slf *RedisMutex) key() string {
	return fmt.Sprintf("redislock:%s", slf.resource)
}

//Unlock Unlock
func (slf *RedisMutex) Unlock() (err error) {
	_, err = slf.doFun("del", slf.key())
	return
}

//AddTimeout Set lock timeout
func (slf *RedisMutex) AddTimeout(exTime int64) (ok bool, err error) {
	ttl, err := redis.Int64(slf.doFun("TTL", slf.key()))
	if err != nil {
		log.Fatal("redis get failed:", err)
	}
	if ttl > 0 {
		_, err := redis.String(slf.doFun("SET", slf.key(), slf.token, "EX", int(ttl+exTime)))
		if err == redis.ErrNil {
			return false, nil
		}
		if err != nil {
			return false, err
		}
	}
	return false, nil
}

//TryLock Try to acquire a lock
func TryLock(doFun RedisMethodDo,
	resouse string,
	token string,
	timeout int) (m *RedisMutex, ok bool, err error) {
	return TryLockWithTimeout(doFun, resouse, token, timeout)
}

//TryLockWithTimeout Try to acquire the lock and set the lock timeout
func TryLockWithTimeout(doFun RedisMethodDo,
	resouse string,
	token string,
	timeout int) (m *RedisMutex, ok bool, err error) {
	m = &RedisMutex{doFun, resouse, token, timeout}
	ok, err = m.tryLock()
	if !ok || err != nil {
		m = nil
	}
	return
}
