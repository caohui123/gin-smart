package lib

import "C"
import (
	"errors"
	"os"
	"os/signal"
	"syscall"
	"time"

	redigo "github.com/gomodule/redigo/redis"
)

type RedisPool struct {
	*redigo.Pool
}

func NewRedis(appConf *Config) (*RedisPool, error) {
	section := "redis"
	if v, _ := appConf.Get(section, "enable"); v == "false" {
		return nil, nil
	}
	conf, err := appConf.Section(section)
	if err != nil {
		return nil, errors.New("start redis failed! couldn't get config")
	}
	pool := &redigo.Pool{
		MaxIdle:     5,                 // idle的列表长度, 空闲的线程数
		MaxActive:   0,                 // 线程池的最大连接数， 0表示没有限制
		Wait:        true,              // 当连接数已满，是否要阻塞等待获取连接。false表示不等待，直接返回错误。
		IdleTimeout: 200 * time.Second, // 最大的空闲连接等待时间，超过此时间后，空闲连接将被关闭
		Dial: func() (redigo.Conn, error) { // 创建链接
			c, err := redigo.Dial("tcp", conf["host"])
			if err != nil {
				return nil, err
			}
			if conf["password"] != "" {
				if _, err := c.Do("AUTH", conf["password"]); err != nil {
					c.Close()
					return nil, err
				}
			}
			if _, err := c.Do("SELECT", conf["db_num"]); err != nil {
				c.Close()
				return nil, err
			}
			return c, nil
		},
		TestOnBorrow: func(c redigo.Conn, t time.Time) error { // 一个测试链接可用性
			if time.Since(t) < time.Minute {
				return nil
			}
			_, err := c.Do("PING")
			return err
		},
	}
	return &RedisPool{pool}, nil
}

func (pool *RedisPool) closePool() {
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	signal.Notify(c, syscall.SIGTERM)
	signal.Notify(c, syscall.SIGKILL)
	go func() {
		<-c
		pool.Close()
	}()
}

// get
func (pool *RedisPool) GetKey(key string) (string, error) {
	rds := pool.Get()
	defer rds.Close()
	return redigo.String(rds.Do("GET", key))
}

// set expires为0时，表示永久性存储
func (pool *RedisPool) SetKey(key string, value interface{}, expires int64) error {
	rds := pool.Get()
	defer rds.Close()
	if expires == 0 {
		_, err := rds.Do("SET", key, value)
		return err
	} else {
		_, err := rds.Do("SETEX", key, expires, value)
		return err
	}
}

// del
func (pool *RedisPool) DelKey(key string) error {
	rds := pool.Get()
	defer rds.Close()
	_, err := rds.Do("DEL", key)
	return err
}

// lrange
func (pool *RedisPool) LRange(key string, start, stop int64) ([]string, error) {
	rds := pool.Get()
	defer rds.Close()
	return redigo.Strings(rds.Do("LRANGE", key, start, stop))
}

// lpop
func (pool *RedisPool) LPop(key string) (string, error) {
	rds := pool.Get()
	defer rds.Close()
	return redigo.String(rds.Do("LPOP", key))
}

// LPush
func (pool *RedisPool) LPush(key, value interface{}) error {
	rds := pool.Get()
	defer rds.Close()
	_, err := rds.Do("LPUSH", key, value)
	return err
}

// LPushAndTrimKey
func (pool *RedisPool) LPushAndTrimKey(key, value interface{}, size int64) error {
	rds := pool.Get()
	defer rds.Close()
	rds.Send("MULTI")
	rds.Send("LPUSH", key, value)
	rds.Send("LTRIM", key, size-2*size, -1)
	_, err := rds.Do("EXEC")
	return err
}

// RPushAndTrimKey
func (pool *RedisPool) RPushAndTrimKey(key, value interface{}, size int64) error {
	rds := pool.Get()
	defer rds.Close()
	rds.Send("MULTI")
	rds.Send("RPUSH", key, value)
	rds.Send("LTRIM", key, size-2*size, -1)
	_, err := rds.Do("EXEC")
	return err
}

// ExistsKey
func (pool *RedisPool) ExistsKey(key string) (bool, error) {
	rds := pool.Get()
	defer rds.Close()
	return redigo.Bool(rds.Do("EXISTS", key))
}

// ttl 返回剩余时间
func (pool *RedisPool) TTLKey(key string) (int64, error) {
	rds := pool.Get()
	defer rds.Close()
	return redigo.Int64(rds.Do("TTL", key))
}

func (pool *RedisPool) ExpireKey(key string, expires int) (bool, error) {
	rds := pool.Get()
	defer rds.Close()
	return redigo.Bool(rds.Do("EXPIRE", key, expires))
}

// incr 自增
func (pool *RedisPool) Incr(key string) (int64, error) {
	rds := pool.Get()
	defer rds.Close()
	return redigo.Int64(rds.Do("INCR", key))
}

// Decr 自减
func (pool *RedisPool) Decr(key string) (int64, error) {
	rds := pool.Get()
	defer rds.Close()
	return redigo.Int64(rds.Do("DECR", key))
}

// mset 批量写入 rds.Do("MSET", "ket1", "value1", "key2","value2")
func (pool *RedisPool) MsetKey(keyValue ...interface{}) error {
	rds := pool.Get()
	defer rds.Close()
	_, err := rds.Do("MSET", keyValue...)
	return err
}

// mget  批量读取 mget key1, key2, 返回map结构
func (pool *RedisPool) MgetKey(keys ...interface{}) map[interface{}]string {
	rds := pool.Get()
	defer rds.Close()
	values, _ := redigo.Strings(rds.Do("MGET", keys...))
	resultMap := map[interface{}]string{}
	keyLen := len(keys)
	for i := 0; i < keyLen; i++ {
		resultMap[keys[i]] = values[i]
	}
	return resultMap
}

// hmset 同时将多个 field-value (域-值)对设置到哈希表 key 中。
func (pool *RedisPool) HMsetKey(key string, simpleObject interface{}) error {
	rds := pool.Get()
	defer rds.Close()
	_, err := rds.Do("HMSET", redigo.Args{}.Add(key).AddFlat(simpleObject)...)
	return err
}

// hmget 返回哈希表 key 中，一个或多个给定域的值
func (pool *RedisPool) HMgetKey(key string, simpleObject interface{}) error {
	rds := pool.Get()
	defer rds.Close()
	values, _ := redigo.Values(rds.Do("HGETALL", key))

	err := redigo.ScanStruct(values, simpleObject)
	if err != nil {
		return err
	}
	return err
}
