package bablrequest

import (
	"strconv"
	"strings"
)

const (
	QAMsgTypeDefault = 1
	QAMsgTypeAsync   = 2
)
const (
	QAMsgSupervisor = "supervisor2"
	QAMsg1          = iota
	QAMsg2
	QAMsg3
	QAMsg4
	QAMsg5
	QAMsg6
	QAMsgTimeout = 999
)

type MessageState struct {
	Default MessageDefault
	Async   MessageAsync
}

func (this *MessageState) Initialize() {
	this.Default = MessageDefault{}
	this.Async = MessageAsync{}
	this.Default.Initialize()
	this.Async.Initialize()
}

func (this *MessageState) GetMessageType(qadata *QAJsonData) int {
	result := QAMsgTypeDefault
	// test Async message, if not defaults to normal message
	if strings.Contains(qadata.Message, this.Async.msg2) {
		result = QAMsgTypeAsync
	}
	return result
}

func (this *MessageState) GetProgress(msgType int, qadata *QAJsonData) int {
	if msgType == QAMsgTypeAsync {
		return this.Async.GetProgress(qadata)
	}
	return this.Default.GetProgress(qadata)
}

func (this *MessageState) GetProgressCompletion(msgType int, qadata *QAJsonData) string {
	if msgType == QAMsgTypeAsync {
		return this.Async.GetProgressCompletion(qadata)
	}
	return this.Default.GetProgressCompletion(qadata)
}

func (this *MessageState) GetProgressFromString(msgType int, message string) int {
	if msgType == QAMsgTypeAsync {
		return this.Async.GetProgressFromString(message, QAMsgSupervisor)
	}
	return this.Default.GetProgressFromString(message, QAMsgSupervisor)
}

func (this *MessageState) GetProgressCompletionFromString(msgType int, message string) string {
	if msgType == QAMsgTypeAsync {
		return this.Async.GetProgressCompletionFromString(message, QAMsgSupervisor)
	}
	return this.Default.GetProgressCompletionFromString(message, QAMsgSupervisor)
}

func (this *MessageState) FirstQAMsg(msgType int) int {
	if msgType == QAMsgTypeAsync {
		return this.Async.FirstQAMsg()
	}
	return this.Default.FirstQAMsg()
}

func (this *MessageState) LastQAMsg(msgType int) int {
	if msgType == QAMsgTypeAsync {
		return this.Async.LastQAMsg()
	}
	return this.Default.LastQAMsg()
}

type MessageDefault struct {
	/*
		#1 [supervisor2]	=>	"msg": "Producer: message sent", Topic:  babl.larskluge.StringUpcase.IO
		#2 [babl-server]	=>	"msg": "New Group Message Received"
		#3 [babl-server]	=>	"msg": "call"
		#4 [babl-server]	=>	"msg": "Producer: message sent", Topic:  supervisor.SevenMBP
		#5 [supervisor2]	=>	"msg": "New Message Received"
		#6 [supervisor2]	=>	"msg": "Module responded"
	*/
	msg1       string
	msg2       string
	msg3       string
	msg4       string
	msg5       string
	msg6       string
	msgTimeout string
}

func (this *MessageDefault) Initialize() {
	this.msg1 = "Producer: message sent"
	this.msg2 = "New Group Message Received"
	this.msg3 = "call"
	this.msg4 = "Producer: message sent"
	this.msg5 = "New Message Received"
	this.msg6 = "Module responded"
}

func (this *MessageDefault) FirstQAMsg() int {
	return QAMsg1
}

func (this *MessageDefault) LastQAMsg() int {
	return QAMsg6
}

func (this *MessageDefault) TimeoutQAMsg() int {
	return QAMsgTimeout
}

func (this *MessageDefault) GetProgress(qadata *QAJsonData) int {
	return this.GetProgressFromString(qadata.Message, qadata.Supervisor)
}

