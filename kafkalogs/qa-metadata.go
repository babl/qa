package kafkalogs

import (
	"encoding/json"
	"fmt"
)

type QAMetadata struct {
	Type          string `json:"type"`
	Service       string `json:"service"`
	Host          string `json:"host"`
	Module        string `json:"module"`
	ModuleVersion string `json:"module_version"`
	ImageName     string `json:"image_name"`
	Y             map[string]interface{}
}

func (qamdata *QAMetadata) UnmarshalJSON(data []byte) error {
	err := json.Unmarshal(data, &qamdata.Y)
	qamdata.Type = getFieldDataString(qamdata.Y["type"])
	qamdata.Host = getFieldDataString(qamdata.Y["host"])
	qamdata.Service = getFieldDataString(qamdata.Y["service"])
	qamdata.Module = getFieldDataString(qamdata.Y["module"])
	qamdata.ModuleVersion = getFieldDataString(qamdata.Y["module_version"])
	qamdata.ImageName = getFieldDataString(qamdata.Y["image_name"])
	return err
}

func (qamdata *QAMetadata) DebugY() {
	rhJson, _ := json.Marshal(qamdata)
	fmt.Printf("%s\n", rhJson)
}
