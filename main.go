package main

import (
	"os"
	"strings"
	"time"

	"github.com/Shopify/sarama"
	log "github.com/Sirupsen/logrus"
	. "github.com/larskluge/babl-qa/kafkalogs"
	"github.com/larskluge/babl-server/kafka"
)

type server struct {
	kafkaClient *sarama.Client
}

const Version = "1.0.0"
const ModuleExecutionWaitTimeout = 5 * time.Minute

var debug bool

func main() {
	log.SetOutput(os.Stderr)
	log.SetFormatter(&log.JSONFormatter{})
	log.SetLevel(log.WarnLevel)

	app := configureCli()
	app.Run(os.Args)
}

func run(listen, kafkaBrokers string, dbg bool) {
	debug = dbg
	if debug {
		log.SetLevel(log.DebugLevel)
	}
	s := server{}

	qaTopic := "logs.qa"
	brokers := strings.Split(kafkaBrokers, ",")
	s.kafkaClient = kafka.NewClient(brokers, qaTopic, debug)
	defer (*s.kafkaClient).Close()

	chQALog := make(chan *QALog)
	chQAMsg := make(chan *QAMessage)
	//go debugLog(chQALog)
	//go debugMsg(chQAMsg)
	go ListenToLogsQA(s.kafkaClient, qaTopic, chQALog, chQAMsg)

	// other higher level go rotines go here
	go MonitorRequestHistory(chQAMsg)
	// -> go MonitorRequestLifeCycle(chQAMsg) // NOTE: Can not consume twice from the same channel !!! (chQAMsg)

	// block main process
	for {
	}
}

// TOBEREMOVED
func debugLog(chQALog chan *QALog) {
	for {
		qalog, ok := <-chQALog
		if !ok {
			panic(qalog)
		}
		qalog.Debug()
	}
}

// TOBEREMOVED
func debugMsg(chQAMsg chan *QAMessage) {
	for {
		qamsg, ok := <-chQAMsg
		if !ok {
			panic(qamsg)
		}
		qamsg.Debug()
	}
}
