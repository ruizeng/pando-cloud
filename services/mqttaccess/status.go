package main

import (
	"github.com/PandoCloud/pando-cloud/pkg/protocol"
)

var StatusChan map[uint64]chan *protocol.Data

func init() {
	StatusChan = make(map[uint64]chan *protocol.Data)
}
