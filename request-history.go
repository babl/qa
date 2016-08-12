package main

import (
	"encoding/json"
	"fmt"
	"strconv"
	"time"

	"github.com/Shopify/sarama"
	. "github.com/larskluge/babl-qa/kafkalogs"
	"github.com/larskluge/babl-server/kafka"
	. "github.com/larskluge/babl-server/utils"
)

type RequestHistory struct {
	Timestamp  time.Time `json:"time"`
	RequestId  int32     `json:"rid"`
	Supervisor string    `json:"supervisor"`
	Module     string    `json:"module"`
	Status     int32     `json:"status"`
	Duration   float64   `json:"duration_ms"`
}

/*
RESTAPI
{
“timestamp” : ”2016-08-04T09:55:06Z”,     -> #1
“requestid” : “123456”,                   -> #1
“supervisor” : “babl-queue1”,             -> #1 (key)
“module” : “larskluge/bablbot”,           -> #1
“duration_ms” : 125.75,                   -> #6
“status” : SUCCESS		// [‘SUCCESS’, ‘FAIL’, ‘TIMEOUT’] #6 (status + stderr)
},
*/
func MonitorRequestHistory(chQAMsg chan *QAMessage, chQAHist chan *RequestHistory) {
	rhList := make(map[int32]RequestHistory) // this needs to be optimized, should use a circular list

	for qalog := range chQAMsg {
		progress := CheckMessageProgress(qalog)
		if progress == QAMsg1 {
			rhList[qalog.RequestId] = RequestHistory{
				Timestamp:  qalog.Timestamp,
				RequestId:  qalog.RequestId,
				Supervisor: SplitGetByIndex(qalog.Key, ".", 0),
				Module: SplitGetByIndex(qalog.Topic, ".", 1) + "." +
					SplitGetByIndex(qalog.Topic, ".", 2),
				Status:   0,
				Duration: 0.0,
			}
		}
		if progress == QAMsg3 {
			data := rhList[qalog.RequestId]
			data.Status = qalog.Status
			rhList[qalog.RequestId] = data
			fmt.Println(qalog.Status)
		}
		if progress == QAMsg6 {
			data := rhList[qalog.RequestId]
			//data.Status = 1
			data.Duration = qalog.Duration
			rhList[qalog.RequestId] = data
			chQAHist <- &data
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
