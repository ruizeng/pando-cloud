// RPCClient implements a rpc client tool with reconnect and load balance.
package server

import (
	"fmt"
	"math/rand"
	"net/rpc"
	"time"
)

type RPCClient struct {
	clients map[string]*rpc.Client
	random  *rand.Rand
}

func NewRPCClient() (*RPCClient, error) {
	if serverInstance == nil {
		return nil, errorf(errServerNotInit)
	}
	if serverInstance.svrmgr == nil {
		return nil, errorf(errServerManagerNotInit)
	}
	return &RPCClient{
		clients: make(map[string]*rpc.Client),
		random:  rand.New(rand.NewSource(time.Now().UnixNano())),
	}, nil
}

func rpcCallWithReconnect(client *rpc.Client, addr string, serverMethod string, args interface{}, reply interface{}) error {
	err := client.Call(serverMethod, args, reply)
	if err == rpc.ErrShutdown {
		client, err = rpc.Dial("tcp", addr)
		if err != nil {
			return err
		}
	} else if err == nil {
		return nil
	}

	err = client.Call(serverMethod, args, reply)

	return err
}

//RPC call with reconnect and retry.
func (client *RPCClient) Call(severName string, serverMethod string, args interface{}, reply interface{}) error {
	addrs, err := serverInstance.svrmgr.GetServerHosts(severName, FlagRPCHost)
	if err != nil {
		return err
	}

	// pick a random start index for round robin
	total := len(addrs)
	start := client.random.Intn(total)

	for idx := 0; idx < total; idx++ {
		addr := addrs[(start+idx)%total]
		mapkey := fmt.Sprintf("%s[%s]", severName, addr)
		if client.clients[mapkey] == nil {
			client.clients[mapkey], err = rpc.Dial("tcp", addr)
			if err != nil {
				Log.Warnf("RPC dial error : %s", err)
				continue
			}
		}

		err = rpcCallWithReconnect(client.clients[mapkey], addr, serverMethod, args, reply)
		if err != nil {
			Log.Warnf("RpcCallWithReconnect error : %s", err)
			continue
		}

		return nil
	}

	return errorf(err.Error())
}
