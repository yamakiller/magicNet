package mutex

import (
	"fmt"
	"log"

	"github.com/gomodule/redigo/redis"
)

//RedisMethodDo doc
//@type (func(commandName string, args ...interface{}) (interface{}, error))
type RedisMethodDo func(commandName string, args ...interface{}) (interface{}, error)

//RedisMutex doc
//@Summary redis mutex object
//@Member (string) resource
//@Member (string) token
//@Member (int)    timeout
type RedisMutex struct {
	_doFun    RedisMethodDo
	_resource string
	_token    string
	_timeout  int
}

func (slf *RedisMutex) tryLock() (ok bool, err error) {
	_, err = slf._doFun("SET", slf.key(), slf._token, "EX", int(slf._timeout), "NX")
	if err == redis.ErrNil {
		return false, nil
	}

	if err != nil {
		return false, err
	}
	return true, nil
}

func (slf *RedisMutex) key() string {
	return fmt.Sprintf("redislock:%s", slf._resource)
}

//Unlock doc
//@Method Unlock @Summary unlocking
//@Return (error) unlock fail returns error informat
func (slf *RedisMutex) Unlock() (err error) {
	_, err = slf._doFun("del", slf.key())
	return
}

//AddTimeOut doc
//@Summary rest/append lock timeout time
//@Param (int64) setting/append time
//@Return (bool)
//@reutrn (error)
func (slf *RedisMutex) AddTimeOut(exTime int64) (ok bool, err error) {
	ttl, err := redis.Int64(slf._doFun("TTL", slf.key()))
	if err != nil {
		log.Fatal("redis get failed:", err)
	}
	if ttl > 0 {
		_, err := redis.String(slf._doFun("SET", slf.key(), slf._token, "EX", int(ttl+exTime)))
		if err == redis.ErrNil {
			return false, nil
		}
		if err != nil {
			return false, err
		}
	}
	return false, nil
}

//TryRedisLock Try to acquire a lock
//@Summary Try to acquire a locking
//@Param (RedisMethodDo) redis do function
//@Param (string) lock object(key/name)
//@Param (string) lock token
//@Param (int)    lock timeout millsec
//@Return (*RedisMutex) redis mutex object
//@Return (bool) redis lock is success
//@Return (error) redis lock fail error informat
func TryRedisLock(doFun RedisMethodDo,
	resouse string,
	token string,
	timeout int) (m *RedisMutex, ok bool, err error) {
	return TryRedisLockWithTimeOut(doFun, resouse, token, timeout)
}

//TryRedisLockWithTimeOut doc
//@Summary Try to acquire the lock and set the lock timeout
//@Param (RedisMethodDo) redis do function
//@Param (string) lock object(key/name)
//@Param (string) lock token
//@Param (int)    lock timeout millsec
//@Return (*RedisMutex) redis mutex object
//@Return (bool) redis lock is success
//@Return (error) redis lock fail error informat
func TryRedisLockWithTimeOut(doFun RedisMethodDo,
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
