package clickhouse

import (
	"context"
	"database/sql"

	"analytics-service/internal/domain"
)

type UserStatRepository struct {
	db *sql.DB
}

func NewUserStatRepository(db *sql.DB) *UserStatRepository {
	return &UserStatRepository{db: db}
}

func (r *UserStatRepository) Add(
	ctx context.Context,
	stat *domain.UserStat,
) error {

	const query = `
		INSERT INTO user_stats (
			user_id,
			first_tx_at,
			last_tx_at,
			total_spent,
			avg_tx_amount,
			tx_count,
			top_category,
			top_category_amount,
			updated_at
		) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?)
	`

	_, err := r.db.ExecContext(
		ctx,
		query,
		stat.UserID,
		stat.FirstTxAt,
		stat.LastTxAt,
		stat.TotalSpent,
		stat.AvgTxAmount,
		stat.TxCount,
		stat.TopCategory,
		stat.TopCategoryAmount,
		stat.UpdatedAt,
	)

	return err
}

func (r *UserStatRepository) Update(
	ctx context.Context,
	stat *domain.UserStat,
) error {

	// В ClickHouse update = insert
	return r.Add(ctx, stat)
}

func (r *UserStatRepository) GetByUserID(
	ctx context.Context,
	userID string,
) (*domain.UserStat, error) {

	const query = `
		SELECT
			user_id,
			first_tx_at,
			last_tx_at,
			total_spent,
			avg_tx_amount,
			tx_count,
			top_category,
			top_category_amount,
			updated_at
		FROM user_stats
		WHERE user_id = ?
		ORDER BY updated_at DESC
		LIMIT 1
	`

	row := r.db.QueryRowContext(ctx, query, userID)

	var stat domain.UserStat

	err := row.Scan(
		&stat.UserID,
		&stat.FirstTxAt,
		&stat.LastTxAt,
		&stat.TotalSpent,
		&stat.AvgTxAmount,
		&stat.TxCount,
		&stat.TopCategory,
		&stat.TopCategoryAmount,
		&stat.UpdatedAt,
	)

	if err == sql.ErrNoRows {
		return nil, nil
	}

	if err != nil {
		return nil, err
	}

	return &stat, nil
}
