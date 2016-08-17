package kafkalogs

import (
	"encoding/json"
	"fmt"
	"reflect"
	"strconv"
	"strings"
	"time"
)

type QAMessage struct {
	Timestamp time.Time `json:"time"`
	RequestId int32     `json:"rid"`
	Key       string    `json:"key"`
	Message   string    `json:"msg"`
	Level     string    `json:"level"`
	Status    int32     `json:"status"`
	Stderr    string    `json:"stderr"`
	Topic     string    `json:"topic"`
	Partition int32     `json:"partition"`
	Offset    int32     `json:"offset"`
	ValueSize int32     `json:"value_size"`
	Duration  float64   `json:"duration_ms"`
	Z         map[string]interface{}
}

func (qamsg *QAMessage) UnmarshalJSON(data []byte) error {
	err := json.Unmarshal(data, &qamsg.Z)

	qamsg.Key = getFieldDataString(qamsg.Z["key"])
	qamsg.Message = getFieldDataString(qamsg.Z["msg"])
	qamsg.Level = getFieldDataString(qamsg.Z["level"])
	qamsg.Status = getFieldDataInt32(qamsg.Z["status"])
	qamsg.Stderr = getFieldDataString(qamsg.Z["stderr"])
	qamsg.Topic = getFieldDataString(qamsg.Z["topic"])
	qamsg.Partition = getFieldDataInt32(qamsg.Z["partition"])
	qamsg.Offset = getFieldDataInt32(qamsg.Z["offset"])
	qamsg.ValueSize = getFieldDataInt32(qamsg.Z["value_size"])
	qamsg.Duration = getFieldDataFloat64(qamsg.Z["duration_ms"])

	// custom fields conversion
	if isValidField(qamsg.Z["time"], reflect.String) {
		t1, _ := time.Parse(time.RFC3339, qamsg.Z["time"].(string))
		qamsg.Timestamp = t1
	}

	reqid := getFieldDataString(qamsg.Z["rid"])
	rid, errParse := strconv.ParseInt(reqid, 10, 64)
	if errParse == nil {
		qamsg.RequestId = int32(rid)
	}

	if isValidField(qamsg.Z["topics"], reflect.Slice) {
		str := ""
		vals := qamsg.Z["topics"].([]interface{})
		for _, val := range vals {
			str += " " + val.(string)
		}
		qamsg.Topic = strings.Replace(strings.Trim(str, " "), " ", ",", -1)
	}
	return err
}

func (qamsg *QAMessage) DebugZ() {
	rhJson, _ := json.Marshal(qamsg)
	fmt.Printf("%s\n", rhJson)
}
