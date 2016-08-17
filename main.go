package main

import (
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/Shopify/sarama"
	log "github.com/Sirupsen/logrus"
	"github.com/gorilla/mux"
	. "github.com/larskluge/babl-qa/httpserver"
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

	const kafkaTopicQA = "logs.qa"
	const kafkaTopicHistory = "logs.history"
	const kafkaTopicLifecycle = "logs.lifecycle"
	brokers := strings.Split(kafkaBrokers, ",")
	s.kafkaClient = kafka.NewClient(brokers, kafkaTopicQA, debug)
	defer (*s.kafkaClient).Close()
	s.kafkaProducer = kafka.NewProducer(brokers, clientID+".producer")
	defer (*s.kafkaProducer).Close()

	chQALog := make(chan *QALog)
	chQAHistory := make(chan *RequestHistory)
	chQADetails := make(chan *[]RequestDetails)

	go ListenToLogsQA(s.kafkaClient, kafkaTopicQA, chQALog)

	// other higher level go rotines go here
	go MonitorRequest(chQALog, chQAHistory, chQADetails)
	go SaveRequestHistory(s.kafkaProducer, kafkaTopicHistory, chQAHistory)
	go SaveRequestLifecycle(s.kafkaProducer, kafkaTopicLifecycle, chQADetails)

	// http callback function handler for Request History
	// $ http 127.0.0.1:8080/api/request/history
	// $ http 127.0.0.1:8080/api/request/history?blocksize=20
	HandlerRequestHistory := func(w http.ResponseWriter, r *http.Request) {
		lastn := GetVarsBlockSize(r, 10)
		rhJson := ReadRequestHistory(s.kafkaClient, kafkaTopicHistory, lastn)
		w.Header().Set("Content-Type", "application/json")
		w.Write(rhJson)
	}
	// http callback function handler for Request Details
	// http 127.0.0.1:8080/api/request/details/12345
	HandlerRequestDetails := func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		w.Write([]byte("Request Details! RequestId=" + vars["requestid"] + "\n"))
	}
	StartHttpServer(listen, HandlerRequestHistory, HandlerRequestDetails)
}
