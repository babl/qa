package kafkalogs

import (
	"encoding/json"
	"fmt"
	"reflect"
	"time"
)

type RequestHistory struct {
	Timestamp     time.Time `json:"time"`
	RequestId     string    `json:"rid"`
	Message       string    `json:"message"`
	Error         string    `json:"message_error"`
	Supervisor    string    `json:"supervisor"`
	Module        string    `json:"module"`
	ModuleVersion string    `json:"moduleversion"`
	Status        int32     `json:"status"`
	Duration      float64   `json:"duration_ms"`
}

type RequestDetails struct {
	RequestHistory
	Host      string `json:"host"`
	Step      int    `json:"step"`
	Topic     string `json:"topic"`
	Partition int32  `json:"partition"`
	Offset    int32  `json:"offset"`
}

func (reqhist *RequestHistory) UnmarshalJSON(data []byte) error {
	Z := make(map[string]interface{})
	err := json.Unmarshal(data, &Z)

	if isValidField(Z["time"], reflect.String) {
		t1, _ := time.Parse(time.RFC3339, Z["time"].(string))
		reqhist.Timestamp = t1
	}
	reqhist.RequestId = getFieldDataString(Z["rid"])
	reqhist.Message = getFieldDataString(Z["message"])
	reqhist.Error = getFieldDataString(Z["message_error"])
	reqhist.Supervisor = getFieldDataString(Z["supervisor"])
	reqhist.Module = getFieldDataString(Z["module"])
	reqhist.ModuleVersion = getFieldDataString(Z["moduleversion"])
	reqhist.Status = getFieldDataInt32(Z["status"])
	reqhist.Duration = getFieldDataFloat64(Z["duration_ms"])
	return err
}

func (reqhist *RequestHistory) Debug() {
	rhJson, _ := json.Marshal(reqhist)
	fmt.Printf("%s\n", rhJson)
}

func (reqdet *RequestDetails) ManualUnmarshalJSON(Z map[string]interface{}) error {

	if isValidField(Z["time"], reflect.String) {
		t1, _ := time.Parse(time.RFC3339, Z["time"].(string))
		reqdet.Timestamp = t1
	}
	reqdet.RequestId = getFieldDataString(Z["rid"])
	reqdet.Message = getFieldDataString(Z["message"])
	reqdet.Error = getFieldDataString(Z["message_error"])
	reqdet.Supervisor = getFieldDataString(Z["supervisor"])
	reqdet.Module = getFieldDataString(Z["module"])
	reqdet.ModuleVersion = getFieldDataString(Z["moduleversion"])
	reqdet.Status = getFieldDataInt32(Z["status"])
	reqdet.Duration = getFieldDataFloat64(Z["duration_ms"])

	reqdet.Host = getFieldDataString(Z["host"])
	reqdet.Step = getFieldDataInt(Z["step"])
	reqdet.Topic = getFieldDataString(Z["topic"])
	reqdet.Partition = getFieldDataInt32(Z["partition"])
	reqdet.Offset = getFieldDataInt32(Z["offset"])

	return nil
}

func (reqdet *RequestDetails) Debug() {
	rhJson, _ := json.Marshal(reqdet)
	fmt.Printf("%s\n", rhJson)
}
