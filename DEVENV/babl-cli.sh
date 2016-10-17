#!/bin/sh

echo Hi There! | $GOPATH/src/github.com/larskluge/babl/babl -c s1.babl.sh:4445 larskluge/string-upcase
