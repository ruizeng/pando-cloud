// package server provides service interfaces and libraries.
// including:
// tcp/http server library.
// rpc service library with addon functionality.
// service discory and registration Logic.
// statistic lib.
package server

import (
	"fmt"
	"github.com/vharitonsky/iniflags"
	"net/http"
	"net/rpc"
	"time"
)

// server is a singleton
var serverInstance *Server = nil

// Server
type Server struct {
	// required
	name string
	// optional
	rpcsvr    *RPCServer  // RPC server
	tcpsvr    *TCPServer  // TCP server
	httpsvr   *HTTPServer // HTTP server
	timertask TimerTask   // timer task
	// functions
	// svcmgr *ServiceManager // service registration&discovery manager
}

// init the server with specific name.
func Init(name string) error {
	if serverInstance == nil {
		// read config
		iniflags.Parse()

		// read network info
		readNetInterfaces()

		// log
		initLog(name, *confLogLevel)

		// server instance
		serverInstance = &Server{
			name: name,
		}

		// init service manager
		/*
			serverInstance.svcmgr, err = NewServiceManager(name)
			if err != nil {
				return err
			}
			fmt.Printf("service manager init! \n")
		*/

	}
	return nil
}

// register TCP handler class
func RegisterTCPHandler(handler TCPHandler) error {
	if serverInstance == nil {
		return errorf(errServerNotInit)
	}
	if serverInstance.tcpsvr == nil {
		if *confTCPHost == "" {
			return errorf(errMissingFlag, flagTCPHost)
		}

		addr, err := fixHostIp(*confTCPHost)
		if err != nil {
			return errorf(errWrongHostAddr, confTCPHost)
		}

		serverInstance.tcpsvr = &TCPServer{
			addr:    addr,
			handler: handler,
			useTls:  *confUseTls,
		}
	}
	return nil
}

// register HTTP handler class
func RegisterHTTPHandler(handler http.Handler) error {
	if serverInstance == nil {
		return errorf(errServerNotInit)
	}
	if serverInstance.httpsvr == nil {
		if *confHTTPHost == "" {
			return errorf(errMissingFlag, flagHTTPHost)
		}

		addr, err := fixHostIp(*confHTTPHost)
		if err != nil {
			return errorf(errWrongHostAddr, flagHTTPHost)
		}

		serverInstance.httpsvr = &HTTPServer{
			addr:     addr,
			handler:  handler,
			useHttps: *confUseHttps,
		}
	}
	return nil
}

// register RPC handler class
func RegisterRPCService(rcvr interface{}) error {
	if serverInstance == nil {
		return errorf(errServerNotInit)
	}
	if serverInstance.rpcsvr == nil {
		if *confRPCHost == "" {
			return errorf(errMissingFlag, flagRPCHost)
		}

		addr, err := fixHostIp(*confRPCHost)
		if err != nil {
			return errorf(errWrongHostAddr, *confRPCHost)
		}

		err = rpc.Register(rcvr)
		if err != nil {
			return errorf("Cannot Resgister RPC service: %s", err)
		}

		handler := rpcHandler{}

		serverInstance.rpcsvr = &RPCServer{
			TCPServer{
				addr:    addr,
				handler: &handler,
				useTls:  false, // rpc service do not use tls because it's in internal network
			},
		}
	}
	return nil
}

//register timer task
func RegisterTimerTask(task TimerTask) error {
	if serverInstance == nil {
		return errorf(errServerNotInit)
	}
	if serverInstance.timertask == nil {
		serverInstance.timertask = task
	}
	return nil
}

// start service
func Run() error {
	if serverInstance == nil {
		return errorf(errServerNotInit)
	}

	if serverInstance.tcpsvr != nil {
		err := serverInstance.tcpsvr.Start()
		if err != nil {
			return err
		}
		Log.Info("starting tcp server ... OK")
	}

	if serverInstance.httpsvr != nil {
		fmt.Println("serverInstance.httpsvr")
		err := serverInstance.httpsvr.Start()
		if err != nil {
			return err
		}
		Log.Info("starting http server ... OK")
	}

	if serverInstance.rpcsvr != nil {
		fmt.Println("serverInstance.rpcsvr \n")
		err := serverInstance.rpcsvr.Start()
		if err != nil {
			return err
		}
		Log.Info("starting rpc server ... OK")
	}

	Log.Info("sever launch successfully!")

	// loop to do something
	for {
		// master update
		/*
			err := serverInstance.svcmgr.RegisterServer()
			if err != nil {
				Log.Error("RegisterServer error: %s", err)
			} else {
				Log.Info("RegisterServer Success")
			}
			err = serverInstance.svcmgr.UpdateServerHosts()
			if err != nil {
				Log.Error("UpdateServerHosts error: %s", err)
			} else {
				Log.Info("UpdateServerHosts Success")
			}
		*/

		//timer task
		if serverInstance.timertask != nil {
			serverInstance.timertask.DoTask()
		}

		time.Sleep(60 * time.Second)
	}

	return nil
}
