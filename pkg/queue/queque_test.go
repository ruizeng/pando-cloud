package queue

import (
	// "fmt"
	"reflect"
	"testing"
)

type test struct {
	Cmd int
	Msg string
}

const testTopic = "test/topic"

var testChan chan test = make(chan test)

func recv(t *testing.T) {
	q := New("localhost:6379")
	msg := test{}
	err := q.Receive(testTopic, &msg)
	if err != nil {
		t.Error(err)
	}
	testChan <- msg
}

func TestQueue(t *testing.T) {
	testMessage := test{123, "hello"}

	q := New("localhost:6379")

	go recv(t)

	err := q.Send(testTopic, testMessage)
	if err != nil {
		t.Error(err)
	}

	recvMessage := <-testChan

	if !reflect.DeepEqual(testMessage, recvMessage) {
		t.Errorf("receive message not match, want: %v, got : %v", testMessage, recvMessage)
	}
}
