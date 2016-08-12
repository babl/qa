# Babl Quality Assurance Service

Initial draft: [https://docs.google.com/babl-qa](https://docs.google.com/presentation/d/1-j8LmxzHOYQMkUNUKSrIgwG9MN3MRfkyG4z3P_u8pf0/edit?usp=sharing)

### How to test babl-qa in the local dev-env
In order to test `babl-qa` locally it is required to launch 4 additional services/modules:

- logstash
- supervisor2
- babl-server
- babl (cli)

`babl` + `supervisor2` + `babl-server` will use [queue.babl.sh:9092](queue.babl.sh:9092) kafka pub-sub service, logstash will receive STDOUT/STDERR output log messages from `supervisor2` + `babl-server` and will forward the logs into kafka `logs.qa` topic.

`babl-qa` will be working independently and consuming data from kafka `logs.qa` topic.

Launch 6 terminal sessions in `$GOPATH/src/github.com/larskluge/babl-qa/DEVENV` and run:


	T1: $ ./logstash_start.sh && ./run_attach.sh

	T2: $ ./supervisor2.sh

	T3: $ ./babl-server.sh

	T4: $ ./babl-qa.sh

	T5: $ ./kafka-listen-logsqa.sh

	T6: $ ./babl-cli.sh


### TODO:

1. Babl module request history view

	- [x] Track request Start/Stop(Duration) from kafka `logs.qa` topic


2. Babl module request lifecycle view
	- [ ] Track request full lifecycle from kafka `logs.qa` topic
	- [ ] Send babl events (Telegram/Slack notifications) when the request lifecycle is not completed within the expected time limit.

3. babl-qa core modules

	- [x] Listen to kafka `logs.qa` topic and parse JSON log messages
	- [ ] Create a request lifecycle manager (receive messages and groups by requestid, should contain the 6 different lifecycle messages to assure the request was successfuly completed)
	- [ ] Create a REST api module (WIP)
