package main

import (
	"flag"
)

const (
	flagMongoHost    = "mongohost"
	defaultMongoHost = "localhost"

	flagRabbitHost    = "rabbithost"
	defaultRabbitHost = "amqp://guest:guest@localhost:5672/"
)

var (
	confMongoHost  = flag.String(flagMongoHost, defaultMongoHost, "mongo host address, ip:port")
	confRabbitHost = flag.String(flagRabbitHost, defaultRabbitHost, "rabbitmq host address, amqp://user:password@ip:port/")
)
