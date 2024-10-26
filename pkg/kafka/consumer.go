package kafka

import (
	"context"
	"encoding/json"
	"github.com/IBM/sarama"
	"twitter/pkg/email"
	"twitter/pkg/logger"
)

type KafkaConsumer struct {
	Consumer sarama.Consumer
	log      logger.Logger
}

type LikeNotificationEvent struct {
	Email   string `json:"email"`
	Message string `json:"message"`
}

func NewKafkaConsumer(client sarama.Client, log logger.Logger) (*KafkaConsumer, error) {
	consumer, err := sarama.NewConsumerFromClient(client)
	if err != nil {
		return nil, err
	}

	return &KafkaConsumer{
		Consumer: consumer,
		log:      log,
	}, nil
}

func (kc *KafkaConsumer) Start(ctx context.Context) {
	partitionConsumer, err := kc.Consumer.ConsumePartition("send_notification", 0, sarama.OffsetNewest)
	if err != nil {
		kc.log.Error("Error while starting consumer", logger.Error(err))
		return
	}
	defer partitionConsumer.Close()

	for {
		select {
		case msg := <-partitionConsumer.Messages():
			var event LikeNotificationEvent
			if err := json.Unmarshal(msg.Value, &event); err != nil {
				kc.log.Error("Error while unmarshalling message", logger.Error(err))
				continue
			}

			// Xabarni email orqali yuborish
			if err := email.SendEmail(event.Email, event.Message); err != nil {
				kc.log.Error("Error while sending email", logger.Error(err))
			}
		case <-ctx.Done():
			return
		}
	}
}

func (kc *KafkaConsumer) Close() error {
	return kc.Consumer.Close()
}
