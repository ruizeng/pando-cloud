// package server provides service interfaces and libraries.
// including:
// tcp/http server library.
// rpc service library with addon functionality.
// service discory and registration Logic.
// statistic lib.
package server

import (
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
	svrmgr *ServerManager // service registration&discovery manager
	rpccli *RPCClient     // rpc client
}

// init the server with specific name.
func Init(name string) error {
	if serverInstance == nil {
		// read config
		iniflags.Parse()

		// read network info
		readNetInterfaces()

		// log
		err := initLog(name, *confLogLevel)
		if err != nil {
			return err
		}

		// server instance
		serverInstance = &Server{
			name: name,
		}

		// init service manager
		serverInstance.svrmgr, err = NewServerManager(name, *confEtcd)
		if err != nil {
			return err
		}

		// create RPC client
		serverInstance.rpccli, err = NewRPCClient()
		if err != nil {
			return err
		}

		Log.Infof("server %s init success.", name)

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
			return errorf(errMissingFlag, FlagTCPHost)
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
			return errorf(errMissingFlag, FlagHTTPHost)
		}

		addr, err := fixHostIp(*confHTTPHost)
		if err != nil {
			return errorf(errWrongHostAddr, FlagHTTPHost)
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
func RegisterRPCHandler(rcvr interface{}) error {
	if serverInstance == nil {
		return errorf(errServerNotInit)
	}
	if serverInstance.rpcsvr == nil {
		if *confRPCHost == "" {
			return errorf(errMissingFlag, FlagRPCHost)
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

// register timer task
func RegisterTimerTask(task TimerTask) error {
	if serverInstance == nil {
		return errorf(errServerNotInit)
	}
	if serverInstance.timertask == nil {
		serverInstance.timertask = task
	}
	return nil
}

// rpc call by name
func RPCCallByName(serverName string, serverMethod string, args interface{}, reply interface{}) error {
	if serverInstance == nil {
		return errorf(errServerNotInit)
	}

	return serverInstance.rpccli.Call(serverName, serverMethod, args, reply)
}

// get server's hosts by server name and service type
func GetServerHosts(serverName string, hostType string) ([]string, error) {
	if serverInstance == nil {
		return nil, errorf(errServerNotInit)
	}

	return serverInstance.svrmgr.GetServerHosts(serverName, hostType)
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
		err := serverInstance.httpsvr.Start()
		if err != nil {
			return err
		}
		Log.Info("starting http server ... OK")
	}

	if serverInstance.rpcsvr != nil {
		err := serverInstance.rpcsvr.Start()
		if err != nil {
			return err
		}
		Log.Info("starting rpc server ... OK")
	}

	Log.Info("sever launch successfully!")

	// loop to do something
	for {
		// server manager update
		err := serverInstance.svrmgr.RegisterServer()
		if err != nil {
			Log.Warnf("RegisterServer error: %s", err)
		} else {
			Log.Info("RegisterServer Success")
		}
		err = serverInstance.svrmgr.UpdateServerHosts()
		if err != nil {
			Log.Error("UpdateServerHosts error: %s", err)
		} else {
			Log.Info("UpdateServerHosts Success")
		}

		//timer task
		if serverInstance.timertask != nil {
			serverInstance.timertask.DoTask()
		}

		time.Sleep(60 * time.Second)
	}

	return nil
}
