package main

import (
	. "github.com/larskluge/babl-qa/kafkalogs"
)

func MonitorRequestLifeCycle(chQAMsg chan *QAMessage) {
	for qalog := range chQAMsg {
		qalog.Debug()
	}
}
