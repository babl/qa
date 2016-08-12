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
func MonitorRequest(chQAMsg chan *QAMessage, chQAHist chan *RequestHistory) {
	rhList := make(map[int32]RequestHistory)

	for qalog := range chQAMsg {
		progress := CheckMessageProgress(qalog)
		//fmt.Println("MonitorRequestHistory: ", progress)
		//qalog.Debug()
		if progress == QAMsg1 {
			rhList[qalog.RequestId] = RequestHistory{
				Timestamp:  qalog.Timestamp,
				RequestId:  qalog.RequestId,
				Supervisor: strings.Split(qalog.Key, ".")[0],
				Module: strings.Split(qalog.Topic, ".")[1] + "." +
					strings.Split(qalog.Topic, ".")[2],
				Status:   0,
				Duration: 0.0,
			}
		}
		if progress == QAMsg3 {
			data := rhList[qalog.RequestId]
			data.Status = qalog.Status
			rhList[qalog.RequestId] = data
			//fmt.Println(qalog.Status)
		}
		if progress == QAMsg6 {
			data := rhList[qalog.RequestId]
			//data.Status = 1
			data.Duration = qalog.Duration
			rhList[qalog.RequestId] = data
			chQAHist <- &data
			delete(rhList, qalog.RequestId)
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
