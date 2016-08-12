package kafkalogs

import (
	"encoding/json"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/Shopify/sarama"
	"github.com/larskluge/babl-server/kafka"
)

type RequestHistory struct {
	Timestamp  time.Time `json:"time"`
	RequestId  int32     `json:"rid"`
	Supervisor string    `json:"supervisor"`
	Module     string    `json:"module"`
	Status     int32     `json:"status"`
	Duration   float64   `json:"duration_ms"`
}

type RequestDetails struct {
	RequestHistory
	Step      int    `json:"step"`
	Topic     string `json:"topic"`
	Partition int32  `json:"partition"`
	Offset    int32  `json:"offset"`
}

/*
RESTAPI: GET /api/request/logs
{
“timestamp” : ”2016-08-04T09:55:06Z”,     -> #1
“requestid” : “123456”,                   -> #1
“supervisor” : “babl-queue1”,             -> #1 (key)
“module” : “larskluge/bablbot”,           -> #1
“duration_ms” : 125.75,                   -> #6
“status” : SUCCESS		// [‘SUCCESS’, ‘FAIL’, ‘TIMEOUT’] #6 (status + stderr)
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

func MonitorRequest(chQAMsg chan *QAMessage,
	chQAHist chan *RequestHistory, chQADetails chan *[]RequestDetails) {
	rhList := make(map[int32]RequestHistory)
	rdList := make(map[int32][]RequestDetails)

	for qamsg := range chQAMsg {
		progress := CheckMessageProgress(qamsg)
		//fmt.Println("MonitorRequestHistory: ", progress)
		//qamsg.Debug()

		if progress == QAMsg1 {
			// TODO: optimize
			data := RequestHistory{}
			data.Timestamp = qamsg.Timestamp
			data.RequestId = qamsg.RequestId
			data.Supervisor = strings.Split(qamsg.Key, ".")[0]
			data.Module = strings.Split(qamsg.Topic, ".")[1] + "." + strings.Split(qamsg.Topic, ".")[2]
			data.Status = qamsg.Status
			data.Duration = qamsg.Duration
			rhList[qamsg.RequestId] = data

			rdList[qamsg.RequestId] = append(rdList[qamsg.RequestId],
				RequestDetails{
					RequestHistory: rhList[qamsg.RequestId],
					Step:           progress,
					Topic:          qamsg.Topic,
					Partition:      qamsg.Partition,
					Offset:         qamsg.Offset,
				})
		}
		if progress == QAMsg2 {
			// TODO: optimize
			data := RequestHistory{}
			data = rhList[qamsg.RequestId]
			data.Timestamp = qamsg.Timestamp
			data.Status = qamsg.Status
			data.Duration = qamsg.Duration
			rhList[qamsg.RequestId] = data

			rdList[qamsg.RequestId] = append(rdList[qamsg.RequestId],
				RequestDetails{
					RequestHistory: rhList[qamsg.RequestId],
					Step:           progress,
					Topic:          qamsg.Topic,
					Partition:      qamsg.Partition,
					Offset:         qamsg.Offset,
				})
		}
		if progress == QAMsg3 {
			// TODO: optimize
			data := RequestHistory{}
			data = rhList[qamsg.RequestId]
			data.Timestamp = qamsg.Timestamp
			data.Status = qamsg.Status
			data.Duration = qamsg.Duration
			rhList[qamsg.RequestId] = data

			rdList[qamsg.RequestId] = append(rdList[qamsg.RequestId],
				RequestDetails{
					RequestHistory: rhList[qamsg.RequestId],
					Step:           progress,
					Topic:          qamsg.Topic,
					Partition:      qamsg.Partition,
					Offset:         qamsg.Offset,
				})
		}
		if progress == QAMsg4 {
			// TODO: optimize
			data := RequestHistory{}
			data = rhList[qamsg.RequestId]
			data.Timestamp = qamsg.Timestamp
			data.Status = qamsg.Status
			data.Duration = qamsg.Duration
			rhList[qamsg.RequestId] = data

			rdList[qamsg.RequestId] = append(rdList[qamsg.RequestId],
				RequestDetails{
					RequestHistory: rhList[qamsg.RequestId],
					Step:           progress,
					Topic:          qamsg.Topic,
					Partition:      qamsg.Partition,
					Offset:         qamsg.Offset,
				})
		}
		if progress == QAMsg5 {
			// TODO: optimize
			data := RequestHistory{}
			data = rhList[qamsg.RequestId]
			data.Timestamp = qamsg.Timestamp
			data.Status = qamsg.Status
			data.Duration = qamsg.Duration
			rhList[qamsg.RequestId] = data

			rdList[qamsg.RequestId] = append(rdList[qamsg.RequestId],
				RequestDetails{
					RequestHistory: rhList[qamsg.RequestId],
					Step:           progress,
					Topic:          qamsg.Topic,
					Partition:      qamsg.Partition,
					Offset:         qamsg.Offset,
				})
		}
		if progress == QAMsg6 {
			// TODO: optimize
			data := RequestHistory{}
			data = rhList[qamsg.RequestId]
			data.Timestamp = qamsg.Timestamp
			data.Status = qamsg.Status
			data.Duration = qamsg.Duration
			rhList[qamsg.RequestId] = data
			chQAHist <- &data

			rdList[qamsg.RequestId] = append(rdList[qamsg.RequestId],
				RequestDetails{
					RequestHistory: rhList[qamsg.RequestId],
					Step:           progress,
					Topic:          qamsg.Topic,
					Partition:      qamsg.Partition,
					Offset:         qamsg.Offset,
				})
			datadetails := rdList[qamsg.RequestId]
			chQADetails <- &datadetails

			delete(rhList, qamsg.RequestId)
			delete(rdList, qamsg.RequestId)
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
