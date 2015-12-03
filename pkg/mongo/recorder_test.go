package mongo

import (
	"github.com/PandoCloud/pando-cloud/pkg/protocol"
	"github.com/PandoCloud/pando-cloud/pkg/rpcs"
	"github.com/PandoCloud/pando-cloud/pkg/tlv"
	"reflect"
	"testing"
	"time"
)

func TestRecorder(t *testing.T) {
	r, err := NewRecorder("localhost", "pandocloud", "commands")
	if err != nil {
		t.Fatal(err)
	}

	tlvs, err := tlv.MakeTLVs([]interface{}{float64(0.1), int64(100), uint64(200)})
	if err != nil {
		t.Error(err)
	}

	deviceid := uint64(12345)
	timestamp := uint64(time.Now().Unix() * 1000)

	subdata := protocol.SubData{
		Head:   protocol.SubDataHead{1, 2, 3},
		Params: tlvs,
	}

	subdatas := []protocol.SubData{}

	subdatas = append(subdatas, subdata)

	data := rpcs.ArgsPutData{
		DeviceId:  deviceid,
		Timestamp: timestamp,
		Subdata:   subdatas,
	}

	err = r.Insert(data)
	if err != nil {
		t.Error(err)
	}

	readData := rpcs.ArgsPutData{}
	err = r.FindLatest(deviceid, &readData)

	if err != nil {
		t.Error(err)
	}

	if !reflect.DeepEqual(data, readData) {
		t.Errorf("read data want %v, but got %v", data, readData)
	}

	readDatas := []rpcs.ArgsPutData{}
	err = r.FindByTimestamp(deviceid, timestamp, timestamp, &readDatas)
	t.Log(readDatas)
	if !reflect.DeepEqual(data, readDatas[0]) {
		t.Errorf("read data by timestamp want %v, but got %v", data, readDatas[0])
	}
}
