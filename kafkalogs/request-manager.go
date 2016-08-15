package kafkalogs

import (
	"encoding/json"
	"fmt"
	"strconv"
	"time"

	"github.com/Shopify/sarama"
	"github.com/larskluge/babl-server/kafka"
)

type RequestHistory struct {
	Timestamp     time.Time `json:"time"`
	RequestId     int32     `json:"rid"`
	Supervisor    string    `json:"supervisor"`
	Module        string    `json:"module"`
	ModuleVersion string    `json:"moduleversion"`
	Status        int32     `json:"status"`
	Duration      float64   `json:"duration_ms"`
}

type RequestDetails struct {
	RequestHistory
	Host      string `json:"host"`
	Step      int    `json:"step"`
	Topic     string `json:"topic"`
	Partition int32  `json:"partition"`
	Offset    int32  `json:"offset"`
}

/*
RESTAPI: GET /api/request/logs
{
“timestamp” : ”2016-08-04T09:55:06Z”,
“requestid” : “123456”,
“supervisor” : “babl-queue1”,
“module” : “larskluge/bablbot”,
“duration_ms” : 125.75,
“status” : SUCCESS		// [‘SUCCESS’, ‘FAIL’, ‘TIMEOUT’]
}

RESTAPI: GET /api/request/details/12345
[{
“timestamp” : ”2016-08-04T09:55:06Z”,
“step” : “1”,
“requestid” : “123456”,
“supervisor” : “babl-queue1”,
“module” : “larskluge/string-upcase”,
“topic” : “qa-logs”,
“partition” : “0”,
“offset” : “55”,
“duration_ms” : 5.325,
},{
“timestamp” : ”2016-08-04T09:55:06Z”,
“step” : “2”,
“requestid” : “123456”,
“supervisor” : “babl-queue1”,
“module” : “larskluge/string-upcase”,
“topic” : “babl.larskluge.StringUpcase.IO”,
“partition” : “0”,
“offset” : “15”,
“duration_ms” : 55.125,
},{...},
{... “step” : 6}]
*/

func MonitorRequest(chQALog chan *QALog,
	chQAHist chan *RequestHistory, chQADetails chan *[]RequestDetails) {
	rhList := make(map[int32]RequestHistory)
	rdList := make(map[int32][]RequestDetails)

	for qalog := range chQALog {
		progress := CheckMessageProgress(qalog)
		//fmt.Println("MonitorRequestHistory: ", progress)
		//qalog.DebugY()
		//qalog.DebugZ()

		data := rhList[qalog.RequestId]
		data.Timestamp = qalog.Timestamp
		data.RequestId = qalog.RequestId
		if qalog.Service == "supervisor2" && data.Supervisor == "" {
			data.Supervisor = qalog.Host
		}
		if qalog.Module != "" && data.Module == "" {
			data.Module = qalog.Module
		}
		if qalog.ModuleVersion != "" && data.ModuleVersion == "" {
			data.ModuleVersion = qalog.ModuleVersion
		}
		if qalog.Status != 0 {
			data.Status = qalog.Status
		}
		//data.Duration = qalog.Duration
		rhList[qalog.RequestId] = data

		rdList[qalog.RequestId] = append(rdList[qalog.RequestId],
			RequestDetails{
				RequestHistory: rhList[qalog.RequestId],
				Host:           qalog.Host,
				Step:           progress,
				Topic:          qalog.Topic,
				Partition:      qalog.Partition,
				Offset:         qalog.Offset,
			})

		if progress == QAMsg6 {
			data := rhList[qalog.RequestId]
			data.Duration = qalog.Duration // final duration_ms from supervisor2
			chQAHist <- &data

			datadetails := rdList[qalog.RequestId]
			chQADetails <- &datadetails

			delete(rhList, qalog.RequestId)
			delete(rdList, qalog.RequestId)
		}
	}
}

func SaveRequestHistory(producer *sarama.SyncProducer, topic string, chQAHist chan *RequestHistory) {
	for reqhist := range chQAHist {
		rhJson, _ := json.Marshal(reqhist)
		fmt.Printf("%s\n", rhJson)
		kafka.SendMessage(producer, strconv.FormatInt(int64(reqhist.RequestId), 10), topic, &rhJson)
	}
}

func SaveRequestLifecycle(producer *sarama.SyncProducer, topic string, chQADetails chan *[]RequestDetails) {
	for reqdetails := range chQADetails {
		rhJson, _ := json.Marshal(reqdetails)
		rid := 123456 //reqdetails[0].RequestId // TODO: Needs to fix this!!!
		fmt.Printf("%s\n", rhJson)
		kafka.SendMessage(producer, strconv.FormatInt(int64(rid), 10), topic, &rhJson)
	}
}
