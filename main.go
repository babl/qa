package main

import (
	"os"
	"strings"
	"time"

	"github.com/Shopify/sarama"
	log "github.com/Sirupsen/logrus"
	. "github.com/larskluge/babl-qa/logs"
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

	ListenToLogsQA(s.kafkaClient, qaTopic)
}
