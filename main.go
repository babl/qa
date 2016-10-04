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
	cache "github.com/muesli/cache2go"
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
	const kafkaTopicDetails = "logs.details"
	brokers := strings.Split(kafkaBrokers, ",")
	s.kafkaClient = kafka.NewClient(brokers, kafkaTopicQA, debug)
	defer (*s.kafkaClient).Close()
	s.kafkaProducer = kafka.NewProducer(brokers, clientID+".producer")
	defer (*s.kafkaProducer).Close()

	chQAData := make(chan *QAJsonData)
	chQAHistory := make(chan *RequestHistory)
	chQADetails := make(chan *[]RequestDetails)
	cacheDetails := cache.Cache("cacheDetails")

	// cache details
	ReadRequestDetailsToCache(s.kafkaClient, kafkaTopicDetails, cacheDetails)

	go ListenToLogsQA(s.kafkaClient, kafkaTopicQA, chQAData)

	// other higher level go rotines go here
	go MonitorRequest(chQAData, chQAHistory, chQADetails)
	go SaveRequestHistory(s.kafkaProducer, kafkaTopicHistory, chQAHistory)
	go SaveRequestDetails(s.kafkaProducer, kafkaTopicDetails, chQADetails, cacheDetails)

	// http callback function handler for Request History
	// $ http 127.0.0.1:8888/api/request/history
	// $ http 127.0.0.1:8888/api/request/history?blocksize=20
	HandlerRequestHistory := func(w http.ResponseWriter, r *http.Request) {
		lastn := GetVarsBlockSize(r, 10)
		rhJson := ReadRequestHistory(s.kafkaClient, kafkaTopicHistory, lastn)
		w.Header().Set("Content-Type", "application/json")
		w.Write(rhJson)
	}
	// http callback function handler for Request Details
	// http 127.0.0.1:8888/api/request/details/12345
	HandlerRequestDetails := func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		rhJson := ReadRequestDetailsFromCache(vars["requestid"], cacheDetails)
		w.Header().Set("Content-Type", "application/json")
		w.Write(rhJson)
	}
	// http callback function handler for Request Payload
	// http 127.0.0.1:8888/api/request/payload/babl.babl.Events.IO/4/3460
	HandlerRequestPayload := func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		rhJson := ReadRequestPayload(s.kafkaClient, vars["topic"], vars["partition"], vars["offset"])
		w.Header().Set("Content-Type", "application/octet-stream")
		w.Write(rhJson)
	}
	StartHttpServer(listen, HandlerRequestHistory, HandlerRequestDetails, HandlerRequestPayload)
}
