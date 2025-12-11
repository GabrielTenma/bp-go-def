package infrastructure

import (
	"context"
	"fmt"
	"log"
	"test-go/config"

	"github.com/IBM/sarama"
)

type KafkaManager struct {
	Producer sarama.SyncProducer
	Brokers  []string
	GroupID  string
}

func NewKafkaManager(cfg config.KafkaConfig) (*KafkaManager, error) {
	if !cfg.Enabled {
		return nil, nil
	}

	config := sarama.NewConfig()
	config.Producer.Return.Successes = true
	config.Producer.RequiredAcks = sarama.WaitForAll
	config.Producer.Retry.Max = 5

	producer, err := sarama.NewSyncProducer(cfg.Brokers, config)
	if err != nil {
		return nil, fmt.Errorf("failed to start kafka producer: %w", err)
	}

	return &KafkaManager{
		Producer: producer,
		Brokers:  cfg.Brokers,
		GroupID:  cfg.GroupID,
	}, nil
}

func (k *KafkaManager) GetStatus() map[string]interface{} {
	stats := make(map[string]interface{})
	if k == nil {
		stats["connected"] = false
		return stats
	}

	if k.Producer == nil && len(k.Brokers) == 0 {

		stats["connected"] = false
		return stats
	}

	stats["connected"] = true // Assuming connected if initialized for now, complex to check liveness without producing
	stats["brokers"] = k.Brokers
	stats["group_id"] = k.GroupID
	return stats
}

// Consume starts a consumer group for the given topic.
// NOTE: This blocks the calling goroutine. Run in a separate goroutine.
func (k *KafkaManager) Consume(ctx context.Context, topic string, handler func(key, value []byte) error) error {
	config := sarama.NewConfig()
	config.Consumer.Group.Rebalance.Strategy = sarama.BalanceStrategyRoundRobin
	config.Consumer.Offsets.Initial = sarama.OffsetOldest

	consumerGroup, err := sarama.NewConsumerGroup(k.Brokers, k.GroupID, config)
	if err != nil {
		return fmt.Errorf("error creating consumer group: %w", err)
	}
	defer consumerGroup.Close()

	consumer := &consumerHandler{
		handler: handler,
	}

	for {
		select {
		case <-ctx.Done():
			return nil
		default:
			if err := consumerGroup.Consume(ctx, []string{topic}, consumer); err != nil {
				return fmt.Errorf("error from consumer: %w", err)
			}
		}
	}
}

// consumerHandler implements sarama.ConsumerGroupHandler
type consumerHandler struct {
	handler func(key, value []byte) error
}

func (h *consumerHandler) Setup(sarama.ConsumerGroupSession) error   { return nil }
func (h *consumerHandler) Cleanup(sarama.ConsumerGroupSession) error { return nil }
func (h *consumerHandler) ConsumeClaim(session sarama.ConsumerGroupSession, claim sarama.ConsumerGroupClaim) error {
	for message := range claim.Messages() {
		if err := h.handler(message.Key, message.Value); err != nil {
			log.Printf("Error handling message: %v", err)
		}
		session.MarkMessage(message, "")
	}
	return nil
}
