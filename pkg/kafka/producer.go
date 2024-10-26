package kafka

import (
	"github.com/IBM/sarama"
	"twitter/pkg/logger"
)

type Producer struct {
	producer sarama.SyncProducer
	log      logger.Logger
	topic    string
}

func NewKafkaProducer(brokers []string, topic string, log logger.Logger) (*Producer, error) {
	config := sarama.NewConfig()
	config.Producer.Return.Successes = true
	config.Producer.Return.Errors = true

	producer, err := sarama.NewSyncProducer(brokers, config)
	if err != nil {
		log.Error("Failed to create Kafka kafka", logger.Error(err))
		return nil, err
	}

	return &Producer{
		producer: producer,
		log:      log,
		topic:    topic,
	}, nil
}

func (p *Producer) SendMessage(message string) error {
	msg := &sarama.ProducerMessage{
		Topic: p.topic,
		Value: sarama.StringEncoder(message),
	}

	partition, offset, err := p.producer.SendMessage(msg)
	if err != nil {
		p.log.Error("Failed to send message to Kafka", logger.Error(err))
		return err
	}

	p.log.Info("Message sent to Kafka", logger.Any("partition", partition), logger.Any("offset", offset))
	return nil
}

func (p *Producer) Close() error {
	return p.producer.Close()
}
