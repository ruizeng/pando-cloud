package redispool

import (
	"github.com/garyburd/redigo/redis"
	"time"
)

var mapRedisPool map[string]*redis.Pool

// GetClient get a redis connection by host
func GetClient(host string) (redis.Conn, error) {

	pool, exist := mapRedisPool[host]
	if !exist {
		pool = &redis.Pool{
			MaxIdle: 10,
			Dial: func() (redis.Conn, error) {
				c, err := redis.Dial("tcp", host)
				if err != nil {
					return nil, err
				}
				return c, err
			},
			TestOnBorrow: func(c redis.Conn, t time.Time) error {
				_, err := c.Do("PING")
				return err
			},
		}
		mapRedisPool[host] = pool
	}
	return pool.Get(), nil
}

func init() {
	mapRedisPool = make(map[string]*redis.Pool)
}
