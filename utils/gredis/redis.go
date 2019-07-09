package gredis

import (
	"encoding/json"
	"github.com/gomodule/redigo/redis"
	"xhblog/utils/logging"
	"xhblog/utils/setting"
)

var RedisConn *redis.Pool

func Setup() error {
	RedisConn = &redis.Pool{
		MaxIdle:     setting.RedisSetting.MaxIdle,
		MaxActive:   setting.RedisSetting.MaxActive,
		IdleTimeout: setting.RedisSetting.IdleTimeout,
		Dial: func() (conn redis.Conn, e error) {
			conn, err := redis.Dial("tcp", setting.RedisSetting.Host)
			if err != nil {
				return nil, err
			}
			if setting.RedisSetting.Password != "" {
				if _, err := conn.Do("AUTH", setting.RedisSetting.Password); err != nil {
					conn.Close()
					return nil, err
				}
			}
			return conn, err
		 },

	}
	return nil
}

func Set(key string, data interface{}, time int) (bool, error) {
 	conn := RedisConn.Get()
 	defer conn.Close()

 	value, err := json.Marshal(data)
	if err != nil {
		return false, err
	}

 	reply, err := redis.Bool(conn.Do("SET", key, value))
 	conn.Do("EXPIRE", key, time)

	return reply, err
}

func Get(key string) ([]byte, error) {
	conn := RedisConn.Get()
	defer conn.Close()

	reply, err := redis.Bytes(conn.Do("GET", key))
	if err != nil {
		return nil, err
	}

	return reply, err
}

func Exists(key string) (bool) {
	conn := RedisConn.Get()
	defer conn.Close()

	reply, err := redis.Bool(conn.Do("EXISTS", key))
	if err != nil {
		return false
	}

	return reply
}

func Delete(key string) (bool, error) {
	conn := RedisConn.Get()
	defer conn.Close()

	reply, err := redis.Bool(conn.Do("DEL", key))
	if err != nil {
		logging.Info("delete cache failed")
		logging.Info(err)
	}
	return reply, err
}

