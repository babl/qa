package kafkalogs

import (
	"encoding/json"
	"fmt"
	"sort"
	"strconv"

	"github.com/Shopify/sarama"
	log "github.com/Sirupsen/logrus"
	"github.com/larskluge/babl-server/kafka"
	. "github.com/larskluge/babl-server/utils"
)

func updateRequestHistory(qadata *QAJsonData, rh RequestHistory) RequestHistory {
	//qadata.DebugJson()
	data := rh
	data.Timestamp = qadata.Timestamp
	data.RequestId = qadata.RequestId
	data.Duration = qadata.Duration
	if data.Supervisor == "" && qadata.Service == "supervisor2" {
		data.Supervisor = qadata.Host
	}
	if data.Module == "" && qadata.Module != "" {
		data.Module = qadata.Module
	}
	if data.ModuleVersion == "" && qadata.ModuleVersion != "" {
		data.ModuleVersion = qadata.ModuleVersion
	}
	if qadata.Status != 0 {
		data.Status = qadata.Status
	}
	data.Duration = qadata.Duration
	return data
}

func updateRequestDetails(progress int, qadata *QAJsonData) RequestDetails {
	data := RequestDetails{
		RequestHistory: RequestHistory{
			Timestamp:     qadata.Timestamp,
			RequestId:     qadata.RequestId,
			Module:        qadata.Module,
			ModuleVersion: qadata.ModuleVersion,
			Status:        qadata.Status,
			Duration:      qadata.Duration,
		},
		Host:      qadata.Host,
		Step:      progress,
		Topic:     qadata.Topic,
		Partition: qadata.Partition,
		Offset:    qadata.Offset,
	}
	data.Duration = qadata.Duration
	if qadata.Service == "supervisor2" {
		data.Supervisor = qadata.Host
	}
	//data.Debug()
	return data
}

func orderbystepRequestDetails(rdOrigin []RequestDetails) []RequestDetails {
	rdAuxList := make(map[int]RequestDetails)
	var rdResult []RequestDetails
	var keys []int
	for _, reqdet := range rdOrigin {
		keys = append(keys, reqdet.Step)
		rdAuxList[reqdet.Step] = reqdet
	}
	sort.Ints(keys)
	for _, k := range keys {
		rdResult = append(rdResult, rdAuxList[k])
	}
	return rdResult
}

func MonitorRequest(chQAData chan *QAJsonData,
	chQAHist chan *RequestHistory, chQADetails chan *[]RequestDetails) {
	rhList := make(map[int32]RequestHistory)
	rdList := make(map[int32][]RequestDetails)

	for qadata := range chQAData {
		progress := CheckMessageProgress(qadata)

		// RequestHistory: update log messages
		rhList[qadata.RequestId] = updateRequestHistory(qadata, rhList[qadata.RequestId])
		// RequestHistory: send data to channel if last message arrived (QAMsg6)
		if progress == QAMsg6 {
			data := rhList[qadata.RequestId]
			chQAHist <- &data
			delete(rhList, qadata.RequestId)
		}

		// RequestDetails log messages
		rdList[qadata.RequestId] = append(rdList[qadata.RequestId], updateRequestDetails(progress, qadata))
		// RequestDetails: send data to channel if all 6 messages arrived (QAMsg1...QAMsg6)
		// NOTE: this is required due to the async nature of log messages: e.g.:
		// -> QAMsg2 -> QAMsg3 -> QAMsg4 -> QAMsg1 -> QAMsg6 -> QAMsg5
		if len(rdList[qadata.RequestId]) >= 6 {
			rdList[qadata.RequestId] = orderbystepRequestDetails(rdList[qadata.RequestId])
			datadetails := rdList[qadata.RequestId]
			chQADetails <- &datadetails
			delete(rdList, qadata.RequestId)
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

func ReadRequestHistory(client *sarama.Client, topic string, lastn int64) []byte {
	log.Debug("Consuming from topic: ", topic)
	rhList := []RequestHistory{}
	ch := make(chan *kafka.ConsumerData)
	go kafka.ConsumeLastN(client, topic, 0, lastn, ch)
	for msg := range ch {
		log.WithFields(log.Fields{"key": msg.Key}).Debug("Read Request History message")

		reqhist := RequestHistory{}
		reqhist.UnmarshalJSON(msg.Value)
		//reqhist.Debug()
		rhList = append(rhList, reqhist)
		msg.Processed <- "success"
	}
	rhJson, _ := json.Marshal(rhList)
	return rhJson
}

func SaveRequestDetails(producer *sarama.SyncProducer, topic string, chQADetails chan *[]RequestDetails) {
	for reqdetails := range chQADetails {
		rid := (*reqdetails)[0].RequestId
		rhJson, _ := json.Marshal(reqdetails)
		//fmt.Printf("%s\n", rhJson)
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
		msg.Processed <- "success"
	}
	rhJson, _ := json.Marshal(rdList)
	return rhJson
}
