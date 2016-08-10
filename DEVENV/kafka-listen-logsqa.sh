#!/bin/sh

$GOPATH/src/github.com/Shopify/sarama/tools/kafka-console-consumer/kafka-console-consumer \
-verbose -brokers queue.babl.sh:9092 -offset oldest -topic logs.qa
