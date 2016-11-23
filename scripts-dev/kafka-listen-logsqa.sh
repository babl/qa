#!/bin/sh

#$GOPATH/src/github.com/Shopify/sarama/tools/kafka-console-consumer/kafka-console-consumer \
#-verbose -brokers sandbox.babl.sh:9092 -offset oldest -topic logs.qa

$GOPATH/src/github.com/Shopify/sarama/tools/kafka-console-consumer/kafka-console-consumer \
-verbose -brokers sandbox.babl.sh:9092 -offset newest -topic logs.qa
