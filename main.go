package main

import (
	"bufio"
	"flag"
	"fmt"
	"log"
	"os"
	"sync"

	"github.com/Pungyeon/message-queue/queue"
)

func main() {
	qhost := flag.String("qhost", "0.0.0.0", "[string] - set server listening address to specified IP Address (example: 10.10.10.10)")
	port := flag.Int("port", 1234, "[int] - set the server listening port")
	server := flag.Bool("server", false, "[boolean] - set to true for starting a server, this will override subscribe and file parameters")

	filename := flag.String("file", "", "[string] - specify a file to send to the message queue. (REMEMBER port and qhost)")
	subscribe := flag.Bool("subscribe", false, "[boolean] - set to true if you wish to start a subscription to a specified server. (REMEMBER port and qhost)")
	flag.Parse()

	wg := sync.WaitGroup{}

	if *server {
		q := queue.NewServer(*qhost, *port)
		q.Start()
		return
	}

	if *subscribe {
		sub := queue.NewClient(*qhost, *port)
		sub.Subscribe()
		wg.Add(1)
	}

	if *filename != "" {
		sender := queue.NewClient(*qhost, *port)
		fmt.Println("attempting to read: " + *filename)
		file, err := os.Open(*filename)
		if err != nil {
			log.Fatal(err)
		}
		defer file.Close()

		scanner := bufio.NewScanner(file)
		for scanner.Scan() {
			sender.Push(scanner.Text())
		}

		if err := scanner.Err(); err != nil {
			log.Println(err)
		}
	}

	wg.Wait()
}
