package main

import (
	"flag"
)

const (
	flagRabbitHost    = "rabbithost"
	defaultRabbitHost = "amqp://guest:guest@localhost:5672/"
)

var (
	confRabbitHost = flag.String(flagRabbitHost, defaultRabbitHost, "rabbitmq host address, amqp://user:password@ip:port/")
)
