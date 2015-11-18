package main

import (
	"github.com/PandoCloud/pando-cloud/pkg/rpcs"
	"testing"
)

func TestDeviceManager(t *testing.T) {
	mgr := NewDeviceManager("localhost:6379")

	deviceid := uint64(123456)

	args1 := rpcs.ArgsGenerateDeviceAccessToken{
		Id: deviceid,
	}
	reply1 := rpcs.ReplyGenerateDeviceAccessToken{}
	err := mgr.GenerateDeviceAccessToken(args1, &reply1)
	if err != nil {
		t.Fatal(err)
	}

	token := reply1.AccessToken

	args2 := rpcs.ArgsValidateDeviceAccessToken{
		Id:          deviceid,
		AccessToken: token,
	}
	reply2 := rpcs.ReplyValidateDeviceAccessToken{}
	err = mgr.ValidateDeviceAccessToken(args2, &reply2)
	if err != nil {
		t.Fatal(err)
	}

	args3 := rpcs.ArgsGetOnline{
		Id:                deviceid,
		ClientIP:          "",
		AccessRPCHost:     "",
		HeartbeatInterval: 10,
	}
	reply3 := rpcs.ReplyGetOnline{}
	err = mgr.GetOnline(args3, &reply3)
	if err != nil {
		t.Fatal(err)
	}

	args4 := rpcs.ArgsHeartBeat{
		Id: deviceid,
	}
	reply4 := rpcs.ReplyHeartBeat{}
	err = mgr.HeartBeat(args4, &reply4)
	if err != nil {
		t.Fatal(err)
	}

	args5 := rpcs.ArgsGetDeviceStatus{
		Id: deviceid,
	}
	reply5 := rpcs.ReplyGetDeviceStatus{}
	err = mgr.GetDeviceStatus(args5, &reply5)
	if err != nil {
		t.Fatal(err)
	}
	t.Log(reply5)

	args6 := rpcs.ArgsGetOffline{
		Id: deviceid,
	}
	reply6 := rpcs.ReplyGetOffline{}
	err = mgr.GetOffline(args6, &reply6)
	if err != nil {
		t.Fatal(err)
	}

}
