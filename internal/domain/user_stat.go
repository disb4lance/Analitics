package domain

import (
	"time"

	"github.com/shopspring/decimal"
)

type UserStat struct {
	UserID            string          `ch:"user_id"`
	FirstTxAt         time.Time       `ch:"first_tx_at"`
	LastTxAt          time.Time       `ch:"last_tx_at"`
	TotalSpent        decimal.Decimal `ch:"total_spent"`
	AvgTxAmount       decimal.Decimal `ch:"avg_tx_amount"`
	TxCount           uint32          `ch:"tx_count"`
	TopCategory       string          `ch:"top_category"`
	TopCategoryAmount decimal.Decimal `ch:"top_category_amount"`
	UpdatedAt         time.Time       `ch:"updated_at"`
}

func NewUserStat(
	userID string,
	category string,
	amount decimal.Decimal,
	createdAt time.Time,
) *UserStat {

	return &UserStat{
		UserID:            userID,
		FirstTxAt:         createdAt,
		LastTxAt:          createdAt,
		TotalSpent:        amount,
		AvgTxAmount:       amount,
		TxCount:           1,
		TopCategory:       category,
		TopCategoryAmount: amount,
		UpdatedAt:         time.Now(),
	}
}

func (u *UserStat) ApplyTransaction(
	category string,
	amount decimal.Decimal,
	createdAt time.Time,
) {
	u.LastTxAt = createdAt
	u.TxCount++

	u.TotalSpent = u.TotalSpent.Add(amount)
	u.AvgTxAmount = u.TotalSpent.Div(
		decimal.NewFromInt(int64(u.TxCount)),
	)

	u.applyTopCategory(category, amount)
	u.UpdatedAt = time.Now()
}

func (u *UserStat) applyTopCategory(
	category string,
	amount decimal.Decimal,
) {
	if category == u.TopCategory {
		u.TopCategoryAmount = u.TopCategoryAmount.Add(amount)
		return
	}

	if amount.GreaterThan(u.TopCategoryAmount) {
		u.TopCategory = category
		u.TopCategoryAmount = amount
	}
}
