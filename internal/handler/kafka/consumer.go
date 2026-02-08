package kafka

import (
	"analytics-service/internal/service"
	"context"
	"encoding/json"
	"time"

	"github.com/shopspring/decimal"
)

type TransactionCreatedEvent struct {
	UserID       string          `json:"user_id"`
	CategoryName string          `json:"category_name"`
	Amount       decimal.Decimal `json:"amount"`
	CreatedAt    time.Time       `json:"created_at"`
}

type AnalyticsService interface {
	HandleTransactionCreated(
		ctx context.Context,
		cmd service.TransactionCreatedCommand,
	) error
}

type Consumer struct {
	service AnalyticsService
}

func NewConsumer(
	service AnalyticsService,
) *Consumer {
	return &Consumer{service: service}
}

func (c *Consumer) HandleMessage(
	ctx context.Context,
	message []byte,
) error {

	var event TransactionCreatedEvent

	if err := json.Unmarshal(message, &event); err != nil {
		return err
	}

	cmd := service.TransactionCreatedCommand{
		UserID:       event.UserID,
		CategoryName: event.CategoryName,
		Amount:       event.Amount,
		CreatedAt:    event.CreatedAt,
	}

	return c.service.HandleTransactionCreated(ctx, cmd)
}