func (this *MessageDefault) GetProgressCompletion(qadata *QAJsonData) string {
	cProgress := strconv.Itoa(this.GetProgressFromString(qadata.Message, qadata.Supervisor))
	tProgress := strconv.Itoa(this.LastQAMsg())
	progress := cProgress + "/" + tProgress
	if cProgress == "0" {
		progress = "+"
	}
	return progress
}

func (this *MessageDefault) GetProgressFromString(message string, service string) int {
	result := int(0)
	if strings.Contains(message, this.msg1) && strings.Contains(service, QAMsgSupervisor) {
		result = QAMsg1
	} else if strings.Contains(message, this.msg2) {
		result = QAMsg2
	} else if strings.Contains(message, this.msg3) {
		result = QAMsg3
	} else if strings.Contains(message, this.msg4) && !strings.Contains(service, QAMsgSupervisor) {
		result = QAMsg4
	} else if strings.Contains(message, this.msg5) {
		result = QAMsg5
	} else if strings.Contains(message, this.msg6) {
		result = QAMsg6
	} else if strings.Contains(message, this.msgTimeout) {
		result = QAMsgTimeout
	}
	return result
}

func (this *MessageDefault) GetProgressCompletionFromString(message string, service string) string {
	cProgress := strconv.Itoa(this.GetProgressFromString(message, service))
	tProgress := strconv.Itoa(this.LastQAMsg())
	progress := cProgress + "/" + tProgress
	if cProgress == "0" {
		progress = "+"
	}
	return progress
}

type MessageAsync struct {
	/*
		#1 [supervisor2]	=>	"msg": "Producer: message sent", Topic:  babl.babl.Events.IO
		#2 [supervisor2]	=>	"msg": "Request processed async"
		#3 [babl-server]	=>	"msg": "New Group Message Received", Topic:  babl.babl.Events.IO
		#4 [babl-server]	=>	"msg": "call"
	*/
	msg1       string
	msg2       string
	msg3       string
	msg4       string
	msgTimeout string
}

func (this *MessageAsync) Initialize() {
	this.msg1 = "Producer: message sent"
	this.msg2 = "Request processed async"
	this.msg3 = "New Group Message Received"
	this.msg4 = "call"
	this.msgTimeout = "qa-service detected timeout"
}

func (this *MessageAsync) FirstQAMsg() int {
	return QAMsg1
}

func (this *MessageAsync) LastQAMsg() int {
	return QAMsg4
}

func (this *MessageAsync) TimeoutQAMsg() int {
	return QAMsgTimeout
}

func (this *MessageAsync) GetProgress(qadata *QAJsonData) int {
	return this.GetProgressFromString(qadata.Message, qadata.Supervisor)
}

func (this *MessageAsync) GetProgressCompletion(qadata *QAJsonData) string {
	cProgress := strconv.Itoa(this.GetProgressFromString(qadata.Message, qadata.Supervisor))
	tProgress := strconv.Itoa(this.LastQAMsg())
	progress := cProgress + "/" + tProgress
	if cProgress == "0" {
		progress = "+"
	}
	return progress
}

func (this *MessageAsync) GetProgressFromString(message string, service string) int {
	result := int(0)
	if strings.Contains(message, this.msg1) && strings.Contains(service, QAMsgSupervisor) {
		result = QAMsg1
	} else if strings.Contains(message, this.msg2) {
		result = QAMsg2
	} else if strings.Contains(message, this.msg3) {
		result = QAMsg3
	} else if strings.Contains(message, this.msg4) {
		result = QAMsg4
	} else if strings.Contains(message, this.msgTimeout) {
		result = QAMsgTimeout
	}
	return result
}

func (this *MessageAsync) GetProgressCompletionFromString(message string, service string) string {
	cProgress := strconv.Itoa(this.GetProgressFromString(message, service))
	tProgress := strconv.Itoa(this.LastQAMsg())
	progress := cProgress + "/" + tProgress
	if cProgress == "0" {
		progress = "+"
	}
	return progress
}
