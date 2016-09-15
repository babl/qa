#!/bin/sh

$GOPATH/src/github.com/larskluge/babl-qa/babl-qa -l=:8888 -kb=queue.babl.sh:9092
#$GOPATH/src/github.com/larskluge/babl-qa/babl-qa -l=:8888 -kb=queue.babl.sh:9092 2>&1 | jq .
