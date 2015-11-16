// package queue implement a message queque api with redis pub/sub command
package queue

import (
	// "fmt"
	"github.com/PandoCloud/pando-cloud/pkg/redispool"
	"github.com/PandoCloud/pando-cloud/pkg/serializer"
	// "github.com/garyburd/redigo/redis"
)

type Queue struct {
	redishost string
}

func New(redishost string) *Queue {
	return &Queue{redishost}
}

func (q *Queue) Send(topic string, msg interface{}) error {
	conn, err := redispool.GetClient(q.redishost)
	if err != nil {
		return err
	}

	msgStr, err := serializer.Struct2String(msg)
	if err != nil {
		return err
	}

	_, err = conn.Do("PUBLISH", topic, msgStr)
	if err != nil {
		return err
	}

	return nil
}

func (q *Queue) Receive(topic string, target interface{}) error {
	conn, err := redispool.GetClient(q.redishost)
	if err != nil {
		return err
	}
	_, err = conn.Do("SUBSCRIBE", topic)
	if err != nil {
		return err
	}

	result, err := conn.Receive()
	if err != nil {
		return err
	}

	strMsg := string(result.([]interface{})[2].([]byte))

	err = serializer.String2Struct(strMsg, target)
	if err != nil {
		return err
	}

	return nil

}
