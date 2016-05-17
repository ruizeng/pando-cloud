// package online manage device online state and store it in redis.
package online

import (
	"errors"
	"github.com/PandoCloud/pando-cloud/pkg/redispool"
	"github.com/PandoCloud/pando-cloud/pkg/serializer"
	"github.com/garyburd/redigo/redis"
	"strconv"
)

const (
	OnlineStatusKeyPrefix = "device:onlinestatus:"
)

type Status struct {
	ClientIP          string
	AccessRPCHost     string
	HeartbeatInterval uint32
}

type Manager struct {
	redishost string
}

func NewManager(host string) *Manager {
	mgr := &Manager{
		redishost: host,
	}
	return mgr
}

func (mgr *Manager) GetStatus(id uint64) (*Status, error) {
	key := OnlineStatusKeyPrefix + strconv.FormatUint(id, 10)
	conn, err := redispool.GetClient(mgr.redishost)
	if err != nil {
		return nil, err
	}

	status := &Status{}
	// get status from redis
	bufferStr, err := redis.String(conn.Do("GET", key))
	if err != nil {
		return nil, err
	}
	err = serializer.String2Struct(bufferStr, status)
	if err != nil {
		return nil, err
	}

	return status, nil
}

func (mgr *Manager) GetOnline(id uint64, status Status) error {
	key := OnlineStatusKeyPrefix + strconv.FormatUint(id, 10)
	conn, err := redispool.GetClient(mgr.redishost)
	if err != nil {
		return err
	}
	// serialize and store the device's online status info in redis
	bufferStr, err := serializer.Struct2String(status)
	if err != nil {
		return err
	}
	_, err = conn.Do("SET", key, bufferStr)
	if err != nil {
		return err
	}
	_, err = conn.Do("EXPIRE", key, status.HeartbeatInterval+status.HeartbeatInterval/2)
	if err != nil {
		return err
	}

	return nil
}

func (mgr *Manager) SetHeartbeat(id uint64) error {
	status, err := mgr.GetStatus(id)
	if err != nil {
		return err
	}

	if status == nil {
		return errors.New("device offline.")
	}

	key := OnlineStatusKeyPrefix + strconv.FormatUint(id, 10)
	conn, err := redispool.GetClient(mgr.redishost)
	if err != nil {
		return err
	}

	_, err = conn.Do("EXPIRE", key, status.HeartbeatInterval+status.HeartbeatInterval/2)
	if err != nil {
		return err
	}

	return nil
}

func (mgr *Manager) GetOffline(id uint64) error {
	key := OnlineStatusKeyPrefix + strconv.FormatUint(id, 10)
	conn, err := redispool.GetClient(mgr.redishost)
	if err != nil {
		return err
	}

	_, err = conn.Do("DEL", key)
	if err != nil {
		return err
	}

	return nil
}
