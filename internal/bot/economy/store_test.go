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

func TestRobRequiresDifferentUsers(t *testing.T) {
	store := NewStore()
	now := time.Now()

	_, err := store.Rob(snowflake.ID(1), "User", snowflake.ID(1), "User", now)
	if !errors.Is(err, ErrSelfPayment) {
		t.Fatalf("expected self-rob error, got %v", err)
	}
}

func TestSlotsAdjustsBalance(t *testing.T) {
	store := NewStore()
	now := time.Now()

	result, err := store.Slots(snowflake.ID(1), "Player", 50, now)
	if err != nil {
		t.Fatalf("slots returned error: %v", err)
	}
	if result.Wager != 50 {
		t.Fatalf("expected wager 50, got %d", result.Wager)
	}
	// Balance should have changed (either won or lost)
	if result.Record.Balance == StartingBalance {
		t.Fatal("expected balance to change after slots")
	}
}

func TestBuyDeductsBalance(t *testing.T) {
	store := NewStore()
	now := time.Now()

	// make the buyer rich with repeated daily claims
	buyerID := snowflake.ID(1)
	for i := 0; i < 20; i++ {
		store.Daily(buyerID, "Buyer", now.Add(time.Duration(i)*25*time.Hour))
	}

	result, err := store.Buy(buyerID, "Buyer", "lucky_pickaxe", now.Add(500*time.Hour))
	if err != nil {
		t.Fatalf("buy returned error: %v", err)
	}
	if result.Item.ID != "lucky_pickaxe" {
		t.Fatalf("expected lucky_pickaxe, got %s", result.Item.ID)
	}
	if result.Record.Balance <= 0 {
		t.Fatalf("expected positive balance after buy, got %d", result.Record.Balance)
	}
}

func TestBuyRejectsUnknownItem(t *testing.T) {
	store := NewStore()
	now := time.Now()

	_, err := store.Buy(snowflake.ID(1), "Buyer", "nonexistent", now)
	if err == nil || err.Error() != "item not found" {
		t.Fatalf("expected item not found error, got %v", err)
	}
}

func TestBuyRejectsInsufficientFunds(t *testing.T) {
	store := NewStore()
	now := time.Now()

	_, err := store.Buy(snowflake.ID(1), "Buyer", "golden_watch", now)
	if !errors.Is(err, ErrInsufficientFund) {
		t.Fatalf("expected insufficient fund error, got %v", err)
	}
}

func TestInventoryAfterBuy(t *testing.T) {
	store := NewStore()
	now := time.Now()

	buyerID := snowflake.ID(1)
	for i := 0; i < 20; i++ {
		store.Daily(buyerID, "Buyer", now.Add(time.Duration(i)*25*time.Hour))
	}

	_, err := store.Buy(buyerID, "Buyer", "lucky_pickaxe", now.Add(500*time.Hour))
	if err != nil {
		t.Fatalf("buy returned error: %v", err)
	}

	items := store.Inventory(buyerID, "Buyer")
	if len(items) != 1 {
		t.Fatalf("expected 1 item, got %d", len(items))
	}
	if items[0].ItemID != "lucky_pickaxe" {
		t.Fatalf("expected lucky_pickaxe, got %s", items[0].ItemID)
	}
	if items[0].Quantity != 1 {
		t.Fatalf("expected quantity 1, got %d", items[0].Quantity)
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
