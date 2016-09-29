package kafkalogs

import (
	"encoding/json"
	"strconv"
	"time"

	"github.com/Shopify/sarama"
	log "github.com/Sirupsen/logrus"
	"github.com/larskluge/babl-server/kafka"
	. "github.com/larskluge/babl-server/utils"
	cache "github.com/muesli/cache2go"
)

const cacheDefaultExpiration = 7 * 24 * time.Hour

func MonitorRequest(chQAData chan *QAJsonData,
	chQAHist chan *RequestHistory, chQADetails chan *[]RequestDetails) {
	mState := MessageState{}
	mState.Initialize()
	rhList := make(map[int32]RequestHistory)
	rdList := make(map[int32][]RequestDetails)
	rdTimeout := make(map[int32]time.Time)
	rdType := make(map[int32]int)
	timeout := 5 * time.Minute
	const timeoutStatus int32 = 408

	// monitor rdTimeout list to check if request takes longer than timeout
	// if request timeout occurs than sends data to chQAData channel
	go monitorRdTimeout(&rdTimeout, timeout, chQAData, timeoutStatus)

	for qadata := range chQAData {
		//qadata.DebugJson()

		// verify message type, defaults to normal message and sets to Async in special case
		if _, ok := rdType[qadata.RequestId]; !ok {
			rdType[qadata.RequestId] = mState.GetMessageType(qadata)
		} else if mState.GetMessageType(qadata) == QAMsgTypeAsync && rdType[qadata.RequestId] == QAMsgTypeDefault {
			rdType[qadata.RequestId] = QAMsgTypeAsync
			// Needs to adjust rdList[qadata.RequestId] messages with new step codes
			rdList[qadata.RequestId] = updateRequestDetailsMsgType(QAMsgTypeAsync, rdList[qadata.RequestId])
		}

		// sets the message step
		progress := mState.GetProgress(rdType[qadata.RequestId], qadata)

		// update rdTimeout with now(): exclude previous timedout requests
		if qadata.Status != timeoutStatus {
			rdTimeout[qadata.RequestId] = time.Now()
		}

		// RequestHistory: update log messages
		rhList[qadata.RequestId] = updateRequestHistory(qadata, rhList[qadata.RequestId])
		// RequestHistory: send data to channel if last message arrived (QAMsg6)
		if progress == mState.LastQAMsg(rdType[qadata.RequestId]) ||
			checkRequestDetailsLastMsg(rdType[qadata.RequestId], rdList[qadata.RequestId]) ||
			qadata.Status == timeoutStatus {
			data := rhList[qadata.RequestId]
			chQAHist <- &data
			delete(rhList, qadata.RequestId)
		}

		// RequestDetails log messages
		rdList[qadata.RequestId] = append(rdList[qadata.RequestId], updateRequestDetails(progress, qadata))
		// RequestDetails: send data to channel if all 6 messages arrived (QAMsg1...QAMsg6)
		// NOTE: this is required due to the async nature of log messages: e.g.:
		// QAMsg2 -> QAMsg3 -> QAMsg4 -> QAMsg1 -> QAMsg6 -> QAMsg5
		if checkRequestDetailsCompleteSequence(rdType[qadata.RequestId], rdList[qadata.RequestId]) ||
			qadata.Status == timeoutStatus {

			rdList[qadata.RequestId] = orderbystepRequestDetails(rdList[qadata.RequestId])
			datadetails := rdList[qadata.RequestId]
			chQADetails <- &datadetails
			delete(rdList, qadata.RequestId)
			delete(rdType, qadata.RequestId)
			// resets rdTimeout, to be deleted by the monitor process
			rdTimeout[qadata.RequestId] = time.Time{}
		}
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
		reqhist.UnmarshalJSON(msg.Value)
		//reqhist.Debug()
		rhList = append(rhList, reqhist)
		msg.Processed <- "success"
	}
	rhJson, _ := json.Marshal(rhList)
	return rhJson
}

func SaveRequestDetails(producer *sarama.SyncProducer, topic string, chQADetails chan *[]RequestDetails,
	cacheDetails *cache.CacheTable) {
	for reqdetails := range chQADetails {
		rid := (*reqdetails)[0].RequestId
		rhJson, _ := json.Marshal(reqdetails)
		//fmt.Printf("%s\n", rhJson)
		requestID := strconv.FormatInt(int64(rid), 10)
		cacheDetails.Add(requestID, cacheDefaultExpiration, rhJson)
		kafka.SendMessage(producer, requestID, topic, &rhJson)
	}
}

func ReadRequestDetailsFromCache(requestid string, cacheDetails *cache.CacheTable) []byte {
	log.Debug("Reading from cache: Details %d", cacheDetails.Count())

	value, err := cacheDetails.Value(requestid)
	if err != nil {
		rdList := []RequestDetails{}
		rdJson, _ := json.Marshal(rdList)
		return rdJson
	}
	valueData := value.Data()
	valueByte := valueData.([]byte)

	var arraymap []map[string]interface{}
	err1 := json.Unmarshal(valueByte, &arraymap)
	Check(err1)

	rdList := []RequestDetails{}
	for _, v := range arraymap {
		reqdet := RequestDetails{}
		reqdet.ManualUnmarshalJSON(v)
		//reqdet.Debug()
		rdList = append(rdList, reqdet)
	}
	rdJson, _ := json.Marshal(rdList)
	return rdJson
}

func ReadRequestDetailsToCache(client *sarama.Client, topic string, cacheDetails *cache.CacheTable) int {
	log.Debug("Consuming from topic: ", topic)

	lastn := int64(999999)
	ch := make(chan *kafka.ConsumerData)
	go kafka.ConsumeLastN(client, topic, 0, lastn, ch)

	for msg := range ch {
		log.WithFields(log.Fields{"key": msg.Key}).Debug("Read Request Details message")
		cacheDetails.Add(msg.Key, cacheDefaultExpiration, msg.Value)
		msg.Processed <- "success"
	}
	return cacheDetails.Count()
}
