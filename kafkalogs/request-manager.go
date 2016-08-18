package kafkalogs

import (
	"encoding/json"
	"fmt"
	"strconv"

	"github.com/Shopify/sarama"
	log "github.com/Sirupsen/logrus"
	"github.com/larskluge/babl-server/kafka"
	. "github.com/larskluge/babl-server/utils"
)

func updateRequestHistory(qalog *QALog, rh RequestHistory) RequestHistory {
	//qalog.DebugY()
	//qalog.DebugZ()
	data := rh
	data.Timestamp = qalog.Timestamp
	data.RequestId = qalog.RequestId
	data.Duration = qalog.Duration
	if data.Supervisor == "" && qalog.Service == "supervisor2" {
		data.Supervisor = qalog.Host
	}
	if data.Module == "" && qalog.Module != "" {
		data.Module = qalog.Module
	}
	if data.ModuleVersion == "" && qalog.ModuleVersion != "" {
		data.ModuleVersion = qalog.ModuleVersion
	}
	if qalog.Status != 0 {
		data.Status = qalog.Status
	}
	data.Duration = qalog.Duration
	return data
}

func updateRequestDetails(progress int, qalog *QALog) RequestDetails {
	data := RequestDetails{
		RequestHistory: RequestHistory{
			Timestamp:     qalog.Timestamp,
			RequestId:     qalog.RequestId,
			Module:        qalog.Module,
			ModuleVersion: qalog.ModuleVersion,
			Status:        qalog.Status,
			Duration:      qalog.Duration,
		},
		Host:      qalog.Host,
		Step:      progress,
		Topic:     qalog.Topic,
		Partition: qalog.Partition,
		Offset:    qalog.Offset,
	}
	data.Duration = qalog.Duration
	if qalog.Service == "supervisor2" {
		data.Supervisor = qalog.Host
	}
	//data.Debug()
	return data
}

func MonitorRequest(chQALog chan *QALog,
	chQAHist chan *RequestHistory, chQADetails chan *[]RequestDetails) {
	rhList := make(map[int32]RequestHistory)
	rdList := make(map[int32][]RequestDetails)

	for qalog := range chQALog {
		progress := CheckMessageProgress(qalog)

		// RequestHistory: update log messages
		rhList[qalog.RequestId] = updateRequestHistory(qalog, rhList[qalog.RequestId])
		// RequestHistory: send data to channel if last message arrived (QAMsg6)
		if progress == QAMsg6 {
			data := rhList[qalog.RequestId]
			chQAHist <- &data
			delete(rhList, qalog.RequestId)
		}

		// RequestDetails log messages
		rdList[qalog.RequestId] = append(rdList[qalog.RequestId], updateRequestDetails(progress, qalog))
		// RequestDetails: send data to channel if all 6 messages arrived (QAMsg1...QAMsg6)
		// NOTE: this is required due to the async nature of log messages: e.g.:
		// -> QAMsg2 -> QAMsg3 -> QAMsg4 -> QAMsg1 -> QAMsg6 -> QAMsg5
		if len(rdList[qalog.RequestId]) >= 6 {
			datadetails := rdList[qalog.RequestId]
			chQADetails <- &datadetails
			delete(rdList, qalog.RequestId)
		}

		/*
			data := rhList[qalog.RequestId]
			data.Timestamp = qalog.Timestamp
			data.RequestId = qalog.RequestId
			data.Duration = qalog.Duration
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
		*/

		/*
			rdList[qalog.RequestId] = append(rdList[qalog.RequestId],
				RequestDetails{
					RequestHistory: rhList[qalog.RequestId],
					Host:           qalog.Host,
					Step:           progress,
					Topic:          qalog.Topic,
					Partition:      qalog.Partition,
					Offset:         qalog.Offset,
				})
		*/
		/*
			if progress == QAMsg6 {
				data := rhList[qalog.RequestId]
				data.Duration = qalog.Duration // final duration_ms from supervisor2
				chQAHist <- &data

				datadetails := rdList[qalog.RequestId]
				chQADetails <- &datadetails

				delete(rhList, qalog.RequestId)
				delete(rdList, qalog.RequestId)

			}
		*/
	}
}

func SaveRequestHistory(producer *sarama.SyncProducer, topic string, chQAHist chan *RequestHistory) {
	for reqhist := range chQAHist {
		rhJson, _ := json.Marshal(reqhist)
		//fmt.Printf("%s\n", rhJson)
		kafka.SendMessage(producer, strconv.FormatInt(int64(reqhist.RequestId), 10), topic, &rhJson)
	}
}

func ReadRequestHistory(client *sarama.Client, topic string, lastn int64) []byte {
	log.Debug("Consuming from topic: ", topic)
	rhList := []RequestHistory{}
	ch := make(chan *kafka.ConsumerData)
	go kafka.ConsumeLastN(client, topic, 0, lastn, ch)
	for msg := range ch {
		log.WithFields(log.Fields{"key": msg.Key}).Debug("Read Request History message")

		reqhist := RequestHistory{}
		err1 := reqhist.UnmarshalJSON(msg.Value)
		//reqhist.Debug()
		Check(err1)
		rhList = append(rhList, reqhist)
		msg.Processed <- true
	}
	rhJson, _ := json.Marshal(rhList)
	return rhJson
}

func SaveRequestDetails(producer *sarama.SyncProducer, topic string, chQADetails chan *[]RequestDetails) {
	for reqdetails := range chQADetails {
		rid := (*reqdetails)[0].RequestId
		rhJson, _ := json.Marshal(reqdetails)
		fmt.Printf("%s\n", rhJson)
		kafka.SendMessage(producer, strconv.FormatInt(int64(rid), 10), topic, &rhJson)
	}
}

func ReadRequestDetails(client *sarama.Client, topic string, requestid string) []byte {
	log.Debug("Consuming from topic: ", topic)
	lastn := int64(100)
	rdList := []RequestDetails{}
	ch := make(chan *kafka.ConsumerData)
	go kafka.ConsumeLastN(client, topic, 0, lastn, ch) // need to replace with a full scan (stop after all steps)

	for msg := range ch {
		log.WithFields(log.Fields{"key": msg.Key}).Debug("Read Request History message")

		var arraymap []map[string]interface{}
		err := json.Unmarshal(msg.Value, &arraymap)
		Check(err)

		for _, v := range arraymap {
			rid, _ := strconv.ParseInt(requestid, 10, 32)
			if getFieldDataInt32(v["rid"]) == int32(rid) {
				reqdet := RequestDetails{}
				reqdet.ManualUnmarshalJSON(v)
				//reqdet.Debug()
				rdList = append(rdList, reqdet)
			}
		}
		msg.Processed <- true
	}
	rhJson, _ := json.Marshal(rdList)
	return rhJson
}
