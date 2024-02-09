package kafka

import (
	"context"
	"time"

	"github.com/invest-scraping/config"
	"github.com/invest-scraping/logg"
	"github.com/invest-scraping/stream"
	"github.com/segmentio/kafka-go"
)

type Producer struct {
	writer kafka.Writer
	log    logg.Logger
}

func NewProducer(cfg *config.Stream, topic string, log logg.Logger) stream.Producer {
	cli := Producer{log: log}
	cli.Connect(cfg, topic)
	return &cli
}

func (c *Producer) Connect(cfg *config.Stream, topic string) {
	c.log.Info("Connecting to kafka broker: ", cfg.Kafka.Brokers)

	w := kafka.Writer{
		Addr:                   kafka.TCP(cfg.Kafka.Brokers...),
		Topic:                  topic,
		BatchTimeout:           time.Millisecond * 1,
		AllowAutoTopicCreation: true,
	}

	c.writer = w
}

func (c *Producer) Produce(msg []byte) error {
	m := kafka.Message{
		Value: msg,
	}

	if err := c.writer.WriteMessages(context.Background(), m); err != nil {
		c.log.Info("Failed to produce message to kafka. Reason: ", err.Error())
		return err
	}

	return nil
}
