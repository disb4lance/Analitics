package kafka

import (
	"context"
	"log"

	"github.com/segmentio/kafka-go"
)

type Reader struct {
	reader *kafka.Reader
}

func NewKafkaReader(
	broker string,
	topic string,
	groupID string,
) *Reader {

	return &Reader{
		reader: kafka.NewReader(kafka.ReaderConfig{
			Brokers: []string{broker},
			Topic:   topic,
			GroupID: groupID,
		}),
	}
}

func (r *Reader) Start(
	ctx context.Context,
	consumer *Consumer,
) error {

	for {
		msg, err := r.reader.ReadMessage(ctx)
		if err != nil {
			return err
		}

		if err := consumer.HandleMessage(ctx, msg.Value); err != nil {
			log.Println("handler error:", err)
		}
	}
}
