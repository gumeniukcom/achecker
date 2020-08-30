package kafka

import (
	"context"
	"github.com/Shopify/sarama"
	jrpc "github.com/gumeniukcom/golang-jsonrpc2"
	"github.com/rs/zerolog/log"
)

// Consumer represents a Sarama consumer group consumer
type Consumer struct {
	Ready chan bool
	jrpc  *jrpc.JSONRPC
}

func NewConsumer(serv *jrpc.JSONRPC) *Consumer {
	return &Consumer{
		Ready: make(chan bool),
		jrpc:  serv,
	}
}

// Setup is run at the beginning of a new session, before ConsumeClaim
func (consumer *Consumer) Setup(sarama.ConsumerGroupSession) error {
	// Mark the consumer as ready
	close(consumer.Ready)
	return nil
}

// Cleanup is run at the end of a session, once all ConsumeClaim goroutines have exited
func (consumer *Consumer) Cleanup(sarama.ConsumerGroupSession) error {
	return nil
}

// ConsumeClaim must start a consumer loop of ConsumerGroupClaim's Messages().
func (consumer *Consumer) ConsumeClaim(session sarama.ConsumerGroupSession, claim sarama.ConsumerGroupClaim) error {
	// NOTE:
	// Do not move the code below to a goroutine.
	// The `ConsumeClaim` itself is called within a goroutine, see:
	// https://github.com/Shopify/sarama/blob/master/consumer_group.go#L27-L29
	for message := range claim.Messages() {
		log.Debug().
			Time("value", message.Timestamp).
			Str("value", string(message.Value)).
			Str("topic", message.Topic).
			Int32("topic", message.Partition).
			Msg("new msg")

		session.MarkMessage(message, "")

		go func(data []byte) {
			res := consumer.jrpc.HandleRPCJsonRawMessage(context.Background(), data)
			log.Debug().
				Msgf("%#v", string(res))
		}(message.Value)

	}

	return nil
}
