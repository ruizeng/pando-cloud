// service registration and discovery

package server

import (
	"errors"
	"github.com/coreos/etcd/client"
	"golang.org/x/net/context"
	"strings"
	"time"
)

const (
	EtcdServersPrefix    = "/pando/servers/"
	EtcdServersPrefixCnt = 2
)

type ServerManager struct {
	serverName string
	// servername -> hosttype -> hostlist
	// eg. var hosts []string = mapServers["testserver"]["rpchost"]
	mapServers map[string](map[string][]string)
	etcdHosts  []string
}

// etcd hosts is config as http://ip1:port1;http://ip2:port2;http://ip3:port3
func NewServerManager(name string, etcd string) (*ServerManager, error) {
	if etcd == "" {
		return nil, errors.New("no etcd host found!")
	}
	return &ServerManager{
		serverName: name,
		etcdHosts:  strings.Split(etcd, ";"),
		mapServers: make(map[string](map[string][]string)),
	}, nil
}

// register server to etcd
func (mgr *ServerManager) RegisterServer() error {
	if serverInstance == nil {
		return errorf(errServerNotInit)
	}
	cfg := client.Config{
		Endpoints: mgr.etcdHosts,
		Transport: client.DefaultTransport,
		// set timeout per request to fail fast when the target endpoint is unavailable
		HeaderTimeoutPerRequest: time.Second,
	}
	c, err := client.New(cfg)
	if err != nil {
		return err
	}
	kapi := client.NewKeysAPI(c)
	prefix := EtcdServersPrefix + mgr.serverName + "/"
	var response *client.Response
	opt := &client.SetOptions{TTL: 180 * time.Second}
	if serverInstance.tcpsvr != nil {
		addr, _ := fixHostIp(*confTCPHost)
		response, err = kapi.Set(context.Background(), prefix+FlagTCPHost+"/"+addr, addr, opt)
	}
	if serverInstance.rpcsvr != nil {
		addr, _ := fixHostIp(*confRPCHost)
		response, err = kapi.Set(context.Background(), prefix+FlagRPCHost+"/"+addr, addr, opt)
	}
	if serverInstance.httpsvr != nil {
		addr, _ := fixHostIp(*confHTTPHost)
		response, err = kapi.Set(context.Background(), prefix+FlagHTTPHost+"/"+addr, addr, opt)
	}
	if err != nil {
		return err
	}
	// print common key info
	Log.Infof("RegisterServer is done. Metadata is %q\n", response)

	return nil
}

// update server hosts
func (mgr *ServerManager) UpdateServerHosts() error {
	if serverInstance == nil {
		return errorf(errServerNotInit)
	}

	cfg := client.Config{
		Endpoints: mgr.etcdHosts,
		Transport: client.DefaultTransport,
		// set timeout per request to fail fast when the target endpoint is unavailable
		HeaderTimeoutPerRequest: time.Second,
	}
	c, err := client.New(cfg)
	if err != nil {
		return err
	}

	kapi := client.NewKeysAPI(c)
	prefix := EtcdServersPrefix
	opt := &client.GetOptions{Recursive: true}
	response, err := kapi.Get(context.Background(), prefix, opt)
	if err != nil {
		return err
	}

	servers := make(map[string](map[string][]string))

	root := response.Node
	if root.Dir != true {
		return errorf(errWrongEtcdPath, root.Key)
	}
	for _, server := range root.Nodes {
		if server.Dir != true {
			return errorf(errWrongEtcdPath, server.Key)
		}
		name := strings.Split(server.Key, "/")[EtcdServersPrefixCnt+1]
		servers[name] = make(map[string][]string)
		for _, hosttype := range server.Nodes {
			if hosttype.Dir != true {
				return errorf(errWrongEtcdPath, hosttype.Key)
			}
			host := strings.Split(hosttype.Key, "/")[EtcdServersPrefixCnt+2]
			servers[name][host] = []string{}
			for _, hostaddr := range hosttype.Nodes {
				addr := strings.Split(hostaddr.Key, "/")[EtcdServersPrefixCnt+3]
				servers[name][host] = append(servers[name][host], addr)
			}
		}
	}

	mgr.mapServers = servers

	Log.Infof("UpdateServerHosts is done: %v", mgr.mapServers)
	return nil

}

// get host ips for the server, now return all hosts
func (mgr *ServerManager) GetServerHosts(serverName string, hostType string) ([]string, error) {
	server, ok := mgr.mapServers[serverName]
	if !ok {
		// try update server hosts mannually
		mgr.UpdateServerHosts()
	}
	server, ok = mgr.mapServers[serverName]
	if !ok {
		return nil, errorf("no server for %s", serverName)
	}
	hosts, ok := server[hostType]
	if !ok || len(hosts) == 0 {
		return nil, errorf("no hosts for %s:%s", serverName, hostType)
	}

	return hosts, nil
}
