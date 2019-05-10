/*
   @Time : 2019-05-10 09:55
   @Author : frozenchen
   @File : config
   @Software: studio-library
*/
package redis

import (
	"github.com/garyburd/redigo/redis"
	"time"
)

type Redis struct {
	pool *redis.Pool
}

type Config struct {
	Addr         string
	Auth         string
	Idle         int
	Active       int
	DialTimeout  time.Duration
	ReadTimeout  time.Duration
	WriteTimeout time.Duration
	IdleTimeout  time.Duration
}

func New(c *Config) *Redis {
	var option []redis.DialOption
	option = append(option, redis.DialReadTimeout(c.ReadTimeout))
	option = append(option, redis.DialWriteTimeout(c.WriteTimeout))
	option = append(option, redis.DialPassword(c.Auth))
	option = append(option, redis.DialConnectTimeout(c.DialTimeout))

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
		IdleTimeout: c.IdleTimeout,
	}
	return &Redis{pool: pool}
}

func (r *Redis) Get() redis.Conn {
	return r.pool.Get()
}
