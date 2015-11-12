// Define the errors
package main

import (
	"fmt"
)

type ErrDeviceManager int

const (
	ErrContextNotInitialized int = -1
	ErrRedisError            int = -2
	ErrGobError              int = -3
	ErrRedisNotInitialized   int = -4
)

func (e ErrDeviceManager) Error() string {
	switch int(e) {
	case ErrContextNotInitialized:
		return fmt.Sprintf("context not initialized, errno:%d", e)
	case ErrRedisError:
		return fmt.Sprintf("redis operation error, errno:%d", e)
	case ErrGobError:
		return fmt.Sprintf("redis not initialized, errno:%d", e)
	case ErrRedisNotInitialized:
		return fmt.Sprintf("gob encode/decode fails, errno:%d", e)
	default:
		return fmt.Sprintf("online_server unknow error, errno:%d", e)
	}
}
