package kafkalogs

import (
	"strings"
)

/*
#1 [supervisor2]	=>	"msg": "Producer: message sent", Topic:  babl.larskluge.StringUpcase.IO
#2 [babl-server]	=>	"msg": "New Group Message Received"
#3 [babl-server]	=>	"msg": "call"
#4 [babl-server]	=>	"msg": "Producer: message sent", Topic:  supervisor.SevenMBP
#5 [supervisor2]	=>	"msg": "New Message Received"
#6 [supervisor2]	=>	"msg": "Module responded"
*/
const (
	supervisor = "supervisor2"
	msg1       = "Producer: message sent"
	msg2       = "New Group Message Received"
	msg3       = "call"
	msg4       = "Producer: message sent"
	msg5       = "New Message Received"
	msg6       = "Module responded"
	msgTimeout = "qa-service detected timeout"
)

const (
	QAMsg1 = iota + 1
	QAMsg2
	QAMsg3
	QAMsg4
	QAMsg5
	QAMsg6
	QAMsgTimeout = 999
)

func CheckMessageProgress(qadata *QAJsonData) int {
	result := int(0)
	if strings.Contains(qadata.Message, msg1) && strings.Contains(qadata.Supervisor, supervisor) {
		result = QAMsg1
	} else if strings.Contains(qadata.Message, msg2) {
		result = QAMsg2
	} else if strings.Contains(qadata.Message, msg3) {
		result = QAMsg3
	} else if strings.Contains(qadata.Message, msg4) && !strings.Contains(qadata.Supervisor, supervisor) {
		result = QAMsg4
	} else if strings.Contains(qadata.Message, msg5) {
		result = QAMsg5
	} else if strings.Contains(qadata.Message, msg6) {
		result = QAMsg6
	} else if strings.Contains(qadata.Message, msgTimeout) {
		result = QAMsgTimeout
	}
	return result
}
