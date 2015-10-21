// config flags from command line or ini conf file.

package server

import (
	"flag"
)

const (
	flagTCPHost  = "tcphost"
	flagUseTls   = "usetls"
	flagHTTPHost = "httphost"
	flagUseHttps = "usehttps"
	flagCAFile   = "cafile"
	flagKeyFile  = "keyfile"
	flagRPCHost  = "rpchost"
	flagEtcd     = "etcd"
	flagLogLevel = "loglevel"
)

var (
	confTCPHost = flag.String(flagTCPHost, "", "tcp server listen address, format ip:port")
	confUseTls  = flag.Bool(flagUseTls, false, "if tcp server uses tls, default false")

	confHTTPHost = flag.String(flagHTTPHost, "", "http server listen address, format ip:port")
	confUseHttps = flag.Bool(flagUseHttps, false, "if http server uses tls, default false")

	confCAFile  = flag.String(flagCAFile, "cacert.pem", "public ca pem file path")
	confKeyFile = flag.String(flagKeyFile, "privkey.pem", "private key pem file path")

	confRPCHost = flag.String(flagRPCHost, "", "rpc server listen address, format ip:port")

	confEtcd = flag.String(flagEtcd, "", "etcd service addr, format ip:port;ip:port")

	confLogLevel = flag.String(flagLogLevel, "warn", "default log level, options are panic|fatal|error|warn|info|debug")
)
