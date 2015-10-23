// tcp server library.
package server

import (
	"crypto/tls"
	"net"
)

type TCPHandler interface {
	Handle(net.Conn)
}

type TCPServer struct {
	addr    string
	handler TCPHandler
	useTls  bool
}

// start will keep accepting and serving tcp connections.
func (ts *TCPServer) Start() error {
	// check for conditions
	if ts.handler == nil {
		return errorf(errTCPHandlerNotRegistered)
	}
	// listen
	var ln net.Listener
	var err error
	if ts.useTls {
		// if use tls, then load pem files and start server
		if *confCAFile == "" {
			return errorf(errMissingFlag, FlagCAFile)
		}
		if *confKeyFile == "" {
			return errorf(errMissingFlag, FlagKeyFile)
		}

		// process key files
		cert, err := tls.LoadX509KeyPair(*confCAFile, *confKeyFile)
		if err != nil {
			return errorf(errLoadSecureKey, err.Error())
		}

		// config server with tls
		config := tls.Config{Certificates: []tls.Certificate{cert}}

		// listen for new connection
		ln, err = tls.Listen("tcp", ts.addr, &config)

		if err != nil {
			return errorf(errListenFailed, ts.addr, err)
		}

	} else {
		// don't use tls, just listen
		ln, err = net.Listen("tcp", ts.addr)
		if err != nil {
			return errorf(errListenFailed, ts.addr, err)
		}
	}

	Log.Infof("TCP Server Listen on %s, use tls: %v", ts.addr, ts.useTls)

	// continously accept connections and serve, nonblock
	go func() {
		for {
			conn, err := ln.Accept()
			if err != nil {
				Log.Errorf(errNewConnection, err.Error())
				continue
			}
			Log.Infof("accepting new connection %s", conn.RemoteAddr())
			go ts.handler.Handle(conn)
		}
	}()

	return nil
}
