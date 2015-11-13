// heartbeat.go handles devices' heartbeat

package main

import (
	log "code.google.com/p/log4go"
	"github.com/garyburd/redigo/redis"
	"pandocloud.com/common/server"
	"strconv"
	"time"
)

// Handle a heartbeat from device
func DeviceHeartBeatFunc(args rpc_protocol.ArgsDeviceHeartBeat, reply *rpc_protocol.ReplyDeviceHeartBeat) error {

	id := args.Id
	id_str := strconv.FormatUint(id, 10)
	redis_conn, err := server.GetRedisClient()
	if err != nil {
		return ErrOnlineServer(Errno_Redis_Not_Initialized)
	}
	key := redisDeviceHeartBeatTimestamp + id_str
	// if this is the device's first heartbeat, then send sync_time_command immediately
	bytes_buffer_str, err := redis.String(redis_conn.Do("GET", key))
	//if len(bytes_buffer_str) < 1 {
	if err == redis.ErrNil { // means we don't have the device's online_status in redis before, so this is the device's first heartbeat
		SendSyncTimeCommand("DeviceHeartBeat", id)
	} else {
		devstatus, err := String2DeviceStatus(bytes_buffer_str)
		if err != nil {
			log.Error("[DeviceHeartBeat]gob decode error:%v", err)
			//no need to return error here
		} else {
			// the device is offline before, but now has heartbeat again, means it get online now
			// need sync time
			if devstatus.Online_status != STATUS_ONLINE {
				SendSyncTimeCommand("DeviceHeartBeat", id)
			}
		}
	}

	// update the device's timestamp in redis,and online_status is STATUS_ONLINE of course
	timestamp := uint64(time.Now().Unix()) //int64
	devstatus := DeviceStatus{id, timestamp, STATUS_ONLINE, args.ClientIP, args.AccessServer_InnerIP, args.AccessServer_OutterIP, args.HeartbeatTimeGap}
	// gob serialization
	bytes_buffer_str, err = DeviceStatus2String(devstatus)
	if err != nil {
		log.Error("[DeviceHeartBeat]gob encode %v error:%v", devstatus, err)
		reply.Result = -1
		return ErrOnlineServer(Errno_Redis_Error)
	}
	// store in redis as a string
	_, err = redis_conn.Do("SET", key, bytes_buffer_str)
	if err != nil {
		reply.Result = -1
		log.Error("[DeviceHeartBeat]redis set error:%v", err)
		return ErrOnlineServer(Errno_Redis_Error)
	}
	// may be we can add EXPIRE in redis, or we can check expire in GetDeviceStatus
	reply.Result = 0
	return nil
}
