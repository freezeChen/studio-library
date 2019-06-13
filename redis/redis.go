/*
   @Time : 2019-05-10 09:55
   @Author : frozenchen
   @File : redis
   @Software: studio-library
*/
package redis

import (
	"context"
	"encoding/json"
	"errors"
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

func (r *Redis) CacheGet(ctx context.Context, key string) (reply *Reply) {
	reply = new(Reply)
	if !legalKey(key) {
		reply.err = errors.New("empty key")
		return
	}
	conn := r.Get()
	defer conn.Close()

	reply2, err := redis.Bytes(conn.Do("GET", key))
	reply.err = err
	reply.item = new(Item)
	reply.item.Value = reply2
	return
}

func (r *Redis) CaCheSet(ctx context.Context, item *Item) (err error) {
	if !legalKey(item.Key) {
		return ErrKey
	}

	conn := r.Get()
	defer conn.Close()

	if item.Value != nil {
		_, err = conn.Do("SET", item.Key, item.Value, "EX", item.Expiration)
		return
	}

	var value string
	switch item.Object.(type) {
	case string:
		value = item.Object.(string)
	case []byte:
		value = string(item.Object.([]byte))
	case interface{}:
		var b []byte
		b, err = json.Marshal(item.Object)
		value = string(b)
		if err != nil {
			return
		}
	}

	_, err = conn.Do("SET", item.Key, value, "EX", item.Expiration)
	return
}

func (r *Redis) CacheGetMulti(ctx context.Context, keys []string) (res *Replies, err error) {
	conn := r.Get()
	defer conn.Close()

	res = new(Replies)
	res.items = make(map[string]*Item, len(keys))
	for _, v := range keys {
		var b []byte
		b, err = redis.Bytes(conn.Do("GET", v))
		if err != nil {
			return
		}
		res.items[v] = &Item{Key: v, Value: b}
	}

	return
}

func (r *Redis) Delete(ctx context.Context, key string) (err error) {
	if !legalKey(key) {
		err = errors.New("empty key")
		return
	}

	conn := r.Get()
	defer conn.Close()

	_, err = conn.Do("DEL", key)
	return
}

func legalKey(key string) bool {
	if len(key) == 0 {
		return false
	}

	return true
}
