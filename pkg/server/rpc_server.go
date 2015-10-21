// rpc server
package server

import (
	"net"
	"net/rpc"
)

type RPCServer struct {
	TCPServer
}

type rpcHandler struct{}

func (handler *rpcHandler) Handle(conn net.Conn) {
	rpc.ServeConn(conn)
}
