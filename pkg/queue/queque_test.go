package queue

import (
	"reflect"
	"testing"
	"time"
)

type test struct {
	Cmd int
	Msg string
}

const testQueueName = "test/queue/somename"

var testChan chan test = make(chan test)

func recv(t *testing.T) {
	q, err := New("amqp://guest:guest@localhost:5672/", testQueueName)
	if err != nil {
		t.Error(err)
	}
	msg := test{}
	err = q.Receive(&msg)
	if err != nil {
		t.Error(err)
	}
	testChan <- msg
}

func TestQueue(t *testing.T) {
	testMessage := test{123, "hello"}

	q, err := New("amqp://guest:guest@localhost:5672/", testQueueName)
	if err != nil {
		t.Fatal(err)
	}

	go recv(t)

	time.Sleep(time.Second)

	err = q.Send(testMessage)
	if err != nil {
		t.Fatal(err)
	}

	recvMessage := <-testChan

	if !reflect.DeepEqual(testMessage, recvMessage) {
		t.Errorf("receive message not match, want: %v, got : %v", testMessage, recvMessage)
	}
}
