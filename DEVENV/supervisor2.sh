#!/bin/sh

$GOPATH/src/github.com/larskluge/supervisor2/supervisor2 -l=:4445 -kb=queue.babl.sh:9092 2>&1 | nc 127.0.0.1 12300
