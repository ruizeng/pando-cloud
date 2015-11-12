package online

import (
	"reflect"
	"testing"
	"time"
)

var testid = uint64(100)

func checkOnlineStatus(t *testing.T, mgr *Manager, status Status) {
	readstatus, err := mgr.GetStatus(testid)
	if err != nil {
		t.Error(err)
	}

	if readstatus == nil {
		t.Errorf("device should be online, but is offline.")
	}

	if !reflect.DeepEqual(status, *readstatus) {
		t.Errorf("get status test error, want %v, got %v", status, *readstatus)
	}
}

func checkOfflineStatus(t *testing.T, mgr *Manager) {
	readstatus, err := mgr.GetStatus(testid)
	if err != nil {
		t.Error(err)
	}

	if readstatus != nil {
		t.Errorf("device should be offline, but got status: %v", readstatus)
	}

}

func TestManager(t *testing.T) {
	mgr := NewManager("localhost:6379")

	status := Status{
		ClientIP:          "3.3.3.3",
		AccessServerIP:    "192.168.9.1",
		HeartbeatInterval: 2,
	}

	err := mgr.GetOnline(testid, status)
	if err != nil {
		t.Error(err)
	}

	checkOnlineStatus(t, mgr, status)

	cnt := 0
	for {
		time.Sleep(time.Second * 2)
		if cnt > 2 {
			break
		}
		err := mgr.SetHeartbeat(testid)
		if err != nil {
			t.Error(err)
		}
		cnt++
	}

	checkOnlineStatus(t, mgr, status)

	err = mgr.GetOffline(testid)
	if err != nil {
		t.Error(err)
	}

	checkOfflineStatus(t, mgr)

}
