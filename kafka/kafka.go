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

	config.Producer.RequiredAcks = sarama.WaitForAll // Wait for all in-sync replicas to ack the message
	config.Producer.Retry.Max = 10                   // Retry up to 10 times to produce the message
	config.Producer.Return.Successes = true

	config.Version, err = sarama.ParseKafkaVersion(conf.Version)

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
