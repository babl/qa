package kafkalogs

import (
	"strings"
)

/*
#1
--------------------------------
[QAMessage] => Timestamp:  2016-08-11 14:29:46 +0100 WEST
[QAMessage] => RequestId:  243632
[QAMessage] => Key:  SevenMBP.243632
[QAMessage] => Message:  Producer: message sent
[QAMessage] => Level:  info
[QAMessage] => Status:  0
[QAMessage] => Stderr:
[QAMessage] => Topic:  babl.larskluge.StringUpcase.IO
[QAMessage] => Partition:  0
[QAMessage] => Offset:  0
[QAMessage] => ValueSize:  0
[QAMessage] => Duration:  129.717736
--------------------------------

#2
--------------------------------
[QAMessage] => Timestamp:  2016-08-11 14:29:46 +0100 WEST
[QAMessage] => RequestId:  243632
[QAMessage] => Key:  SevenMBP.243632
[QAMessage] => Message:  New Group Message Received
[QAMessage] => Level:  info
[QAMessage] => Status:  0
[QAMessage] => Stderr:
[QAMessage] => Topic:  babl.larskluge.StringUpcase.IO,babl.larskluge.StringUpcase.Ping
[QAMessage] => Partition:  5
[QAMessage] => Offset:  20
[QAMessage] => ValueSize:  12
[QAMessage] => Duration:  0
--------------------------------

#3
--------------------------------
[QAMessage] => Timestamp:  2016-08-11 14:29:46 +0100 WEST
[QAMessage] => RequestId:  243632
[QAMessage] => Key:
[QAMessage] => Message:  call
[QAMessage] => Level:  info
[QAMessage] => Status:  200
[QAMessage] => Stderr:
[QAMessage] => Topic:
[QAMessage] => Partition:  0
[QAMessage] => Offset:  0
[QAMessage] => ValueSize:  0
[QAMessage] => Duration:  6.994899
--------------------------------

#4
--------------------------------
[QAMessage] => Timestamp:  2016-08-11 14:29:47 +0100 WEST
[QAMessage] => RequestId:  243632
[QAMessage] => Key:  243632
[QAMessage] => Message:  Producer: message sent
[QAMessage] => Level:  info
[QAMessage] => Status:  0
[QAMessage] => Stderr:
[QAMessage] => Topic:  supervisor.SevenMBP
[QAMessage] => Partition:  0
[QAMessage] => Offset:  0
[QAMessage] => ValueSize:  0
[QAMessage] => Duration:  465.035653
--------------------------------

#5
--------------------------------
[QAMessage] => Timestamp:  2016-08-11 14:29:47 +0100 WEST
[QAMessage] => RequestId:  243632
[QAMessage] => Key:  243632
[QAMessage] => Message:  New Message Received
[QAMessage] => Level:  info
[QAMessage] => Status:  0
[QAMessage] => Stderr:
[QAMessage] => Topic:  supervisor.SevenMBP
[QAMessage] => Partition:  0
[QAMessage] => Offset:  260
[QAMessage] => ValueSize:  12
[QAMessage] => Duration:  0
--------------------------------

#6
--------------------------------
[QAMessage] => Timestamp:  2016-08-11 14:29:47 +0100 WEST
[QAMessage] => RequestId:  243632
[QAMessage] => Key:
[QAMessage] => Message:  Module responded
[QAMessage] => Level:  info
[QAMessage] => Status:  0
[QAMessage] => Stderr:
[QAMessage] => Topic:
[QAMessage] => Partition:  0
[QAMessage] => Offset:  0
[QAMessage] => ValueSize:  0
[QAMessage] => Duration:  810.957566
--------------------------------
*/

/*
#1 [supervisor2]	=>	"msg": "Producer: message sent", Topic:  babl.larskluge.StringUpcase.IO
#2 [babl-server]	=>	"msg": "New Group Message Received"
#3 [babl-server]	=>	"msg": "call"
#4 [babl-server]	=>	"msg": "Producer: message sent", Topic:  supervisor.SevenMBP
#5 [supervisor2]	=>	"msg": "New Message Received"
#6 [supervisor2]	=>	"msg": "Module responded"
*/
const (
	supervisor = "supervisor"
	msg1       = "Producer: message sent"
	msg2       = "New Group Message Received"
	msg3       = "call"
	msg4       = "Producer: message sent"
	msg5       = "New Message Received"
	msg6       = "Module responded"
)

const (
	QAMsg1 = iota + 1
	QAMsg2
	QAMsg3
	QAMsg4
	QAMsg5
	QAMsg6
)

func CheckMessageProgress(qamsg *QAMessage) int {
	result := int(0)
	if strings.Contains(qamsg.Message, msg1) && !strings.Contains(qamsg.Topic, supervisor) {
		result = QAMsg1
	} else if strings.Contains(qamsg.Message, msg2) {
		result = QAMsg2
	} else if strings.Contains(qamsg.Message, msg3) {
		result = QAMsg3
	} else if strings.Contains(qamsg.Message, msg4) && strings.Contains(qamsg.Topic, supervisor) {
		result = QAMsg4
	} else if strings.Contains(qamsg.Message, msg5) {
		result = QAMsg5
	} else if strings.Contains(qamsg.Message, msg6) {
		result = QAMsg6
	}
	return result
}
