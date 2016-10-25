#!/bin/sh

$GOPATH/src/github.com/larskluge/babl-server/babl-server -kb=queue.babl.sh:9092 -m larskluge/string-upcase 2>&1 | nc 127.0.0.1 12500
#$GOPATH/src/github.com/larskluge/babl-server/babl-server -kb=queue.babl.sh:9092 -m larskluge/string-upcase

# Simulate ERROR response
# $GOPATH/src/github.com/larskluge/babl-server/babl-server -kb=queue.babl.sh:9092 -m larskluge/string-upcase --cmd "/Users/nelson/work/GO/src/github.com/nneves/stderr/stderr" 2>&1 | nc 127.0.0.1 12500

# Simulate timeout
# $GOPATH/src/github.com/larskluge/babl-server/babl-server -kb=queue.babl.sh:9092 -m larskluge/string-upcase --cmd "/Users/nelson/work/GO/src/github.com/nneves/stdin-delay/stdin-delay" 2>&1 | nc 127.0.0.1 12500
