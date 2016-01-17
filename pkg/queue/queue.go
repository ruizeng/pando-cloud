// package queue implement a message queque api with rabbitmq
package queue

import (
	"errors"
	"github.com/PandoCloud/pando-cloud/pkg/serializer"
	"github.com/streadway/amqp"
)

const defaultRecvChanLen = 8

type Queue struct {
	rabbithost   string
	conn         *amqp.Connection
	ch           *amqp.Channel
	queue        amqp.Queue
	recvChan     chan ([]byte)
	beginReceive bool
}

func New(rabbithost string, name string) (*Queue, error) {
	conn, err := amqp.Dial(rabbithost)
	if err != nil {
		return nil, err
	}

	ch, err := conn.Channel()
	if err != nil {
		return nil, err
	}

	queue, err := ch.QueueDeclare(
		name,  // name
		true,  // durable
		false, // delete when unused
		false, // exclusive
		false, // no-wait
		nil,   // arguments
	)
	if err != nil {
		return nil, errors.New("Failed to declare a queue")
	}

	err = ch.Qos(
		1,     // prefetch count
		0,     // prefetch size
		false, // global
	)
	if err != nil {
		return nil, errors.New("Failed to set QoS")
	}

	q := &Queue{rabbithost, conn, ch, queue, nil, false}

	return q, nil
}

func (q *Queue) keepReceivingFromQueue() {
	if q.ch == nil || q.recvChan == nil {
		//Message Queue Not Initialzed.
		return
	}

	defer func() {
		if q.recvChan != nil {
			close(q.recvChan)
		}
	}()

	msgs, err := q.ch.Consume(
		q.queue.Name, // queue
		"",           // consumer
		false,        // auto-ack
		false,        // exclusive
		false,        // no-local
		false,        // no-wait
		nil,          // args
	)

	if err != nil {
		return
	}

	for d := range msgs {
		q.recvChan <- d.Body
		d.Ack(false)
	}

}

// Send will send a message to the queue.
func (q *Queue) Send(msg interface{}) error {
	if q.ch == nil {
		return errors.New("Message Queue Not Initialzed.")
	}
	msgStr, err := serializer.Struct2String(msg)
	if err != nil {
		return err
	}
	err = q.ch.Publish(
		"",           // exchange
		q.queue.Name, // routing key
		false,        // mandatory
		false,
		amqp.Publishing{
			DeliveryMode: amqp.Persistent,
			ContentType:  "text/plain",
			Body:         []byte(msgStr),
		})

	return nil
}

// Receive will reveive a message from the queue. may be blocked if there is no message in queue.
func (q *Queue) Receive(target interface{}) error {
	if !q.beginReceive {
		q.recvChan = make(chan ([]byte), defaultRecvChanLen)
		go q.keepReceivingFromQueue()
		q.beginReceive = true
	}

	if q.recvChan == nil {
		return errors.New("Message Queue Has Not Been Initialized.")
	}

	msg, ok := <-q.recvChan

	if !ok {
		return errors.New("Message Queue Has Been Closed.")
	}

	strMsg := string(msg)
	err := serializer.String2Struct(strMsg, target)
	if err != nil {
		return err
	}

	return nil

}
