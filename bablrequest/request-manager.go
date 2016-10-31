package bablrequest

import (
	"encoding/json"
	"fmt"
	"strconv"
	"time"

	"github.com/Shopify/sarama"
	log "github.com/Sirupsen/logrus"
	. "github.com/larskluge/babl-qa/httpserver"
	"github.com/larskluge/babl-server/kafka"
	. "github.com/larskluge/babl-server/utils"
	cache "github.com/muesli/cache2go"
)

const cacheDefaultExpiration = 7 * 24 * time.Hour

func MonitorRequest(chQAData chan *QAJsonData,
	chHist chan *RequestHistory, chWSHist chan *[]byte,
	chDetails chan *[]RequestDetails, cacheDetails *cache.CacheTable) {
	mState := MessageState{}
	mState.Initialize()
	rhList := make(map[string]RequestHistory)
	rdList := make(map[string][]RequestDetails)
	rdTimeout := make(map[string]time.Time)
	rdType := make(map[string]int)
	reqCompleted := make(map[string]bool)
	timeout := 6 * time.Minute
	waitBeforeSave := 10 * time.Second
	const timeoutStatus int32 = 408

	// monitor rdTimeout list to check if request takes longer than timeout
	// if request timeout occurs than sends data to chQAData channel
	go monitorRdTimeout(&rdTimeout, timeout, chQAData, timeoutStatus)

	for qadata := range chQAData {
		//qadata.DebugJson()

		//--------------------------------------------------------------------------
		// force sync. request to avoid map multi access (Add/Delete)
		FlushRequestData(&rhList, &rdList, &rdType, &reqCompleted)
		//--------------------------------------------------------------------------

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
		progressCompletion := mState.GetProgressCompletion(rdType[qadata.RequestId], qadata)

		// update rdTimeout with now():
		// + exclude previous timedout requests
		// + excludes messages with unkown progress
		if qadata.Status != timeoutStatus && progress > 0 {
			rdTimeout[qadata.RequestId] = time.Now()
		}

		//--------------------------------------------------------------------------
		// RequestHistory: update log messages
		//--------------------------------------------------------------------------
		rhList[qadata.RequestId] = updateRequestHistory(qadata, rhList[qadata.RequestId])
		// realtime update of the websockets channel
		updateWSHistory(rhList[qadata.RequestId], chWSHist)
		//--------------------------------------------------------------------------

		//--------------------------------------------------------------------------
		// RequestDetails log messages
		//--------------------------------------------------------------------------
		rdList[qadata.RequestId] = append(rdList[qadata.RequestId], updateRequestDetails(progress, progressCompletion, qadata))
		// realtime update of the details cache
		SaveCacheRequestDetails(qadata.RequestId, rdList[qadata.RequestId], cacheDetails)
		//--------------------------------------------------------------------------

		// RequestDetails: send data to channel if all 6 messages arrived (QAMsg1...QAMsg6)
		// NOTE: this is required due to the async nature of log messages: e.g.:
		// QAMsg2 -> QAMsg3 -> QAMsg4 -> QAMsg1 -> QAMsg6 -> QAMsg5
		if checkRequestDetailsCompleteSequence(rdType[qadata.RequestId], rdList[qadata.RequestId]) ||
			qadata.Status == timeoutStatus {

			// re-order RequestDetails data
			//rdList[qadata.RequestId] = orderbystepRequestDetails(rdList[qadata.RequestId])

			// resets rdTimeout, to be deleted by the monitor process
			rdTimeout[qadata.RequestId] = time.Time{}
			if _, ok := reqCompleted[qadata.RequestId]; !ok {
				// request is in the process to be complete (completed when true)
				reqCompleted[qadata.RequestId] = false
				fmt.Println("Process is about to complete: ", qadata.RequestId)
				//go WriteRequestData(waitBeforeSave, qadata.RequestId, &rhList, &rdList, &rdType, &rdCompleted, chHist, chDetails)
				go WriteRequestData(waitBeforeSave, qadata.RequestId, &rhList, &rdList, &reqCompleted, chHist, chDetails)
			}
		}
	}
}

func WriteRequestData(waitBeforeSave time.Duration, rid string,
	rhList *map[string]RequestHistory, rdList *map[string][]RequestDetails, reqCompleted *map[string]bool,
	chHist chan *RequestHistory, chDetails chan *[]RequestDetails) {

	fmt.Println("waitBeforeSave ...")
	timer1 := time.NewTimer(waitBeforeSave)
	<-timer1.C
	fmt.Println("saving data for rid=", rid)

	// write to kafka logs.history
	data := (*rhList)[rid]
	chHist <- &data

	// write to kafka logs.details
	// re-order RequestDetails data
	datadetails := orderbystepRequestDetails((*rdList)[rid])
	chDetails <- &datadetails

	// request completed
	(*reqCompleted)[rid] = true
	fmt.Println("Process is completed: ", rid)
}

