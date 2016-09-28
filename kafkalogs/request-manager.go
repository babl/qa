package kafkalogs

import (
	"encoding/json"
	"sort"
	"strconv"
	"time"

	"github.com/Shopify/sarama"
	log "github.com/Sirupsen/logrus"
	"github.com/larskluge/babl-server/kafka"
	. "github.com/larskluge/babl-server/utils"
	cache "github.com/muesli/cache2go"
)

const cacheDefaultExpiration = 7 * 24 * time.Hour

func updateRequestHistory(qadata *QAJsonData, rh RequestHistory) RequestHistory {
	//qadata.DebugJson()
	data := rh
	data.Timestamp = qadata.Timestamp
	data.RequestId = qadata.RequestId
	data.Duration = qadata.Duration
	if data.Supervisor == "" && qadata.Supervisor != "" {
		data.Supervisor = qadata.Supervisor
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
	data.Message = qadata.Message
	//data.Debug()
	return data
}

func updateRequestDetails(progress int, qadata *QAJsonData) RequestDetails {
	//qadata.DebugJson()
	data := RequestDetails{
		RequestHistory: RequestHistory{
			Timestamp:     qadata.Timestamp,
			RequestId:     qadata.RequestId,
			Message:       qadata.Message,
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
	data.Supervisor = qadata.Supervisor
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

func monitorRdTimeout(rdTL *map[int32]time.Time, timeout time.Duration, chQAData chan *QAJsonData, statuscode int32) {
	for {
		timer1 := time.NewTimer(time.Second)
		<-timer1.C
		for k, v := range *rdTL {
			// verify if request was completed with success
			if v == (time.Time{}) {
				delete(*rdTL, k)
			} else {
				elapsed := time.Since(v)
				if elapsed > timeout {
					timeoutData := QAJsonData{}
					timeoutData.RequestId = k
					timeoutData.Duration = float64(elapsed) / float64(time.Millisecond)
					timeoutData.Level = "error"
					timeoutData.Status = statuscode
					timeoutData.Timestamp = time.Now()
					timeoutData.Message = "qa-service detected timeout!"
					chQAData <- &timeoutData
					delete(*rdTL, k)
				}
			}
		}
	}
}

func MonitorRequest(chQAData chan *QAJsonData,
	chQAHist chan *RequestHistory, chQADetails chan *[]RequestDetails) {
	rhList := make(map[int32]RequestHistory)
	rdList := make(map[int32][]RequestDetails)
	rdTimeout := make(map[int32]time.Time)
	timeout := 5 * time.Minute
	const timeoutStatus int32 = 408

	// monitor rdTimeout list to check if request takes longer than timeout
	// if request timeout occurs than sends data to chQAData channel
	go monitorRdTimeout(&rdTimeout, timeout, chQAData, timeoutStatus)

	for qadata := range chQAData {
		//qadata.DebugJson()
		progress := CheckMessageProgress(qadata)

		// update rdTimeout with now(): exclude previous timedout requests
		if qadata.Status != timeoutStatus {
			rdTimeout[qadata.RequestId] = time.Now()
		}

		// RequestHistory: update log messages
		rhList[qadata.RequestId] = updateRequestHistory(qadata, rhList[qadata.RequestId])
		// RequestHistory: send data to channel if last message arrived (QAMsg6)
		if progress == QAMsg6 || qadata.Status == timeoutStatus {
			data := rhList[qadata.RequestId]
			chQAHist <- &data
			delete(rhList, qadata.RequestId)
		}

		// RequestDetails log messages
		rdList[qadata.RequestId] = append(rdList[qadata.RequestId], updateRequestDetails(progress, qadata))
		// RequestDetails: send data to channel if all 6 messages arrived (QAMsg1...QAMsg6)
		// NOTE: this is required due to the async nature of log messages: e.g.:
		// -> QAMsg2 -> QAMsg3 -> QAMsg4 -> QAMsg1 -> QAMsg6 -> QAMsg5
		if len(rdList[qadata.RequestId]) >= 6 || qadata.Status == timeoutStatus {
			rdList[qadata.RequestId] = orderbystepRequestDetails(rdList[qadata.RequestId])
			datadetails := rdList[qadata.RequestId]
			chQADetails <- &datadetails
			delete(rdList, qadata.RequestId)
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
