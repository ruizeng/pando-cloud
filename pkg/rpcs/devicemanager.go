package rpcs

import (
	"github.com/PandoCloud/pando-cloud/pkg/online"
)

type ArgsGenerateDeviceAccessToken ArgsDeviceId
type ReplyGenerateDeviceAccessToken struct {
	AccessToken []byte
}

type ArgsValidateDeviceAccessToken struct {
	Id          uint64
	AccessToken []byte
}
type ReplyValidateDeviceAccessToken ReplyEmptyResult

type ArgsGetOnline struct {
	Id                uint64
	ClientIP          string
	AccessRPCHost     string
	HeartbeatInterval uint32
}
type ReplyGetOnline ReplyEmptyResult

type ArgsGetOffline ArgsDeviceId
type ReplyGetOffline ReplyEmptyResult

type ArgsHeartBeat struct {
	Id uint64
}
type ReplyHeartBeat ReplyEmptyResult

type ArgsGetDeviceStatus ArgsDeviceId
type ReplyGetDeviceStatus online.Status
