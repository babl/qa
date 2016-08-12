#!/bin/sh

docker stop logstash
docker rm logstash

docker pull registry.babl.sh/logstash:logmatic-v2

docker run -d --name logstash \
  -p 12300:12300 \
  -p 12500:12500 \
  -v "$PWD/logstash.conf:/etc/logstash.conf" \
  registry.babl.sh/logstash:logmatic-v2
