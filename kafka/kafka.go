package kafka

import (
	"crypto/tls"
	"crypto/x509"
	"io/ioutil"

	"github.com/Shopify/sarama"
	"github.com/gumeniukcom/achecker/configs"
	"github.com/rs/zerolog/log"
)

// NewClient init sarama Client
func NewClient(conf configs.KafkaConf) sarama.Client {

	config := sarama.NewConfig()

	var err error

	config.Version, err = sarama.ParseKafkaVersion(conf.Version)
	config.Producer.RequiredAcks = sarama.WaitForAll // Wait for all in-sync replicas to ack the message
	config.Producer.Retry.Max = 10                   // Retry up to 10 times to produce the message
	config.Producer.Return.Successes = true

	if err != nil {
		log.Fatal().
			Err(err).
			Str("version", conf.Version).
			Msg("cannot convert version to uint slice")
	}

	if conf.SSl {
		config.Net.TLS.Enable = true
		tlsConfig, err := newTLSConfig(conf)
		if err != nil {
			log.Fatal().
				Err(err).
				Msg("error on create tls config")
		}
		config.Net.TLS.Config = tlsConfig
		log.Debug().Msg("tls ok")
	}

	client, err := sarama.NewClient(conf.Brokers, config)

	if err != nil {
		log.Fatal().
			Err(err).
			Strs("brokers", conf.Brokers).
			Msg("Failed to start Sarama client")
	}

	//producer, err := sarama.NewAsyncProducerFromClient(client)
	//
	//if err != nil {
	//	log.Fatal().Err(err).Msg("Failed to start Sarama producer")
	//}
	//
	//consumer, err := sarama.NewConsumerGroupFromClient(conf.Group, client)
	//
	//if err != nil {
	//	log.Fatal().Err(err).Msg("Failed to start Sarama consumer")
	//}
	//
	//// We will just log to STDOUT if we're not able to produce messages.
	//// Note: messages will only be returned here after all retry attempts are exhausted.
	//go func() {
	//	for err := range producer.Errors() {
	//		log.Error().Err(err).Msg("Failed to write access log entry")
	//	}
	//}()

	return client
}

func newTLSConfig(config configs.KafkaConf) (*tls.Config, error) {

	tlsConfig := tls.Config{}

	// Load client cert
	cert, err := tls.LoadX509KeyPair(config.FileCertPath, config.FileKeyPath)
	if err != nil {
		return &tlsConfig, err
	}
	tlsConfig.Certificates = []tls.Certificate{cert}

	// Load CA cert
	caCert, err := ioutil.ReadFile(config.FileCAPath)
	if err != nil {
		return &tlsConfig, err
	}
	caCertPool := x509.NewCertPool()
	caCertPool.AppendCertsFromPEM(caCert)
	tlsConfig.RootCAs = caCertPool

	return &tlsConfig, err
}
