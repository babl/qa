package kafkalogs

import (
	"encoding/json"
	"fmt"
)

type QAMetadata struct {
	Type          string `json:"type"`
	Level         string `json:"level"`
	Service       string `json:"service"`
	Module        string `json:"module"`
	ModuleVersion string `json:"module_version"`
	ImageName     string `json:"image_name"`
	Message       string `json:"message"`
	Port          int32  `json:"port"`
	Z             map[string]interface{}
}

func (qamdata *QAMetadata) UnmarshalJSON(data []byte) error {
	err := json.Unmarshal(data, &qamdata.Z)
	qamdata.Message = getFieldDataString(qamdata.Z["msg"])
	qamdata.Type = getFieldDataString(qamdata.Z["type"])
	qamdata.Level = getFieldDataString(qamdata.Z["level"])
	qamdata.Service = getFieldDataString(qamdata.Z["service"])
	qamdata.Module = getFieldDataString(qamdata.Z["module"])
	qamdata.ModuleVersion = getFieldDataString(qamdata.Z["module_version"])
	qamdata.ImageName = getFieldDataString(qamdata.Z["image_name"])
	qamdata.Port = getFieldDataInt32(qamdata.Z["port"])
	return err
}

func (qamdata *QAMetadata) Debug() {
	fmt.Println("--------------------------------")
	fmt.Println("[QAMetadata] => Message: ", qamdata.Message)
	fmt.Println("[QAMetadata] => Type: ", qamdata.Type)
	fmt.Println("[QAMetadata] => Level: ", qamdata.Level)
	fmt.Println("[QAMetadata] => Service: ", qamdata.Service)
	fmt.Println("[QAMetadata] => Module: ", qamdata.Module)
	fmt.Println("[QAMetadata] => ModuleVersion: ", qamdata.ModuleVersion)
	fmt.Println("[QAMetadata] => ImageName: ", qamdata.ImageName)
	fmt.Println("[QAMetadata] => Port: ", qamdata.Port)
	fmt.Println("--------------------------------")
	fmt.Println("")
}
