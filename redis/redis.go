/*
   @Time : 2019-05-10 09:55
   @Author : frozenchen
   @File : redis
   @Software: studio-library
*/
package redis

import (
	"github.com/garyburd/redigo/redis"
	"time"
)

const (
	JSONFLAG = "json"
)

type Redis struct {
	pool *redis.Pool
}

func New(c *Config) *Redis {
	var option []redis.DialOption
	option = append(option, redis.DialReadTimeout(time.Duration(c.ReadTimeout)))
	option = append(option, redis.DialWriteTimeout(time.Duration(c.WriteTimeout)))
	option = append(option, redis.DialPassword(c.Auth))
	option = append(option, redis.DialConnectTimeout(time.Duration(c.DialTimeout)))

	pool := &redis.Pool{
		Dial: func() (conn redis.Conn, e error) {
			return redis.Dial("tcp", c.Addr, option...)
		},
		TestOnBorrow: func(c redis.Conn, t time.Time) error {
			_, err := c.Do("PING")
			return err
		},
		MaxIdle:     c.Idle,
		MaxActive:   c.Active,
		IdleTimeout: time.Duration(c.IdleTimeout),
	}
	return &Redis{pool: pool}
}

func (r *Redis) Get() redis.Conn {
	return r.pool.Get()
}

//func (r *Redis) CacheGet(key string) (*Item, error) {
//	if !legalKey(key) {
//		return nil, ErrorKey
//	}
//
//}

func legalKey(key string) bool {
	if len(key) == 0 {
		return false
	}

	return true
}
