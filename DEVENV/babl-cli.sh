#!/bin/sh

echo Hi There! | $GOPATH/src/github.com/larskluge/babl/babl -c localhost:4445 larskluge/string-upcase
