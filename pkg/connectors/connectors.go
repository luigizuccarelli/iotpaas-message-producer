package connectors

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/confluentinc/confluent-kafka-go/kafka"
	"github.com/luigizuccarelli/iotpaas-message-producer/pkg/schema"
	"github.com/microlib/simple"
)

type Clients interface {
	Error(string, ...interface{})
	Info(string, ...interface{})
	Debug(string, ...interface{})
	Trace(string, ...interface{})
	SendMessageSync(body []byte) error
	Close()
}

type Connectors struct {
	Producer *kafka.Producer
	Logger   *simple.Logger
	Name     string
}

func NewClientConnectors(logger *simple.Logger) Clients {
	logger.Trace("Creating kafka connections for message producer")
	p, err := kafka.NewProducer(&kafka.ConfigMap{"bootstrap.servers": os.Getenv("KAFKA_BROKERS")})
	if err != nil {
		logger.Error(fmt.Sprintf("Creating kafka connections %v", err))
		panic(err)
	}

	//defer p.Close()
	return &Connectors{Producer: p, Logger: logger, Name: "RealConnectors"}
}

func (conn *Connectors) SendMessageSync(b []byte) error {

	topic := os.Getenv("TOPIC")
	var data *schema.IOTPaaS

	err := json.Unmarshal(b, &data)
	if err != nil {
		return err
	}

	conn.Producer.Produce(&kafka.Message{
		TopicPartition: kafka.TopicPartition{Topic: &topic},
		Key:            []byte(data.Id),
		Value:          b,
	}, nil)

	e := <-conn.Producer.Events()
	msg := e.(*kafka.Message)
	if msg.TopicPartition.Error != nil {
		return msg.TopicPartition.Error
	}

	return nil
}

func (conn *Connectors) Close() {
	conn.Producer.Close()
}

func (conn *Connectors) Error(msg string, val ...interface{}) {
	conn.Logger.Error(fmt.Sprintf(msg, val...))
}

func (conn *Connectors) Info(msg string, val ...interface{}) {
	conn.Logger.Info(fmt.Sprintf(msg, val...))
}

func (conn *Connectors) Debug(msg string, val ...interface{}) {
	conn.Logger.Debug(fmt.Sprintf(msg, val...))
}

func (conn *Connectors) Trace(msg string, val ...interface{}) {
	conn.Logger.Trace(fmt.Sprintf(msg, val...))
}
