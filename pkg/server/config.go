// config flags from command line or ini conf file.

package server

import (
	"flag"
)

const (
	FlagTCPHost  = "tcphost"
	FlagUseTls   = "usetls"
	FlagHTTPHost = "httphost"
	FlagUseHttps = "usehttps"
	FlagCAFile   = "cafile"
	FlagKeyFile  = "keyfile"
	FlagRPCHost  = "rpchost"
	FlagEtcd     = "etcd"
	FlagLogLevel = "loglevel"
)

var (
	confTCPHost = flag.String(FlagTCPHost, "", "tcp server listen address, format ip:port")
	confUseTls  = flag.Bool(FlagUseTls, false, "if tcp server uses tls, default false")

	confHTTPHost = flag.String(FlagHTTPHost, "", "http server listen address, format ip:port")
	confUseHttps = flag.Bool(FlagUseHttps, false, "if http server uses tls, default false")

	confCAFile  = flag.String(FlagCAFile, "cacert.pem", "public ca pem file path")
	confKeyFile = flag.String(FlagKeyFile, "privkey.pem", "private key pem file path")

	confRPCHost = flag.String(FlagRPCHost, "", "rpc server listen address, format ip:port")

	confEtcd = flag.String(FlagEtcd, "", "etcd service addr, format ip:port;ip:port")

	confLogLevel = flag.String(FlagLogLevel, "info", "default log level, options are panic|fatal|error|warn|info|debug")
)
