package queue

import (
	"bufio"
	"log"
	"net"
	"strconv"
)

// Queue is a struct representing a message queue
type Queue struct {
	host        string
	port        int
	running     bool
	connections chan net.Conn
	operations  chan Operation
	subscribers map[string]net.Conn
	messages    []string
}

// New initialises a new messaging Queue
func NewServer(host string, port int) *Queue {
	return &Queue{
		host:        host,
		port:        port,
		connections: make(chan net.Conn, 0),
		operations:  make(chan Operation, 0),
		subscribers: map[string]net.Conn{},
	}
}

// Start will start the queue messaging service
func (q *Queue) Start() error {
	go q.listenForMessages()

	listener, err := net.Listen("tcp", "0.0.0.0:"+strconv.Itoa(q.port))
	if err != nil {
		return err
	}
	q.running = true

	for q.running == true {
		q.acceptConnection(listener)
	}
	return nil
}

func (q *Queue) listenForMessages() {
	for {
		select {
		case conn := <-q.connections:
			log.Printf("Accepted connection from: %v\n", conn.RemoteAddr())
			if operation, err := q.handleConnection(conn); err != nil {
				log.Println(err)
			} else {
				q.handleOperation(operation, conn)
			}
		}
	}
}

func (q *Queue) handleConnection(conn net.Conn) (Operation, error) {
	b, err := bufio.NewReader(conn).ReadBytes('\000')
	if err != nil {
		return Operation{}, err
	}

	op, err := NewOperationFromQueueBytes(b)
	if err != nil {
		log.Println(err)
		return Operation{}, err
	}
	return op, nil
}

func (q *Queue) handleOperation(operation Operation, conn net.Conn) {
	switch operation.Op {
	case Message:
		q.push(operation.Value)
		log.Printf("[Queue Size: %d] New Message Received: %s", len(q.messages), operation.Value)
		if len(q.subscribers) > 0 {
			q.flushMessageQueueToSubscribers()
		}
	case Subscribe:
		q.subscribed(conn)
		log.Printf("[Queue Size: %d] New Subscriber: %s", len(q.messages), conn.RemoteAddr().String())
		if len(q.messages) > 0 {
			q.flushMessageQueueToSubscribers()
		}
	case Unsubscribe:
		q.unsubscribed(conn)
	}
}

func (q *Queue) flushMessageQueueToSubscribers() {
	for len(q.messages) > 0 {
		operation := Operation{Op: Message, Value: q.messages[0]}
		for _, sub := range q.subscribers {
			sub.Write(operation.QueueBytes())
			_, err := bufio.NewReader(sub).ReadBytes('\000')
			if err != nil {
				log.Println(err)
			}
		}
		q.pop()
	}
	log.Printf("[Queue Size: %d] Message Queue Flushed", len(q.messages))
}

func (q *Queue) acceptConnection(listener net.Listener) {
	conn, err := listener.Accept()
	if err != nil {
		log.Println(err)
		return
	}

	q.connections <- conn
	operation := Operation{Op: Notify, Value: "ACK"}
	conn.Write(operation.QueueBytes())
}

// Close shuts down the queue messaging server
func (q *Queue) Close() {
	q.running = false
}

func (q *Queue) push(msg string) {
	q.messages = append(q.messages, msg)
}

func (q *Queue) pop() {
	q.messages = q.messages[1:]
}

func (q *Queue) subscribed(conn net.Conn) {
	q.subscribers[conn.RemoteAddr().String()] = conn
}

func (q *Queue) unsubscribed(conn net.Conn) {
	delete(q.subscribers, conn.RemoteAddr().String())
}
