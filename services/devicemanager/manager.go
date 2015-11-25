package main

import (
	"github.com/PandoCloud/pando-cloud/pkg/online"
	"github.com/PandoCloud/pando-cloud/pkg/rpcs"
	"github.com/PandoCloud/pando-cloud/pkg/token"
)

type DeviceManager struct {
	onlineManager *online.Manager
	tokenHelper   *token.Helper
}

func NewDeviceManager(redishost string) *DeviceManager {
	mgr := online.NewManager(redishost)

	helper := token.NewHelper(redishost)

	return &DeviceManager{
		onlineManager: mgr,
		tokenHelper:   helper,
	}
}

func (dm *DeviceManager) GenerateDeviceAccessToken(args rpcs.ArgsGenerateDeviceAccessToken, reply *rpcs.ReplyGenerateDeviceAccessToken) error {
	token, err := dm.tokenHelper.GenerateToken(args.Id)
	if err != nil {
		return err
	}

	reply.AccessToken = token
	return nil
}

func (dm *DeviceManager) ValidateDeviceAccessToken(args rpcs.ArgsValidateDeviceAccessToken, reply *rpcs.ReplyValidateDeviceAccessToken) error {
	dm.onlineManager.SetHeartbeat(args.Id)
	return dm.tokenHelper.ValidateToken(args.Id, args.AccessToken)
}

func (dm *DeviceManager) GetOnline(args rpcs.ArgsGetOnline, reply *rpcs.ReplyGetOnline) error {
	return dm.onlineManager.GetOnline(args.Id, online.Status{
		ClientIP:          args.ClientIP,
		AccessRPCHost:     args.AccessRPCHost,
		HeartbeatInterval: args.HeartbeatInterval,
	})
}

func (dm *DeviceManager) HeartBeat(args rpcs.ArgsHeartBeat, reply *rpcs.ReplyHeartBeat) error {
	return dm.onlineManager.SetHeartbeat(args.Id)
}

func (dm *DeviceManager) GetOffline(args rpcs.ArgsGetOffline, reply *rpcs.ReplyGetOffline) error {
	return dm.onlineManager.GetOffline(args.Id)
}

func (dm *DeviceManager) GetDeviceOnlineStatus(args rpcs.ArgsGetDeviceOnlineStatus, reply *rpcs.ReplyGetDeviceOnlineStatus) error {
	status, err := dm.onlineManager.GetStatus(args.Id)
	if err != nil {
		return err
	}

	reply.ClientIP = status.ClientIP
	reply.AccessRPCHost = status.AccessRPCHost
	reply.HeartbeatInterval = status.HeartbeatInterval
	return nil
}
