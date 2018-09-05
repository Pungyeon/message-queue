package queue

import (
	"encoding/json"
	"testing"
)

func TestQueueServerInitialisation(t *testing.T) {
	q := NewServer("0.0.0.0", 1234)
	go func() {
		err := q.Start()
		if err != nil {
			t.Error(err)
		}
	}()
	q.Close()

	if q.running != false {
		t.Error("did not stop queue successfully")
	}
}

func TestQueueClientInitialisation(t *testing.T) {
	c := NewClient("localhost", 1235)
	if c.connectionString() != "localhost:1235" {
		t.Error("expected: localhost:1235, actual: " + c.connectionString())
	}
}

func TestClientServerConnection(t *testing.T) {
	q := NewServer("0.0.0.0", 1236)
	go func() {
		err := q.Start()
		if err != nil {
			t.Error(err)
		}
	}()
	defer q.Close()
	c := NewClient("localhost", 1236)
	_, err := c.connect()
	if err != nil {
		t.Error(err)
	}
}

func TestSendOperation(t *testing.T) {
	q := NewServer("0.0.0.0", 1238)
	go func() {
		err := q.Start()
		if err != nil {
			t.Error(err)
		}
	}()
	defer q.Close()

	c := NewClient("localhost", 1238)
	for i := 0; i < 10; i++ {
		msg, err := c.Push("hello world.")
		if err != nil {
			t.Error(err)
		}

		operation, err := NewOperationFromQueueBytes(msg)
		if err != nil {
			t.Error(err)
		}

		if operation.Value != "ACK" {
			t.Error("something went wrong with the server response")
		}
	}
}

func TestSubscribeOperation(t *testing.T) {
	q := NewServer("0.0.0.0", 1003)
	go func() {
		err := q.Start()
		if err != nil {
			t.Error(err)
		}
	}()
	defer q.Close()

	c := NewClient("localhost", 1003)
	err := c.Subscribe()
	if err != nil {
		t.Error(err)
	}
}

func TestSendAndReceive(t *testing.T) {
	q := NewServer("0.0.0.0", 1004)
	go func() {
		err := q.Start()
		if err != nil {
			t.Error(err)
		}
	}()

	defer q.Close()

	sender := NewClient("localhost", 1004)
	subscriber := NewClient("localhost", 1004)

	sender.Push("this is the first line")
	sender.Push("this is the second line")
	sender.Push("this is the third line")

	subscriber.Subscribe()

	sender.Push("this is the fourth line")
}

func TestOperationToString(t *testing.T) {
	op := Operation{
		Op:    Register,
		Value: "QueueName",
	}

	if op.String() != `{"operation":"register","value":"QueueName"}` {
		t.Error("string convertion deemed unacceptable result:")
		t.Error(op.String())
	}
}

func TestOperationToBytes(t *testing.T) {
	op := Operation{
		Op:    Register,
		Value: "QueueName",
	}

	d, err := json.Marshal(op)
	if err != nil {
		t.Error(err)
	}
	if string(d) != string(op.Bytes()) {
		t.Error("byte conversation does not yield same results as json.Marshal")
		t.Error(string(d))
		t.Error(string(op.Bytes()))
	}
}

func TestQueueMessagePush(t *testing.T) {
	q := NewServer("0.0.0.0", 1001)
	q.push("dingeling ding dong!")

	if len(q.messages) != 1 {
		t.Error("message not added to the message queue")
	}
}

func TestQueueMessagePop(t *testing.T) {
	q := NewServer("0.0.0.0", 1002)
	q.push("dingeling ding dong!")
	q.push("dingeling ding dong!")
	q.pop()

	if len(q.messages) != 1 {
		t.Error("message not added to the message queue")
	}
}
