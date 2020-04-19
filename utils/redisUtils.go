package utils

import (
redisPool "StorePanAPI/catch/redis"
	"github.com/garyburd/redigo/redis"
)

const (
	redisLockTimeout = 10 // 10 seconds
)

func TryLock(k string) (isLock bool, err error) {
	// 获取连接
	conn := redisPool.RedisPool().Get()
	// 关闭连接
	defer conn.Close()
	// 这里需要redis.String包一下，才能返回redis.ErrNil
	_, err = redis.String(conn.Do("SET", k, 1, "ex", redisLockTimeout, "nx"))
	if err != nil {
		if err == redis.ErrNil {
			err = nil
			return
		}
		return
	}
	isLock = true
	return
}

func Unlock(k string) (err error) {
	// 获取连接
	conn := redisPool.RedisPool().Get()
	// 关闭连接
	defer conn.Close()
	_, err = conn.Do("DEL", k)
	if err != nil {
		return
	}
	return
}


func Get(k string) (interface{}, error) {
	// 获取连接
	conn := redisPool.RedisPool().Get()
	defer conn.Close()
	v, err := conn.Do("GET", k)
	if err != nil {
		return nil, err
	}
	return v, nil
}

func Set(k string,v interface{}) error {
	// 获取连接
	conn := redisPool.RedisPool().Get()
	defer conn.Close()
	_, err := conn.Do("SET", k, v)
	return err
}

func SetEx(k string, v interface{}, ex int64) error {
	// 获取连接
	conn := redisPool.RedisPool().Get()
	defer conn.Close()
	_, err := conn.Do("SET", k, v, "EX", ex)
	return err
}
