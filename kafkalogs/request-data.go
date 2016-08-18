package kafkalogs

import (
	"encoding/json"
	"fmt"
	"reflect"
	"time"
)

type RequestHistory struct {
	Timestamp     time.Time `json:"time"`
	RequestId     int32     `json:"rid"`
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
	reqhist.RequestId = getFieldDataInt32(Z["rid"])
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
	reqdet.RequestId = getFieldDataInt32(Z["rid"])
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

/*
RESTAPI: GET /api/request/logs
{
“timestamp” : ”2016-08-04T09:55:06Z”,
“requestid” : “123456”,
“supervisor” : “babl-queue1”,
“module” : “larskluge/bablbot”,
“duration_ms” : 125.75,
“status” : SUCCESS		// [‘SUCCESS’, ‘FAIL’, ‘TIMEOUT’]
}

RESTAPI: GET /api/request/details/12345
[{
“timestamp” : ”2016-08-04T09:55:06Z”,
“step” : “1”,
“requestid” : “123456”,
“supervisor” : “babl-queue1”,
“module” : “larskluge/string-upcase”,
“topic” : “qa-logs”,
“partition” : “0”,
“offset” : “55”,
“duration_ms” : 5.325,
},{
“timestamp” : ”2016-08-04T09:55:06Z”,
“step” : “2”,
“requestid” : “123456”,
“supervisor” : “babl-queue1”,
“module” : “larskluge/string-upcase”,
“topic” : “babl.larskluge.StringUpcase.IO”,
“partition” : “0”,
“offset” : “15”,
“duration_ms” : 55.125,
},{...},
{... “step” : 6}]
*/

/*
// MonitorRequest: ASYNC message arrival
{
  "time": "2016-08-18T13:54:41Z",
  "rid": 133937,
  "supervisor": "",
  "module": "larskluge/string-upcase",
  "moduleversion": "v46",
  "status": 0,
  "duration_ms": 0,
  "host": "babl-slave1",
  "step": 2,
  "topic": "babl.larskluge.StringUpcase.IO,babl.larskluge.StringUpcase.Ping",
  "partition": 7,
  "offset": 42
}
{
  "time": "2016-08-18T13:54:41Z",
  "rid": 133937,
  "supervisor": "",
  "module": "larskluge/string-upcase",
  "moduleversion": "v46",
  "status": 200,
  "duration_ms": 2.0692510000000004,
  "host": "babl-slave1",
  "step": 3,
  "topic": "",
  "partition": 0,
  "offset": 0
}
{
  "time": "2016-08-18T13:54:41Z",
  "rid": 133937,
  "supervisor": "",
  "module": "larskluge/string-upcase",
  "moduleversion": "v46",
  "status": 0,
  "duration_ms": 1.334052,
  "host": "babl-slave1",
  "step": 4,
  "topic": "supervisor.babl-queue1",
  "partition": 0,
  "offset": 0
}
{
  "time": "2016-08-18T13:54:41Z",
  "rid": 133937,
  "supervisor": "babl-queue1",
  "module": "",
  "moduleversion": "",
  "status": 0,
  "duration_ms": 112.66903700000002,
  "host": "babl-queue1",
  "step": 1,
  "topic": "babl.larskluge.StringUpcase.IO",
  "partition": 0,
  "offset": 0
}
{
  "time": "2016-08-18T13:54:41Z",
  "rid": 133937,
  "supervisor": "babl-queue1",
  "module": "",
  "moduleversion": "",
  "status": 0,
  "duration_ms": 285.694795,
  "host": "babl-queue1",
  "step": 6,
  "topic": "",
  "partition": 0,
  "offset": 0
}
{
  "time": "2016-08-18T13:54:41Z",
  "rid": 133937,
  "supervisor": "babl-queue1",
  "module": "",
  "moduleversion": "",
  "status": 0,
  "duration_ms": 0,
  "host": "babl-queue1",
  "step": 5,
  "topic": "supervisor.babl-queue1",
  "partition": 0,
  "offset": 250
}
*/