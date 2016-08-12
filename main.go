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
	kafkaClient   *sarama.Client
	kafkaProducer *sarama.SyncProducer
}

const Version = "1.0.0"
const clientID = "babl-qa"
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

	kafkaTopicQA := "logs.qa"
	kafkaTopicHistory := "logs.history"
	brokers := strings.Split(kafkaBrokers, ",")
	s.kafkaClient = kafka.NewClient(brokers, kafkaTopicQA, debug)
	defer (*s.kafkaClient).Close()
	s.kafkaProducer = kafka.NewProducer(brokers, clientID+".producer")
	defer (*s.kafkaProducer).Close()

	chQAMData := make(chan *QAMetadata)
	chQAMsg := make(chan *QAMessage)
	chQAHistory := make(chan *RequestHistory)

	go ListenToLogsQA(s.kafkaClient, kafkaTopicQA, chQAMData, chQAMsg)

	// other higher level go rotines go here
	go MonitorRequest(chQAMsg, chQAHistory)
	go SaveRequestHistory(s.kafkaProducer, kafkaTopicHistory, chQAHistory)
	//go SaveRequestLifecycle(s.kafkaProducer, kafkaTopicHistory, chQALifecycle)

	// block main process (will be replaced with HTTP server call)
	for {
		timer1 := time.NewTimer(time.Second * 30)
		<-timer1.C
	}
}
