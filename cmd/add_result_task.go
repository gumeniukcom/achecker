package main

import (
	"context"
	"github.com/Shopify/sarama"
	"github.com/gumeniukcom/achecker/configs"
	"github.com/gumeniukcom/achecker/kafka"
	"github.com/gumeniukcom/achecker/resultdaemon/structs"
	jrpc "github.com/gumeniukcom/golang-jsonrpc2"
	"github.com/rs/zerolog/log"
)

func main() {
	cfg := configs.ReadConfig("./config.toml")

	client := kafka.NewClient(cfg.Kafka)
	producer, err := sarama.NewSyncProducerFromClient(client)

	if err != nil {
		log.Fatal().Err(err).Msg("Failed to start Sarama producer")
	}

	domains := []string{"ya.ru", "gumeniuk.cpm"}
	for _, domain := range domains {
		if err := addResult(producer, cfg.CheckDaemon.ResultTopic, domain); err != nil {
			log.Error().Err(err).Str("domain", domain).Msg("failed add task")
		}
	}
	if err := producer.Close(); err != nil {
		log.Error().Err(err).Msg("failed close producer")
	}
}

func addResult(producer sarama.SyncProducer, topic string, domain string) error {
	result := structs.CheckResult{
		Domain:     domain,
		StatusCode: 200,
	}

	jreq, err := jrpc.Request(context.Background(), "save_check_domain", result)
	if err != nil {
		log.Error().
			Err(err).
			Msg("error while creating jsonrpc20 request")
		return err
	}

	if jreq == nil {
		log.Error().
			Err(err).
			Msg("empty request after creating jsonrpc20 request")
		return err
	}

	msgByte, err := jreq.MarshalJSON()
	if err != nil {
		log.Error().
			Err(err).
			Msg("error on marshall jsonrpc20 request")
		return err
	}

	if _, _, err = producer.SendMessage(&sarama.ProducerMessage{
		Topic: topic,
		Value: sarama.ByteEncoder(msgByte),
	}); err != nil {
		log.Error().
			Err(err).
			Str("domain", domain).
			Str("topic", topic).
			Msg("failed send result")
		return err
	}

	log.Info().
		Str("domain", domain).
		Str("topic", topic).
		Msg("result added")

	return nil
}
