package kafkalogs

import (
	"github.com/Shopify/sarama"
	log "github.com/Sirupsen/logrus"
	"github.com/larskluge/babl-server/kafka"
	. "github.com/larskluge/babl-server/utils"
)

func ListenToLogsQA(client *sarama.Client, topic string, chQAData chan *QAJsonData) {
	log.Debug("Consuming from topic: ", topic)
	ch := make(chan *kafka.ConsumerData)
	go kafka.Consume(client, topic, ch, kafka.ConsumerOptions{Offset: sarama.OffsetNewest})
	for msg := range ch {
		log.WithFields(log.Fields{"key": msg.Key}).Debug("QA message received")

		qadata := QAJsonData{}
		err1 := qadata.UnmarshalJSON(msg.Value)
		qadata.DebugJson()
		Check(err1)

		if qadata.RequestId > 0 {
			go func() { chQAData <- &qadata }()
		}
		msg.Processed <- "success"
	}
	panic("listenToQAMessages: Lost connection to Kafka")
}
