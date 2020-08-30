package checkdaemon

import (
	"context"
	"encoding/json"
	"fmt"
	"sync"

	"github.com/gumeniukcom/achecker/checkdaemon/structs"
	"github.com/gumeniukcom/achecker/checker"
	"github.com/gumeniukcom/achecker/configs"
	"github.com/gumeniukcom/achecker/kafka"

	"github.com/Shopify/sarama"
	jrpc "github.com/gumeniukcom/golang-jsonrpc2"
	"github.com/rs/zerolog/log"
)

// Daemon container for app
type Daemon struct {
	cfg            configs.Config
	client         sarama.Client
	resultProducer sarama.AsyncProducer
	consumer       sarama.ConsumerGroup
	serv           *jrpc.JSONRPC
	checker        checker.Checker
	resultTopic    string
	checkTopic     string
	wg             *sync.WaitGroup
	consumerCancel context.CancelFunc
}

const (
	CheckDomainMethodName = "check_domain"
)

// New return new instance of Daemon
func New(cfg configs.Config) *Daemon {
	daemon := &Daemon{
		cfg:         cfg,
		checker:     checker.New(&cfg),
		checkTopic:  cfg.CheckDaemon.CheckTopic,
		resultTopic: cfg.CheckDaemon.ResultTopic,
		wg:          &sync.WaitGroup{},
	}

	daemon.client = kafka.NewClient(daemon.cfg.Kafka)

	var err error

	daemon.resultProducer, err = sarama.NewAsyncProducerFromClient(daemon.client)

	if err != nil {
		log.Fatal().Err(err).Msg("Failed to start Sarama producer")
	}

	daemon.consumer, err = sarama.NewConsumerGroupFromClient(cfg.Kafka.Group, daemon.client)

	if err != nil {
		log.Fatal().Err(err).Msg("Failed to start Sarama consumer")
	}

	// JSONRPC20 section
	daemon.serv = jrpc.New()

	if err := daemon.serv.RegisterMethod(CheckDomainMethodName, daemon.checkDomain); err != nil {
		panic(err)
	}

	return daemon
}

// Run run application
func (daemon *Daemon) Run() error {
	log.Debug().
		Str("daemon_name", "checker").
		Msg("start daemon")

	handler := kafka.NewConsumer(daemon.serv)

	var ctx context.Context

	ctx, daemon.consumerCancel = context.WithCancel(context.Background())
	daemon.wg.Add(1)
	go func() {
		defer daemon.wg.Done()
		for {
			if err := daemon.consumer.Consume(ctx, []string{daemon.checkTopic}, handler); err != nil {
				log.Fatal().
					Err(err).
					Msg("consume error")
			}

			if ctx.Err() != nil {
				log.Error().
					Msg("ctx finished")
				return
			}

			select {
			case <-ctx.Done():
				return
			default:
				continue
			}
		}
	}()

	<-handler.Ready // Await till the consumer has been set up

	log.Info().
		Str("daemon_name", "checker").
		Msg("started daemon")

	return nil
}

// Stop stop application
func (daemon *Daemon) Stop() {
	log.Info().
		Str("daemon_name", "checker").
		Msg("trying to stop daemon")

	daemon.consumerCancel()
	daemon.wg.Wait()

	if err := daemon.resultProducer.Close(); err != nil {
		log.Error().
			Str("daemon_name", "checker").
			Err(err).
			Msg("error on stop producer")
	}

	log.Info().
		Str("daemon_name", "checker").
		Msg("daemon stopped")
}

// TODO: think to not check in ths method -- send domain to chan and than check it in pool worker
func (daemon *Daemon) checkDomain(ctx context.Context, data json.RawMessage) (json.RawMessage, int, error) {
	if data == nil {
		return nil, jrpc.InvalidRequestErrorCode, fmt.Errorf("empty request")
	}
	task := &structs.Task{}
	err := task.UnmarshalJSON(data)
	if err != nil {
		return nil, jrpc.InvalidRequestErrorCode, err
	}

	resoluter, err := daemon.checker.CheckDomain(task.Domain)

	if err != nil {
		return nil, jrpc.InternalErrorCode, err
	}

	result := structs.CheckResult{
		Domain:     resoluter.Domain(),
		StatusCode: resoluter.StatusCode(),
	}

	if resoluter.Error() != nil {
		result.Error = resoluter.Error().Error()
	}

	jreq, err := jrpc.Request(context.Background(), "save_check_domain", result)
	if err != nil {
		log.Error().
			Err(err).
			Msg("error while creating jsonrpc20 request")
		return nil, jrpc.InternalErrorCode, err
	}

	if jreq == nil {
		log.Error().
			Err(err).
			Msg("empty request after creating jsonrpc20 request")
		return nil, jrpc.InternalErrorCode, err
	}

	msgByte, err := jreq.MarshalJSON()
	if err != nil {
		log.Error().
			Err(err).
			Msg("error on marshall jsonrpc20 request")
		return nil, jrpc.InternalErrorCode, err
	}

	daemon.resultProducer.Input() <- &sarama.ProducerMessage{
		Topic: daemon.resultTopic,
		Value: sarama.ByteEncoder(msgByte),
	}

	return nil, jrpc.OK, nil
}
