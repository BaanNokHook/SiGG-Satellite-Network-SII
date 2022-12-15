// SiGG-Satellite-Network-SII  //

package nativelog

import (
	"fmt"
	"reflect"

	"github.com/Shopify/sarama"

	"github.com/apache/skywalking-satellite/internal/pkg/config"
	"github.com/apache/skywalking-satellite/internal/satellite/event"

	v1 "skywalking.apache.org/repo/goapi/satellite/data/v1"
)

const (
	Name     = "native-log-kafka-forwarder"
	ShowName = "Native Log Kafka Forwarder"
)

type Forwarder struct {
	config.CommonFields
	Topic    string `mapstructure:"topic"` // The forwarder topic.
	producer sarama.SyncProducer
}

func (f *Forwarder) Name() string {
	return Name
}

func (f *Forwarder) ShowName() string {
	return ShowName
}

func (f *Forwarder) Description() string {
	return "This is a synchronization Kafka forwarder with the SkyWalking native log protocol."
}

func (f *Forwarder) DefaultConfig() string {
	return `
# The remote topic. 
topic: "log-topic"
`
}

func (f *Forwarder) Prepare(connection interface{}) error {
	client, ok := connection.(sarama.Client)
	if !ok {
		return fmt.Errorf("the %s is only accepet the kafka client, but receive a %s",
			f.Name(), reflect.TypeOf(connection).String())
	}
	producer, err := sarama.NewSyncProducerFromClient(client)
	if err != nil {
		return err
	}
	f.producer = producer
	return nil
}

func (f *Forwarder) Forward(batch event.BatchEvents) error {
	var message []*sarama.ProducerMessage
	for _, e := range batch {
		data, ok := e.GetData().(*v1.SniffData_LogList)
		if !ok {
			continue
		}
		for _, l := range data.LogList.Logs {
			message = append(message, &sarama.ProducerMessage{
				Topic: f.Topic,
				Value: sarama.ByteEncoder(l),
			})
		}
	}
	return f.producer.SendMessages(message)
}

func (f *Forwarder) ForwardType() v1.SniffType {
	return v1.SniffType_Logging
}

func (f *Forwarder) SyncForward(_ *v1.SniffData) (*v1.SniffData, error) {
	return nil, fmt.Errorf("unsupport sync forward")
}

func (f *Forwarder) SupportedSyncInvoke() bool {
	return false
}
