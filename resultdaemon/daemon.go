package resultdaemon

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/gumeniukcom/achecker/configs"
	"github.com/gumeniukcom/achecker/kafka"
	"github.com/gumeniukcom/achecker/resultdaemon/structs"
	"sync"

	"github.com/Shopify/sarama"
	jrpc "github.com/gumeniukcom/golang-jsonrpc2"
	"github.com/rs/zerolog/log"
)

type Resoluter struct {
	client         sarama.Client
	consumer       sarama.ConsumerGroup
	resultTopic    string
	cfg            configs.Config
	serv           *jrpc.JSONRPC
	wg             *sync.WaitGroup
	consumerCancel context.CancelFunc
}

const daemonName = "resultdaemon"

func New(cfg configs.Config) *Resoluter {
	daemon := &Resoluter{
		cfg:         cfg,
		resultTopic: cfg.ResultDaemon.ResultTopic,
		wg:          &sync.WaitGroup{},
	}

	daemon.client = kafka.NewClient(daemon.cfg.ResultKafka)

	var err error
	daemon.consumer, err = sarama.NewConsumerGroupFromClient(daemon.cfg.ResultDaemon.KafkaGroup, daemon.client)

	if err != nil {
		log.Fatal().Err(err).Msg("Failed to start Sarama consumer")
	}

	// JSONRPC20 section
	daemon.serv = jrpc.New()

	if err := daemon.serv.RegisterMethod("save_check_domain", daemon.saveCheckDomain); err != nil {
		panic(err)
	}

	return daemon
}

func (daemon *Resoluter) Run() error {
	log.Debug().
		Str("daemon_name", daemonName).
		Msg("start daemon")

	handler := kafka.NewConsumer(daemon.serv)

	var ctx context.Context

	ctx, daemon.consumerCancel = context.WithCancel(context.Background())
	daemon.wg.Add(1)
	go func() {
		defer daemon.wg.Done()
		for {
			if err := daemon.consumer.Consume(ctx, []string{daemon.resultTopic}, handler); err != nil {
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
		Str("daemon_name", daemonName).
		Msg("started daemon")

	return nil

}

func (daemon *Resoluter) Stop() {
	log.Info().
		Str("daemon_name", daemonName).
		Msg("trying to stop daemon")

	daemon.consumerCancel()
	daemon.wg.Wait()

	log.Info().
		Str("daemon_name", daemonName).
		Msg("daemon stopped")
}

func (daemon *Resoluter) saveCheckDomain(ctx context.Context, data json.RawMessage) (json.RawMessage, int, error) {
	if data == nil {
		return nil, jrpc.InvalidRequestErrorCode, fmt.Errorf("empty request")
	}
	result := &structs.CheckResult{}
	err := result.UnmarshalJSON(data)
	if err != nil {
		return nil, jrpc.InvalidRequestErrorCode, err
	}

	log.Debug().Msgf("%#v", result)
	return nil, jrpc.OK, nil
}
