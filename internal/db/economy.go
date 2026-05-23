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
		INSERT INTO user_economy (user_id, balance, total_earned) VALUES ($1, $2, $2)
		ON CONFLICT (user_id) DO UPDATE SET balance = user_economy.balance + $2, total_earned = user_economy.total_earned + $2
		RETURNING balance`, userID, amount).Scan(&balance)
	return balance, err
}

func (db *Database) RemoveBalance(ctx context.Context, userID string, amount int) error {
	_, err := db.Conn.ExecContext(ctx, `
		UPDATE user_economy SET balance = GREATEST(0, balance - $2) WHERE user_id = $1`, userID, amount)
	return err
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
		INSERT INTO user_economy (user_id, balance, last_daily_claim, total_earned) VALUES ($1, $2, $3, $2)
		ON CONFLICT (user_id) DO UPDATE SET balance = user_economy.balance + $2, last_daily_claim = $3, total_earned = user_economy.total_earned + $2
		RETURNING balance`, userID, amount, now).Scan(&balance)

	return balance, true, err
}

func (db *Database) GetNextDailyTime(ctx context.Context, userID string) (time.Time, error) {
	var lastClaim *time.Time
	err := db.Conn.QueryRowContext(ctx, "SELECT last_daily_claim FROM user_economy WHERE user_id = $1", userID).Scan(&lastClaim)
	if err != nil || lastClaim == nil {
		return time.Now(), nil
	}
	return lastClaim.Add(24 * time.Hour), nil
}

func (db *Database) CanWork(ctx context.Context, userID string) (bool, time.Duration) {
	var lastWork *time.Time
	err := db.Conn.QueryRowContext(ctx, "SELECT last_work FROM user_economy WHERE user_id = $1", userID).Scan(&lastWork)
	if err != nil || lastWork == nil {
		return true, 0
	}
	
	remaining := 5*time.Minute - time.Since(*lastWork)
	if remaining <= 0 {
		return true, 0
	}
	return false, remaining
}

func (db *Database) SetWorkCooldown(ctx context.Context, userID string) error {
	_, err := db.Conn.ExecContext(ctx, `
		INSERT INTO user_economy (user_id, last_work) VALUES ($1, $2)
		ON CONFLICT (user_id) DO UPDATE SET last_work = $2`, userID, time.Now())
	return err
}

func (db *Database) CanRob(ctx context.Context, userID string) (bool, time.Duration) {
	var lastRob *time.Time
	err := db.Conn.QueryRowContext(ctx, "SELECT last_rob FROM user_economy WHERE user_id = $1", userID).Scan(&lastRob)
	if err != nil || lastRob == nil {
		return true, 0
	}
	
	remaining := 15*time.Minute - time.Since(*lastRob)
	if remaining <= 0 {
		return true, 0
	}
	return false, remaining
}

func (db *Database) SetRobCooldown(ctx context.Context, userID string) error {
	_, err := db.Conn.ExecContext(ctx, `
		INSERT INTO user_economy (user_id, last_rob) VALUES ($1, $2)
		ON CONFLICT (user_id) DO UPDATE SET last_rob = $2`, userID, time.Now())
	return err
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

	_, err = tx.ExecContext(ctx, "UPDATE user_economy SET balance = balance - $1, transactions = transactions + 1 WHERE user_id = $2", amount, fromID)
	if err != nil {
		return err
	}

	_, err = tx.ExecContext(ctx, `
		INSERT INTO user_economy (user_id, balance, transactions) VALUES ($1, $2, 1)
		ON CONFLICT (user_id) DO UPDATE SET balance = user_economy.balance + $2, transactions = user_economy.transactions + 1`, toID, amount)
	if err != nil {
		return err
	}

	return tx.Commit()
}

func (db *Database) GetUserRank(ctx context.Context, userID string) (int, error) {
	var balance int
	err := db.Conn.QueryRowContext(ctx, "SELECT balance FROM user_economy WHERE user_id = $1", userID).Scan(&balance)
	if err != nil {
		return 0, err
	}

	var rank int
	err = db.Conn.QueryRowContext(ctx, "SELECT COUNT(*) + 1 FROM user_economy WHERE balance > $1", balance).Scan(&rank)
	return rank, err
}

func (db *Database) GetLeaderboard(ctx context.Context, limit, offset int) ([]LeaderboardEntry, error) {
	rows, err := db.Conn.QueryContext(ctx, `
		SELECT user_id, balance FROM user_economy 
		ORDER BY balance DESC 
		LIMIT $1 OFFSET $2`, limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var entries []LeaderboardEntry
	for rows.Next() {
		var entry LeaderboardEntry
		if err := rows.Scan(&entry.UserID, &entry.Balance); err != nil {
			return nil, err
		}
		entries = append(entries, entry)
	}
	return entries, rows.Err()
}

func (db *Database) GetUserStats(ctx context.Context, userID string) (UserStats, error) {
	var stats UserStats
	var totalEarned, totalGambled, transactions sql.NullInt64

	err := db.Conn.QueryRowContext(ctx, `
		SELECT total_earned, COALESCE(total_gambled, 0), COALESCE(transactions, 0) 
		FROM user_economy WHERE user_id = $1`, userID).Scan(&totalEarned, &totalGambled, &transactions)
	
	if err == sql.ErrNoRows {
		return stats, nil
	}
	if err != nil {
		return stats, err
	}

	if totalEarned.Valid {
		stats.TotalEarned = int(totalEarned.Int64)
	}
	stats.TotalGambled = int(totalGambled.Int64)
	stats.Transactions = int(transactions.Int64)

	return stats, nil
}

func (db *Database) AddItem(ctx context.Context, userID, itemID, itemName string, quantity int) error {
	_, err := db.Conn.ExecContext(ctx, `
		INSERT INTO user_inventory (user_id, item_id, item_name, quantity) VALUES ($1, $2, $3, $4)
		ON CONFLICT (user_id, item_id) DO UPDATE SET quantity = user_inventory.quantity + $4`, userID, itemID, itemName, quantity)
	return err
}

func (db *Database) GetInventory(ctx context.Context, userID string) ([]InventoryItem, error) {
	rows, err := db.Conn.QueryContext(ctx, `
		SELECT item_name, quantity FROM user_inventory WHERE user_id = $1 AND quantity > 0`, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var items []InventoryItem
	for rows.Next() {
		var item InventoryItem
		if err := rows.Scan(&item.Name, &item.Quantity); err != nil {
			return nil, err
		}
		items = append(items, item)
	}
	return items, rows.Err()
}

// LeaderboardEntry represents a user on the leaderboard
type LeaderboardEntry struct {
	UserID  string
	Balance int
}

// UserStats represents a user's economy statistics
type UserStats struct {
	TotalEarned  int
	TotalGambled int
	Transactions int
}

// InventoryItem represents an item in a user's inventory
type InventoryItem struct {
	Name     string
	Quantity int
}
