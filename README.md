## Overview
This project consists of 2 components, the actual messaging queue and the client operators, which communicate with the messaging queue, retrieving and sending messages.

## Build
Building this project is relatively simple. There are no 3rd party dependencies as everything is written with the standard  go library. To build the binary file for this project, simply write the following command: 

> go build -o dist/message-queue main.go 

This will place an executable file in the `dist` folder. 

## Running the application
This project is always run with the main binary, using CLI parameters to determine the behaviour and role of the running software. All examples will use an already compiled binary with the filename `message-queue`. Please note, on windows systems, this file will be named `message-queue.exe`.

### Start Message Queue
To start the message queue listening on all address (0.0.0.0) on port `1234` use the following command:

> ./message-queue -qhost="0.0.0.0" -port=1234 -server

### Send a file to the Message Queue
To send a file with the name `textfile.txt` use the following command:

> ./message-queue -file="textfile.txt"

### Start a Subscriber
To start the service as a subscriber, use the following command, specifying the queue service to connect to:

> ./message-queue -qhost="localhost" -port=1234 -subscribe

This will attempt to connect to a queue service on `localhost:1234/tcp` and subscribe to any messages in the queue. Note, that adding the `-subscribe` parameter will always start a subscriber, so it is possible to send a file <b>and</b> subscribe on the same service, using the following command:

> ./message-queue -file="textfile.txt" -subscribe

### Help
For a list of all commands, please use the following command:

> ./message-queue -h

## Todo
- [ ] Implement unsubscribe on graceful shutdown
    - [ ] And implement unsubscribe on bad connection
- [ ] Implement multiple exchanges / routes
- [ ] Implement send retries
    - [ ] Implement send retries pr. client

## Contact
For inquiries that cannot be handled with an issue request, please contact lja@pungy.dk
