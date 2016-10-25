package bablrequest

import (
	"sort"
	"time"
)

func updateRequestHistory(qadata *QAJsonData, rh RequestHistory) RequestHistory {
	data := rh
	data.Timestamp = qadata.Timestamp
	data.RequestId = qadata.RequestId
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
	if data.Duration <= qadata.Duration {
		data.Duration = qadata.Duration
	}
	data.Message = qadata.Message
	if data.Error == "" && qadata.Error != "" {
		data.Error = qadata.Error
	}
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
			Error:         qadata.Error,
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

func updateRequestDetailsMsgType(msgType int, rdOrigin []RequestDetails) []RequestDetails {
	rdAuxList := make(map[int]RequestDetails)
	var rdResult []RequestDetails
	var keys []int
	mState := MessageState{}
	mState.Initialize()

	for _, reqdet := range rdOrigin {
		reqdet.Debug()
		// adjust step for the new msgType
		reqdet.Step = mState.GetProgressFromString(msgType, reqdet.Message)
		keys = append(keys, reqdet.Step)
		rdAuxList[reqdet.Step] = reqdet
	}
	sort.Ints(keys)
	for _, k := range keys {
		rdResult = append(rdResult, rdAuxList[k])
	}
	return rdResult
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

func checkRequestDetailsLastMsg(msgType int, rdOrigin []RequestDetails) bool {
	mState := MessageState{}
	mState.Initialize()
	result := false

	for _, reqdet := range rdOrigin {
		if reqdet.Step == mState.LastQAMsg(msgType) {
			result = true
			break
		}
	}
	return result
}

func checkRequestDetailsCompleteSequence(msgType int, rdOrigin []RequestDetails) bool {
	mState := MessageState{}
	mState.Initialize()

	// create empty array for QAMsg1...QAMsgx
	indexArray := make([]int, mState.LastQAMsg(msgType))
	// check the array with the steps found => messages arrived and identified correctly
	for _, reqdet := range rdOrigin {
		// test QAMsgTimeout condition
		if reqdet.Step == QAMsgTimeout {
			return false
		}
		if reqdet.Step >= mState.FirstQAMsg(msgType) && reqdet.Step <= mState.LastQAMsg(msgType) {
			indexArray[reqdet.Step-1] = 1
		}
	}
	// verify if all the array contains 1 in all elements
	result := true
	for _, v := range indexArray {
		if v != 1 {
			result = false
			break
		}
	}
	return result
}

func monitorRdTimeout(rdTL *map[string]time.Time, timeout time.Duration, chQAData chan *QAJsonData, statuscode int32) {
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
					timeoutData.Error = "QA-Service: Internal timer detected timeout!"
					chQAData <- &timeoutData
					delete(*rdTL, k)
				}
			}
		}
	}
}
