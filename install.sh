#!/bin/bash
echo "installing nats client..."
go get github.com/nats-io/nats
echo "installing nats server..."
go get github.com/nats-io/gnatsd
echo "installing linked list library..."
go get github.com/emirpasic/gods/lists/doublylinkedlist
