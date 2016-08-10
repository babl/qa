package kafkalogs

import (
	"fmt"

	"github.com/Shopify/sarama"
	log "github.com/Sirupsen/logrus"
	"github.com/larskluge/babl-server/kafka"
	. "github.com/larskluge/babl-server/utils"
)

func ListenToLogsQA(client *sarama.Client, topic string) {
	log.Debug("Consuming from qa topic")
	ch := make(chan *kafka.ConsumerData)
	go kafka.Consume(client, topic, ch)
	for msg := range ch {
		log.WithFields(log.Fields{"key": msg.Key}).Debug("QA message received")

		// parse top level log message (include logstash metadata)
		qalog := QALog{}
		err1 := qalog.UnmarshalJSON(msg.Value)
		Check(err1)

		// parse low level log message (logstash "message" property)
		qamessage := QAMessage{}
		err2 := qamessage.UnmarshalJSON([]byte(qalog.Z["message"].(string)))
		Check(err2)

		fmt.Println("--------------------------------")
		fmt.Println("Timestamp: ", qamessage.Timestamp)
		fmt.Println("RequestId: ", qamessage.RequestId)
		fmt.Println("Message: ", qamessage.Message)
		fmt.Println("Level: ", qamessage.Level)
		fmt.Println("Status: ", qamessage.Status)
		fmt.Println("Stderr: ", qamessage.Stderr)
		fmt.Println("Topic: ", qamessage.Topic)
		fmt.Println("Partition: ", qamessage.Partition)
		fmt.Println("Offset: ", qamessage.Offset)
		fmt.Println("ValueSize: ", qamessage.ValueSize)
		fmt.Println("Duration: ", qamessage.Duration)
		fmt.Println("--------------------------------")
		fmt.Println("")

		msg.Processed <- true
	}
	panic("listenToQAMessages: Lost connection to Kafka")
}
