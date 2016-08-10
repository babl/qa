package kafkalogs

import (
	"encoding/json"
	"fmt"
)

type QALog struct {
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

func (qalog *QALog) UnmarshalJSON(data []byte) error {
	err := json.Unmarshal(data, &qalog.Z)
	qalog.Message = getFieldDataString(qalog.Z["msg"])
	qalog.Type = getFieldDataString(qalog.Z["type"])
	qalog.Level = getFieldDataString(qalog.Z["level"])
	qalog.Service = getFieldDataString(qalog.Z["service"])
	qalog.Module = getFieldDataString(qalog.Z["module"])
	qalog.ModuleVersion = getFieldDataString(qalog.Z["module_version"])
	qalog.ImageName = getFieldDataString(qalog.Z["image_name"])
	qalog.Port = getFieldDataInt32(qalog.Z["port"])
	return err
}

func (qalog *QALog) Debug() {
	fmt.Println("--------------------------------")
	fmt.Println("[QALog] => Message: ", qalog.Message)
	fmt.Println("[QALog] => Type: ", qalog.Type)
	fmt.Println("[QALog] => Level: ", qalog.Level)
	fmt.Println("[QALog] => Service: ", qalog.Service)
	fmt.Println("[QALog] => Module: ", qalog.Module)
	fmt.Println("[QALog] => ModuleVersion: ", qalog.ModuleVersion)
	fmt.Println("[QALog] => ImageName: ", qalog.ImageName)
	fmt.Println("[QALog] => Port: ", qalog.Port)
	fmt.Println("--------------------------------")
	fmt.Println("")
}
