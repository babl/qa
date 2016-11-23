package bablrequest

import (
	"encoding/json"
	"fmt"
	"reflect"
	"strings"
	"time"
)

type QAJsonData struct {
	Type          string `json:"type"`
	Supervisor    string `json:"supervisor"`
	Host          string `json:"host"`
	Module        string `json:"module"`
	ModuleVersion string `json:"module_version"`
	ImageName     string `json:"image_name"`

	RequestId string  `json:"rid"`
	Key       string  `json:"key"`
	Message   string  `json:"message"`
	Error     string  `json:"message_error"`
	Level     string  `json:"level"`
	Code      string  `json:"code"`
	Status    string  `json:"status"`
	Stderr    string  `json:"stderr"`
	Topic     string  `json:"topic"`
	Partition int32   `json:"partition"`
	Offset    int32   `json:"offset"`
	ValueSize int32   `json:"value_size"`
	Duration  float64 `json:"duration_ms"`

	Timestamp time.Time `json:"timestamp"`

	Z map[string]interface{}
}

func (qadata *QAJsonData) UnmarshalJSON(data []byte) error {
	err := json.Unmarshal(data, &qadata.Z)

	qadata.Type = getFieldDataString(qadata.Z["type"])
	qadata.Host = getFieldDataString(qadata.Z["host"])
	qadata.Supervisor = getFieldDataString(qadata.Z["supervisor"])
	qadata.Module = getFieldDataString(qadata.Z["module"])
	qadata.ModuleVersion = getFieldDataString(qadata.Z["module_version"])
	qadata.ImageName = getFieldDataString(qadata.Z["image_name"])

	qadata.RequestId = getFieldDataString(qadata.Z["rid"])
	qadata.Key = getFieldDataString(qadata.Z["key"])
	qadata.Message = getFieldDataString(qadata.Z["message"])
	qadata.Error = getFieldDataString(qadata.Z["message_error"])
	qadata.Level = getFieldDataString(qadata.Z["level"])
	qadata.Code = getFieldDataString(qadata.Z["code"])
	qadata.Status = getFieldDataString(qadata.Z["status"])
	qadata.Stderr = getFieldDataString(qadata.Z["stderr"])
	qadata.Topic = getFieldDataString(qadata.Z["topic"])
	qadata.Partition = getFieldDataInt32(qadata.Z["partition"])
	qadata.Offset = getFieldDataInt32(qadata.Z["offset"])
	qadata.ValueSize = getFieldDataInt32(qadata.Z["value_size"])
	qadata.Duration = getFieldDataFloat64(qadata.Z["duration_ms"])

	// custom fields conversion
	if isValidField(qadata.Z["timestamp"], reflect.String) {
		t1, _ := time.Parse(time.RFC3339, qadata.Z["timestamp"].(string))
		qadata.Timestamp = t1
	}

	if isValidField(qadata.Z["topics"], reflect.Slice) {
		str := ""
		vals := qadata.Z["topics"].([]interface{})
		for _, val := range vals {
			str += " " + val.(string)
		}
		qadata.Topic = strings.Replace(strings.Trim(str, " "), " ", ",", -1)
	}
	return err
}

func (qadata *QAJsonData) Debug() {
	rhJson, _ := json.Marshal(qadata)
	fmt.Printf("%s\n", rhJson)
}

func (qadata *QAJsonData) DebugJson() {
	var result map[string]interface{}
	qatemp := qadata
	qatemp.Z = result
	rhJson, _ := json.Marshal(qatemp)
	fmt.Printf("%s\n", rhJson)
}
