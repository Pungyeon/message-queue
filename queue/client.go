package queue

import (
	"bufio"
	"fmt"
	"log"
	"net"
	"strconv"
)

// Client is a struct for interacting with the queue messaging service
type Client struct {
	host         string
	port         int
	subscription net.Conn
}

// NewClient will initialise a new message queue client
func NewClient(host string, port int) *Client {
	return &Client{
		host: host,
		port: port,
	}
}

func (c *Client) connectionString() string {
	return c.host + ":" + strconv.Itoa(c.port)
}

func (c *Client) connect() (net.Conn, error) {
	return net.Dial("tcp", c.connectionString())
}

// Push will send a message to the message queue
func (c *Client) Push(msg string) ([]byte, error) {
	conn, err := c.connect()
	if err != nil {
		return nil, err
	}
	defer conn.Close()
	operation := Operation{
		Op:    Message,
		Value: msg,
	}
	conn.Write(operation.QueueBytes())
	return bufio.NewReader(conn).ReadBytes('\000')
}

// Subscribe will subscribe to a message queue
func (c *Client) Subscribe() error {
	var err error
	c.subscription, err = c.connect()
	if err != nil {
		return err
	}
	op := Operation{Op: Subscribe}
	_, err = c.subscription.Write(op.QueueBytes())
	if err != nil {
		return err
	}
	go c.startRetrieveCycle()
	return nil
}

// UnSubscribe closes the connection to the message queu
func (c *Client) UnSubscribe() {
	c.subscription.Close()
}

func (c *Client) startRetrieveCycle() {
	var operation Operation
	for {
		data, err := bufio.NewReader(c.subscription).ReadBytes('\000')
		if err != nil {
			log.Println(err)
			continue
		}
		c.subscription.Write([]byte{'\x00'})
		operation, err = NewOperationFromQueueBytes(data)
		if err != nil {
			log.Println(err)
			continue
		}
		fmt.Println(operation.Value)
	}
}
