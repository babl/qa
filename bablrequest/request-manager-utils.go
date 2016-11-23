package bablrequest

import (
	"encoding/json"
	"sort"
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
	if qadata.Status != "" {
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

func updateRequestDetails(progress int, progressCompletion string, qadata *QAJsonData) RequestDetails {
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
		Progress:  progressCompletion,
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
		//reqdet.Debug()
		// adjust step for the new msgType
		reqdet.Step = mState.GetProgressFromString(msgType, reqdet.Message)
		reqdet.Progress = mState.GetProgressCompletionFromString(msgType, reqdet.Message)
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

	// adjust optional messages Step == 0 to sequencial step > 100
	var optionalMessageStep = int(100)
	for i, reqdet := range rdOrigin {
		if reqdet.Step == 999 {
			rdOrigin[i].Step = optionalMessageStep
			rdOrigin[i].Progress = "+"
			optionalMessageStep++
		}
	}

	for _, reqdet := range rdOrigin {
		if reqdet.Step > 0 {
			keys = append(keys, reqdet.Step)
		}
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

func updateWSHistory(rh RequestHistory, chWSHist chan *[]byte) {
	type DataHistory struct {
		Type string         `json:"type"`
		Data RequestHistory `json:"data"`
	}
	wsData := &DataHistory{Type: "history", Data: rh}
	rhJson, _ := json.Marshal(wsData)
	chWSHist <- &rhJson
}

func updateWSDetails(rd []RequestDetails, chWSHist chan *[]byte) {
	type DataDetails struct {
		Type string           `json:"type"`
		Data []RequestDetails `json:"data"`
	}
	wsData := &DataDetails{Type: "details", Data: rd}
	rhJson, _ := json.Marshal(wsData)
	chWSHist <- &rhJson
}
