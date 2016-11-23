#!/bin/sh

echo Hi There! | $GOPATH/src/github.com/larskluge/babl/babl -c sandbox.babl.sh:4445 larskluge/string-upcase
