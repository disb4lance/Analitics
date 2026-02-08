package service

import (
	"analytics-service/internal/domain"
	"context"
	"time"

	"github.com/shopspring/decimal"
)

type TransactionCreatedCommand struct {
	UserID       string
	CategoryName string
	Amount       decimal.Decimal
	CreatedAt    time.Time
}

type UserStatRepository interface {
	Add(
		ctx context.Context,
		stat *domain.UserStat,
	) error
	Update(
		ctx context.Context,
		stat *domain.UserStat,
	) error
	GetByUserID(
		ctx context.Context,
		userID string,
	) (*domain.UserStat, error)
}

type AnalyticsService struct {
	repo UserStatRepository
}

func NewAnalyticsService(
	repo UserStatRepository,
) *AnalyticsService {
	return &AnalyticsService{repo: repo}
}

func (s *AnalyticsService) HandleTransactionCreated(
	ctx context.Context,
	cmd TransactionCreatedCommand,
) error {

	stat, err := s.repo.GetByUserID(ctx, cmd.UserID)
	if err != nil {
		return err
	}

	if stat == nil {
		stat = domain.NewUserStat(
			cmd.UserID,
			cmd.CategoryName,
			cmd.Amount,
			cmd.CreatedAt,
		)

		return s.repo.Add(ctx, stat)
	}

	stat.ApplyTransaction(
		cmd.CategoryName,
		cmd.Amount,
		cmd.CreatedAt,
	)

	return s.repo.Update(ctx, stat)
}
