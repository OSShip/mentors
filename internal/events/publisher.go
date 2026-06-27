package events

import (
	"context"

	"github.com/OSShip/utils/kafka"
)

type Publisher struct {
	producer *kafka.Producer
}

func New(brokers string) *Publisher {
	return &Publisher{producer: kafka.NewProducer(brokers, "mentor.events")}
}

func (p *Publisher) Close() {
	p.producer.Close()
}

func (p *Publisher) PublishApproved(ctx context.Context, userID, mentorEmail string) error {
	return p.producer.Publish(ctx, "mentor.approved", map[string]string{
		"user_id":      userID,
		"mentor_email": mentorEmail,
	})
}
