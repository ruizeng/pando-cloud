package mqtt

import (
        "testing"
	"reflect"
)

func TestQosLevel(t *testing.T) {
	qos := QosAtMostOnce

	if !qos.IsValid() {
		t.Errorf("qos valid judge is wrong\n")
	}

	if qos.HasId() {
		t.Errorf("qos hasID judge is wrong\n")
	}
}

func TestConnectMsg(t *testing.T) {
	rw := &TestReadWriter{
                w:   0,
                r:   0,
                buf: make([]byte, 10000),
        }


	head := &Header{
		DupFlag: false,
		QosLevel: QosAtMostOnce,
		Retain: true,
	}
	connMsg := &Connect{
		Header: *head,	
		ProtocolName: "mqtt",
		ProtocolVersion: 1,
		WillRetain: false,
		WillFlag: true,
		ClientId: "itachili",
	}

	err := connMsg.Encode(rw)	
	if err != nil {
		t.Error(err)
	}

	newHead := &Header{}
	newConnMsg := &Connect{}
	_, _, err = newHead.Decode(rw)
	if err != nil {
		t.Error(err)
	}
	packetRemaining := int32(rw.Size())
	err = newConnMsg.Decode(rw, *newHead, packetRemaining)	
	if err != nil {
		t.Error(err)
	}
	if !reflect.DeepEqual(connMsg, newConnMsg) {
		t.Errorf("the origin:\n%x\n, now:\n%x\n", connMsg, newConnMsg)
	}	
}

func TestConnAck(t *testing.T) {
	rw := &TestReadWriter{
                w:   0,
                r:   0,
                buf: make([]byte, 10000),
        }


	head := &Header{
		DupFlag: false,
		QosLevel: QosAtMostOnce,
		Retain: true,
	}
	conAck := &ConnAck{
		Header: *head,	
		ReturnCode: RetCodeAccepted,
	}

	err := conAck.Encode(rw)	
	if err != nil {
		t.Error(err)
	}

	newHead := &Header{}
	newConnAck := &ConnAck{}
	_, _, err = newHead.Decode(rw)
	if err != nil {
		t.Error(err)
	}
	packetRemaining := int32(rw.Size())
	err = newConnAck.Decode(rw, *newHead, packetRemaining)	
	if err != nil {
		t.Error(err)
	}
	if !reflect.DeepEqual(conAck, newConnAck) {
		t.Errorf("the origin:\n%x\n, now:\n%x\n", conAck, newConnAck)
	}
}

func TestPublish(t *testing.T) { 
	rw := &TestReadWriter{
                w:   0,
                r:   0,
                buf: make([]byte, 10000),
        }


	head := &Header{
		DupFlag: false,
		QosLevel: QosAtLeastOnce,
		Retain: true,
	}
	conAck := &Publish{
		Header: *head,	
		TopicName: "itachili",
		MessageId: 18,
		Payload: BytesPayload([]byte{'1', '2', '3'}),
	}

	err := conAck.Encode(rw)	
	if err != nil {
		t.Error(err)
	}

	newHead := &Header{}
	newConnAck := &Publish{}
	_, _, err = newHead.Decode(rw)
	if err != nil {
		t.Error(err)
	}
	packetRemaining := int32(rw.Size())
	err = newConnAck.Decode(rw, *newHead, packetRemaining)	
	if err != nil {
		t.Error(err)
	}
	if !reflect.DeepEqual(conAck, newConnAck) {
		t.Errorf("the origin:\n%x\n, now:\n%x\n", conAck, newConnAck)
	}
}

func TestPubAck(t *testing.T) { 
	rw := &TestReadWriter{
                w:   0,
                r:   0,
                buf: make([]byte, 10000),
        }


	head := &Header{
		DupFlag: false,
		QosLevel: QosAtLeastOnce,
		Retain: true,
	}
	conAck := &PubAck{
		Header: *head,	
		MessageId: 18,
	}

	err := conAck.Encode(rw)	
	if err != nil {
		t.Error(err)
	}

	newHead := &Header{}
	newConnAck := &PubAck{}
	_, _, err = newHead.Decode(rw)
	if err != nil {
		t.Error(err)
	}
	packetRemaining := int32(rw.Size())
	err = newConnAck.Decode(rw, *newHead, packetRemaining)	
	if err != nil {
		t.Error(err)
	}
	if !reflect.DeepEqual(conAck, newConnAck) {
		t.Errorf("the origin:\n%x\n, now:\n%x\n", conAck, newConnAck)
	}
}

func TestPubRel(t *testing.T) { 
	rw := &TestReadWriter{
                w:   0,
                r:   0,
                buf: make([]byte, 10000),
        }


	head := &Header{
		DupFlag: false,
		QosLevel: QosAtLeastOnce,
		Retain: true,
	}
	conAck := &PubRel{
		Header: *head,	
		MessageId: 18,
	}

	err := conAck.Encode(rw)	
	if err != nil {
		t.Error(err)
	}

	newHead := &Header{}
	newConnAck := &PubRel{}
	_, _, err = newHead.Decode(rw)
	if err != nil {
		t.Error(err)
	}
	packetRemaining := int32(rw.Size())
	err = newConnAck.Decode(rw, *newHead, packetRemaining)	
	if err != nil {
		t.Error(err)
	}
	if !reflect.DeepEqual(conAck, newConnAck) {
		t.Errorf("the origin:\n%x\n, now:\n%x\n", conAck, newConnAck)
	}
}

func TestPubRec(t *testing.T) { 
	rw := &TestReadWriter{
                w:   0,
                r:   0,
                buf: make([]byte, 10000),
        }


	head := &Header{
		DupFlag: false,
		QosLevel: QosAtLeastOnce,
		Retain: true,
	}
	conAck := &PubRec{
		Header: *head,	
		MessageId: 18,
	}

	err := conAck.Encode(rw)	
	if err != nil {
		t.Error(err)
	}

	newHead := &Header{}
	newConnAck := &PubRec{}
	_, _, err = newHead.Decode(rw)
	if err != nil {
		t.Error(err)
	}
	packetRemaining := int32(rw.Size())
	err = newConnAck.Decode(rw, *newHead, packetRemaining)	
	if err != nil {
		t.Error(err)
	}
	if !reflect.DeepEqual(conAck, newConnAck) {
		t.Errorf("the origin:\n%x\n, now:\n%x\n", conAck, newConnAck)
	}
}

func TestPubComp(t *testing.T) { 
	rw := &TestReadWriter{
                w:   0,
                r:   0,
                buf: make([]byte, 10000),
        }


	head := &Header{
		DupFlag: false,
		QosLevel: QosAtLeastOnce,
		Retain: true,
	}
	conAck := &PubComp{
		Header: *head,	
		MessageId: 18,
	}

	err := conAck.Encode(rw)	
	if err != nil {
		t.Error(err)
	}

	newHead := &Header{}
	newConnAck := &PubComp{}
	_, _, err = newHead.Decode(rw)
	if err != nil {
		t.Error(err)
	}
	packetRemaining := int32(rw.Size())
	err = newConnAck.Decode(rw, *newHead, packetRemaining)	
	if err != nil {
		t.Error(err)
	}
	if !reflect.DeepEqual(conAck, newConnAck) {
		t.Errorf("the origin:\n%x\n, now:\n%x\n", conAck, newConnAck)
	}
}
