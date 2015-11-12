package redispool

import (
	"github.com/garyburd/redigo/redis"
	"testing"
)

func TestRedisCli(t *testing.T) {
	cli, err := GetClient("localhost:6379")
	if err != nil {
		t.Fatal(err)
	}

	_, err = redis.String(cli.Do("GET", "testkey"))
	if err != nil && err != redis.ErrNil {
		t.Error(err)
	}
}
