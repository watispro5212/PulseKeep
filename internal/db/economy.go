package db

import (
	"context"
	"database/sql"
	"time"
)

func (db *Database) GetBalance(ctx context.Context, userID string) (int, error) {
	var balance int
	err := db.Conn.QueryRowContext(ctx, "SELECT balance FROM user_economy WHERE user_id = $1", userID).Scan(&balance)
	if err == sql.ErrNoRows {
		_, err = db.Conn.ExecContext(ctx, "INSERT INTO user_economy (user_id, balance) VALUES ($1, 0)", userID)
		return 0, err
	}
	return balance, err
}

func (db *Database) AddBalance(ctx context.Context, userID string, amount int) (int, error) {
	var balance int
	err := db.Conn.QueryRowContext(ctx, `
		INSERT INTO user_economy (user_id, balance) VALUES ($1, $2)
		ON CONFLICT (user_id) DO UPDATE SET balance = user_economy.balance + $2
		RETURNING balance`, userID, amount).Scan(&balance)
	return balance, err
}

func (db *Database) ClaimDaily(ctx context.Context, userID string, amount int) (int, bool, error) {
	var lastClaim *time.Time
	err := db.Conn.QueryRowContext(ctx, "SELECT last_daily_claim FROM user_economy WHERE user_id = $1", userID).Scan(&lastClaim)

	if err == nil && lastClaim != nil {
		if time.Since(*lastClaim) < 24*time.Hour {
			return 0, false, nil
		}
	}

	var balance int
	now := time.Now()
	err = db.Conn.QueryRowContext(ctx, `
		INSERT INTO user_economy (user_id, balance, last_daily_claim) VALUES ($1, $2, $3)
		ON CONFLICT (user_id) DO UPDATE SET balance = user_economy.balance + $2, last_daily_claim = $3
		RETURNING balance`, userID, amount, now).Scan(&balance)

	return balance, true, err
}

func (db *Database) Transfer(ctx context.Context, fromID, toID string, amount int) error {
	tx, err := db.Conn.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()

	var fromBalance int
	err = tx.QueryRowContext(ctx, "SELECT balance FROM user_economy WHERE user_id = $1 FOR UPDATE", fromID).Scan(&fromBalance)
	if err != nil {
		return err
	}

	if fromBalance < amount {
		return sql.ErrNoRows
	}

	_, err = tx.ExecContext(ctx, "UPDATE user_economy SET balance = balance - $1 WHERE user_id = $2", amount, fromID)
	if err != nil {
		return err
	}

	_, err = tx.ExecContext(ctx, `
		INSERT INTO user_economy (user_id, balance) VALUES ($1, $2)
		ON CONFLICT (user_id) DO UPDATE SET balance = user_economy.balance + $2`, toID, amount)
	if err != nil {
		return err
	}

	return tx.Commit()
}
