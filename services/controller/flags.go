package main

import (
	"flag"
)

const (
	flagMongoHost = "mongohost"

	defaultMongoHost = "localhost"
)

var (
	confMongoHost = flag.String(flagMongoHost, defaultMongoHost, "mongo host address, ip:port")
)
