package economy

import (
	"errors"
	"testing"
	"time"

	"github.com/disgoorg/snowflake/v2"
)

func TestDailyRewardAndCooldown(t *testing.T) {
	store := NewStore()
	now := time.Date(2026, 5, 23, 12, 0, 0, 0, time.UTC)
	userID := snowflake.ID(1001)

	first := store.Daily(userID, "Ada", now)
	if first.OnCooldown {
		t.Fatal("first daily claim should not be on cooldown")
	}
	if first.Reward != 250 {
		t.Fatalf("expected first daily reward 250, got %d", first.Reward)
	}
	if first.Record.Balance != StartingBalance+250 {
		t.Fatalf("unexpected balance after daily: %d", first.Record.Balance)
	}

	second := store.Daily(userID, "Ada", now.Add(time.Hour))
	if !second.OnCooldown {
		t.Fatal("second daily claim inside cooldown should be blocked")
	}
}

func TestPayTransfersBalance(t *testing.T) {
	store := NewStore()
	now := time.Date(2026, 5, 23, 12, 0, 0, 0, time.UTC)

	result, err := store.Pay(snowflake.ID(1), "Sender", snowflake.ID(2), "Receiver", 125, now)
	if err != nil {
		t.Fatalf("pay returned error: %v", err)
	}
	if result.Sender.Balance != StartingBalance-125 {
		t.Fatalf("unexpected sender balance: %d", result.Sender.Balance)
	}
	if result.Receiver.Balance != StartingBalance+125 {
		t.Fatalf("unexpected receiver balance: %d", result.Receiver.Balance)
	}
}

func TestPayRejectsBadTransfers(t *testing.T) {
	store := NewStore()
	now := time.Now()

	_, err := store.Pay(snowflake.ID(1), "Sender", snowflake.ID(1), "Sender", 10, now)
	if !errors.Is(err, ErrSelfPayment) {
		t.Fatalf("expected self payment error, got %v", err)
	}

	_, err = store.Pay(snowflake.ID(1), "Sender", snowflake.ID(2), "Receiver", StartingBalance+1, now)
	if !errors.Is(err, ErrInsufficientFund) {
		t.Fatalf("expected insufficient fund error, got %v", err)
	}
}

func TestLeaderboardSortsByBalance(t *testing.T) {
	store := NewStore()
	now := time.Now()

	if _, err := store.Pay(snowflake.ID(1), "First", snowflake.ID(2), "Second", 50, now); err != nil {
		t.Fatalf("pay returned error: %v", err)
	}
	records := store.Leaderboard(2)
	if len(records) != 2 {
		t.Fatalf("expected 2 records, got %d", len(records))
	}
	if records[0].Name != "Second" {
		t.Fatalf("expected Second on top, got %s", records[0].Name)
	}
}
