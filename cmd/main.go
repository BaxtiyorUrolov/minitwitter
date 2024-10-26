package main

import (
	"context"
	"github.com/IBM/sarama"
	"go.uber.org/zap"
	"os"
	"os/signal"
	"syscall"
	"twitter/api"
	"twitter/api/handler"
	"twitter/config"
	"twitter/pkg/kafka" // kafka package, producer va consumer uchun
	"twitter/pkg/logger"
	"twitter/service"
	"twitter/storage/postgres"
)

func main() {

	cfg := config.Load()
	log := logger.New(cfg.ServiceName)

	// PostgreSQL bilan ulanish
	store, err := postgres.New(context.Background(), cfg, log)
	if err != nil {
		log.Error("Error while connecting to DB: %v", logger.Error(err))
		return
	}
	defer store.Close()

	// Kafka producer yaratish
	kafkaProducer, err := kafka.NewKafkaProducer([]string{cfg.KafkaHost + ":" + cfg.KafkaPort}, cfg.KafkaTopic, log)
	if err != nil {
		log.Error("Error while creating Kafka producer", logger.Error(err))
		return
	}
	defer kafkaProducer.Close()

	// Kafka consumer yaratish
	kafkaCfg := sarama.NewConfig()
	kafkaCfg.Consumer.Return.Errors = true
	kafkaClient, err := sarama.NewClient([]string{cfg.KafkaHost + ":" + cfg.KafkaPort}, kafkaCfg)
	if err != nil {
		log.Error("Error while connecting to Kafka client", logger.Error(err))
		return
	}
	defer kafkaClient.Close()

	kafkaConsumer, err := kafka.NewKafkaConsumer(kafkaClient, log)
	if err != nil {
		log.Error("Error while creating Kafka consumer", logger.Error(err))
		return
	}
	defer kafkaConsumer.Close()

	// Servislarni yaratish
	services := service.New(store, log, kafkaProducer)
	h := handler.New(services, log)

	// HTTP server yaratish
	server := api.New(cfg, services, store, log, h)

	// Signallarni kuzatish uchun kontekst yaratish
	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer cancel()

	// Kafka consumer'ni ishga tushirish
	go kafkaConsumer.Start(ctx)

	// HTTP serverni ishga tushirish
	go func() {
		if err := server.Run(); err != nil {
			log.Error("Error run http server", zap.Error(err))
		}
	}()

	// Signalni qabul qilish
	<-ctx.Done()

	// Serverni to'xtatish
	if err := server.Shutdown(ctx); err != nil {
		log.Error("failed http graceful shutdown", zap.Error(err))
	}
}
