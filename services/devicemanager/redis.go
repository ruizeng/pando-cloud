package main

import (
	"flag"
)

const (
	flagRedisHost = "redishost"

	defaultRedisHost = "localhost:6379"
)

var (
	confRedisHost = flag.String(flagRedisHost, defaultRedisHost, "redis host address, ip:port")
)
