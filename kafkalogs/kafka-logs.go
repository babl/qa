package kafkalogs

import (
	"github.com/Shopify/sarama"
	log "github.com/Sirupsen/logrus"
	"github.com/larskluge/babl-server/kafka"
	. "github.com/larskluge/babl-server/utils"
)

type QALog struct {
	QAMetadata
	QAMessage
}

func ListenToLogsQA(client *sarama.Client, topic string, chQALog chan *QALog) {
	log.Debug("Consuming from topic: ", topic)
	ch := make(chan *kafka.ConsumerData)
	go kafka.Consume(client, topic, ch, kafka.ConsumerOptions{Offset: sarama.OffsetNewest})
	for msg := range ch {
		log.WithFields(log.Fields{"key": msg.Key}).Debug("QA message received")

		// parse top level log message (include logstash metadata)
		qamdata := QAMetadata{}
		err1 := qamdata.UnmarshalJSON(msg.Value)
		//qamdata.DebugY()
		Check(err1)

		// parse low level log message (logstash "message" property)
		qamsg := QAMessage{}
		err2 := qamsg.UnmarshalJSON([]byte(qamdata.Y["message"].(string)))
		//qamsg.DebugZ()
		Check(err2)

		if qamsg.RequestId > 0 {
			qalog := QALog{qamdata, qamsg}
			go func() { chQALog <- &qalog }()
		}
		msg.Processed <- "success"
	}
	panic("listenToQAMessages: Lost connection to Kafka")
}
