#!/bin/sh

$GOPATH/src/github.com/larskluge/babl-server/babl-server -kb=queue.babl.sh:9092 -m larskluge/string-upcase 2>&1 | nc 127.0.0.1 12500
