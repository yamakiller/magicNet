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

func (m *RedisMutex) tryLock() (ok bool, err error) {
	_, err = m.doFun("SET", m.key(), m.token, "EX", int(m.timeout), "NX")
	if err == redis.ErrNil {
		return false, nil
	}

	if err != nil {
		return false, err
	}
	return true, nil
}

func (m *RedisMutex) key() string {
	return fmt.Sprintf("redislock:%s", m.resource)
}

//Unlock Unlock
func (m *RedisMutex) Unlock() (err error) {
	_, err = m.doFun("del", m.key())
	return
}

//AddTimeout Set lock timeout
func (m *RedisMutex) AddTimeout(exTime int64) (ok bool, err error) {
	ttl, err := redis.Int64(m.doFun("TTL", m.key()))
	if err != nil {
		log.Fatal("redis get failed:", err)
	}
	if ttl > 0 {
		_, err := redis.String(m.doFun("SET", m.key(), m.token, "EX", int(ttl+exTime)))
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
