package kafkalogs

import (
	"github.com/Shopify/sarama"
	log "github.com/Sirupsen/logrus"
	"github.com/larskluge/babl-server/kafka"
	. "github.com/larskluge/babl-server/utils"
)

func ListenToLogsQA(client *sarama.Client, topic string, chQALog chan *QALog, chQAMsg chan *QAMessage) {
	log.Debug("Consuming from qa topic")
	ch := make(chan *kafka.ConsumerData)
	//go kafka.Consume(client, topic, ch, kafka.ConsumerOptions{Offset: sarama.OffsetOldest})
	go kafka.Consume(client, topic, ch)
	for msg := range ch {
		log.WithFields(log.Fields{"key": msg.Key}).Debug("QA message received")

		// parse top level log message (include logstash metadata)
		qalog := QALog{}
		err1 := qalog.UnmarshalJSON(msg.Value)
		Check(err1)

		// parse low level log message (logstash "message" property)
		qamsg := QAMessage{}
		err2 := qamsg.UnmarshalJSON([]byte(qalog.Z["message"].(string)))
		Check(err2)

		go func() { chQALog <- &qalog }()
		go func() { chQAMsg <- &qamsg }()
		// chQALog <- &qalog
		// chQAMsg <- &qamsg

		msg.Processed <- true
	}
	panic("listenToQAMessages: Lost connection to Kafka")
}
