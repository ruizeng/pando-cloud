// error messages

package server

const (
	errServerNotInit           = "Server has not been initialized...You must call server.Init(name) first !"
	errTCPHandlerNotRegistered = "Start TCP Server error : tcp handler not registered !"
	errMissingFlag             = "Missing flag: %s !"
	errLoadSecureKey           = "Load secret key file failed - %s"
	errListenFailed            = "FATAL: tcp listen (%s) failed - %s"
	errNewConnection           = "receive new connection error (%s)"
	errWrongHostAddr           = "wrong address : %s"
	errWrongEtcdPath           = "wrong path in etcd: %s"
)