func FlushRequestData(rhList *map[string]RequestHistory, rdList *map[string][]RequestDetails,
	rdType *map[string]int, rdCompleted *map[string]bool) {

	fmt.Println("rhList (before) => ", len(*rhList))
	fmt.Println("rdList (before) => ", len(*rdList))
	fmt.Println("rdType (before) => ", len(*rdType))
	fmt.Println("rdCompleted (before) => ", len(*rdCompleted))

	for rid, val := range *rdCompleted {
		fmt.Printf("[%s] => [%t]\n", rid, val)
		if val == true {
			if _, ok := (*rhList)[rid]; ok {
				delete(*rhList, rid)
				fmt.Println("rhList (after) => ", len(*rhList))
			}

			if _, ok := (*rdList)[rid]; ok {
				delete(*rdList, rid)
				fmt.Println("rdList (after) => ", len(*rdList))
			}

			if _, ok := (*rdType)[rid]; ok {
				delete(*rdType, rid)
				fmt.Println("rdType (after) => ", len(*rdType))
			}

			if _, ok := (*rdCompleted)[rid]; ok {
				delete(*rdCompleted, rid)
				fmt.Println("rdCompleted (after) => ", len(*rdCompleted))
			}
		}
	}
	/*
		for rid, val := range *rdCompleted {
			fmt.Printf("[%s] => [%t]\n", rid, val)
			if val == true {
				if _, ok := (*rdCompleted)[rid]; ok {
					delete(*rdCompleted, rid)
					fmt.Println("rdCompleted (after) => ", len(*rdCompleted))
				}
			}
		}
	*/
	fmt.Println("")
}

func SaveRequestHistory(producer *sarama.SyncProducer, topic string, chHist chan *RequestHistory) {
	for reqhist := range chHist {
		rhJson, _ := json.Marshal(reqhist)
		//fmt.Printf("%s\n", rhJson)
		fmt.Println("SaveRequestHistory rid=", reqhist.RequestId)
		kafka.SendMessage(producer, reqhist.RequestId, topic, &rhJson)
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

func SaveRequestDetails(producer *sarama.SyncProducer, topic string, chDetails chan *[]RequestDetails,
	cacheDetails *cache.CacheTable) {
	for reqdetails := range chDetails {
		requestID := (*reqdetails)[0].RequestId
		rhJson, _ := json.Marshal(reqdetails)
		//fmt.Printf("%s\n", rhJson)
		fmt.Println("SaveRequestDetails rid=", requestID)
		//cacheDetails.Add(requestID, cacheDefaultExpiration, rhJson)
		SaveCacheRequestDetails(requestID, *reqdetails, cacheDetails)
		kafka.SendMessage(producer, requestID, topic, &rhJson)
	}
}

func SaveCacheRequestDetails(requestID string, rdOrigin []RequestDetails,
	cacheDetails *cache.CacheTable) {
	rhJson, _ := json.Marshal(rdOrigin)
	//fmt.Printf("%s\n", rhJson)
	fmt.Println("SaveCacheRequestDetails rid=", requestID)
	cacheDetails.Add(requestID, cacheDefaultExpiration, rhJson)
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

	_, to := kafka.ConsumeGetOffsetValues(client, topic, 0)
	if to == 0 {
		log.Debug("Topic does not contain data to consume: ", topic)
		return 0
	}

	lastn := int64(10000)
	ch := make(chan *kafka.ConsumerData)
	go kafka.ConsumeLastN(client, topic, 0, lastn, ch)

	counter := 0
	for msg := range ch {
		log.WithFields(log.Fields{"key": msg.Key}).Debug("Read Request Details message")
		cacheDetails.Add(msg.Key, cacheDefaultExpiration, msg.Value)
		msg.Processed <- "success"
		if counter%1000 == 0 {
			log.Warn("Cache Loading: ", counter)
		}
		counter++
	}
	return cacheDetails.Count()
}

func ReadRequestPayload(client *sarama.Client, topic string, partition string, offset string) []byte {
	log.Debug("Consuming from topic/partition/offset: ", topic, partition, offset)
	partitionInt, _ := strconv.ParseInt(partition, 10, 64)
	offsetInt, _ := strconv.ParseInt(offset, 10, 64)
	result := kafka.ConsumeTopicPartitionOffset(client, topic, int32(partitionInt), offsetInt)
	return result.Value
}

func WSBroadcastRequestHistory(wsHub *Hub, chWSHist chan *[]byte) {
	for rhJson := range chWSHist {
		//fmt.Printf("%s\n", *rhJson)
		wsHub.Broadcast <- *rhJson
	}
}
