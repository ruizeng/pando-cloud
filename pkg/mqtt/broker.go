package mqtt

import (
	"net"
)

type Broker struct {
	context *Context
}

func NewBroker(p Provider) *Broker {
	// context
	cxt := NewContext(p)

	handler := &Broker{context: cxt}

	return handler
}

func (h *Broker) Handle(conn net.Conn) {
	host := conn.RemoteAddr().String()
	h.context.NewSub(host, conn)
}
