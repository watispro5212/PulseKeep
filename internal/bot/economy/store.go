package economy

import (
	"errors"
	"math/rand"
	"sort"
	"sync"
	"time"

	"github.com/disgoorg/snowflake/v2"
)

const (
	StartingBalance = 250
	DailyCooldown   = 24 * time.Hour
	WorkCooldown    = 45 * time.Minute
	RobCooldown     = 4 * time.Hour
)

var (
	ErrCooldown         = errors.New("command is on cooldown")
	ErrInvalidAmount    = errors.New("amount must be greater than zero")
	ErrInsufficientFund = errors.New("insufficient balance")
	ErrSelfPayment      = errors.New("you cannot pay yourself")
)

type Store struct {
	mu      sync.RWMutex
	records map[snowflake.ID]*Record
}

type Record struct {
	UserID      snowflake.ID
	Name        string
	Balance     int
	Earned      int
	Spent       int
	DailyStreak int
	LastDaily   time.Time
	LastWork    time.Time
	LastRob     time.Time
	FlipWins    int
	FlipLosses  int
	SlotWins    int
	SlotLosses  int
	RobWins     int
	RobLosses   int
	Inventory   map[string]*InventoryEntry
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

type DailyResult struct {
	Record        Record
	Reward        int
	Streak        int
	NextAvailable time.Time
	OnCooldown    bool
}

type WorkResult struct {
	Record        Record
	Reward        int
	Job           string
	NextAvailable time.Time
	OnCooldown    bool
}

type TransferResult struct {
	Sender    Record
	Receiver  Record
	Amount    int
	Recipient string
}

type FlipResult struct {
	Record Record
	Side   string
	Result string
	Wager  int
	Won    bool
}

type RobResult struct {
	Record        Record
	Target        Record
	Stolen        int
	Fine          int
	Success       bool
	NextAvailable time.Time
	OnCooldown    bool
}

type ShopItem struct {
	ID          string
	Name        string
	Description string
	Price       int
}

var ShopItems = []ShopItem{
	{ID: "lucky_pickaxe", Name: "Lucky Pickaxe", Description: "+15% coinflip win chance", Price: 5000},
	{ID: "xp_boost", Name: "XP Boost", Description: "2x work earnings for 1 hour", Price: 3000},
	{ID: "golden_watch", Name: "Golden Watch", Description: "Reduces daily cooldown by 4 hours", Price: 8000},
	{ID: "shield_token", Name: "Shield Token", Description: "Protects you from one robbery", Price: 6000},
	{ID: "lucky_clover", Name: "Lucky Clover", Description: "+1 slot reel position", Price: 4000},
}

type BuyResult struct {
	Record Record
	Item   ShopItem
}

type SlotsResult struct {
	Record     Record
	Wager      int
	Won        bool
	Multiplier int
	Symbols    [3]string
}

type InventoryEntry struct {
	ItemID   string
	ItemName string
	Quantity int
}

func NewStore() *Store {
	return &Store{
		records: make(map[snowflake.ID]*Record),
	}
}

func (s *Store) Balance(userID snowflake.ID, name string) Record {
	s.mu.Lock()
	defer s.mu.Unlock()

	return s.ensureRecord(userID, name).copy()
}

func (s *Store) Daily(userID snowflake.ID, name string, now time.Time) DailyResult {
	s.mu.Lock()
	defer s.mu.Unlock()

	record := s.ensureRecord(userID, name)
	if !record.LastDaily.IsZero() && now.Sub(record.LastDaily) < DailyCooldown {
		return DailyResult{
			Record:        record.copy(),
			NextAvailable: record.LastDaily.Add(DailyCooldown),
			OnCooldown:    true,
		}
	}

	if record.LastDaily.IsZero() || now.Sub(record.LastDaily) > 48*time.Hour {
		record.DailyStreak = 0
	}

	record.DailyStreak++
	reward := 225 + min(record.DailyStreak*25, 275)
	record.Balance += reward
	record.Earned += reward
	record.LastDaily = now
	record.UpdatedAt = now

	return DailyResult{
		Record:        record.copy(),
		Reward:        reward,
		Streak:        record.DailyStreak,
		NextAvailable: now.Add(DailyCooldown),
	}
}

func (s *Store) Work(userID snowflake.ID, name string, now time.Time) WorkResult {
	s.mu.Lock()
	defer s.mu.Unlock()

	record := s.ensureRecord(userID, name)
	if !record.LastWork.IsZero() && now.Sub(record.LastWork) < WorkCooldown {
		return WorkResult{
			Record:        record.copy(),
			NextAvailable: record.LastWork.Add(WorkCooldown),
			OnCooldown:    true,
		}
	}

	jobs := []string{
		"triaged support tickets",
		"reviewed audit logs",
		"patched a moderation workflow",
		"organized the command center",
		"cleaned up deploy notes",
		"tested the economy sandbox",
	}
	reward := 90 + rand.Intn(161)
	job := jobs[rand.Intn(len(jobs))]

	record.Balance += reward
	record.Earned += reward
	record.LastWork = now
	record.UpdatedAt = now

	return WorkResult{
		Record:        record.copy(),
		Reward:        reward,
		Job:           job,
		NextAvailable: now.Add(WorkCooldown),
	}
}

func (s *Store) Pay(senderID snowflake.ID, senderName string, receiverID snowflake.ID, receiverName string, amount int, now time.Time) (TransferResult, error) {
	if amount <= 0 {
		return TransferResult{}, ErrInvalidAmount
	}
	if senderID == receiverID {
		return TransferResult{}, ErrSelfPayment
	}

	s.mu.Lock()
	defer s.mu.Unlock()

	sender := s.ensureRecord(senderID, senderName)
	receiver := s.ensureRecord(receiverID, receiverName)

	if sender.Balance < amount {
		return TransferResult{}, ErrInsufficientFund
	}

	sender.Balance -= amount
	sender.Spent += amount
	sender.UpdatedAt = now
	receiver.Balance += amount
	receiver.Earned += amount
	receiver.UpdatedAt = now

	return TransferResult{
		Sender:    sender.copy(),
		Receiver:  receiver.copy(),
		Amount:    amount,
		Recipient: receiver.Name,
	}, nil
}

func (s *Store) Coinflip(userID snowflake.ID, name string, side string, wager int, now time.Time) (FlipResult, error) {
	if wager <= 0 {
		return FlipResult{}, ErrInvalidAmount
	}

	s.mu.Lock()
	defer s.mu.Unlock()

	record := s.ensureRecord(userID, name)
	if record.Balance < wager {
		return FlipResult{}, ErrInsufficientFund
	}

	if side != "heads" && side != "tails" {
		side = "heads"
	}

	result := "heads"
	if rand.Intn(2) == 1 {
		result = "tails"
	}

	won := side == result
	if won {
		record.Balance += wager
		record.Earned += wager
		record.FlipWins++
	} else {
		record.Balance -= wager
		record.Spent += wager
		record.FlipLosses++
	}
	record.UpdatedAt = now

	return FlipResult{
		Record: record.copy(),
		Side:   side,
		Result: result,
		Wager:  wager,
		Won:    won,
	}, nil
}

func (s *Store) Leaderboard(limit int) []Record {
	s.mu.RLock()
	defer s.mu.RUnlock()

	records := make([]Record, 0, len(s.records))
	for _, record := range s.records {
		records = append(records, record.copy())
	}

	sort.Slice(records, func(i, j int) bool {
		if records[i].Balance == records[j].Balance {
			return records[i].Earned > records[j].Earned
		}
		return records[i].Balance > records[j].Balance
	})

	if len(records) > limit {
		records = records[:limit]
	}
	return records
}

func (s *Store) Rob(userID snowflake.ID, name string, targetID snowflake.ID, targetName string, now time.Time) (RobResult, error) {
	if userID == targetID {
		return RobResult{}, ErrSelfPayment
	}

	s.mu.Lock()
	defer s.mu.Unlock()

	record := s.ensureRecord(userID, name)
	if !record.LastRob.IsZero() && now.Sub(record.LastRob) < RobCooldown {
		return RobResult{
			Record:        record.copy(),
			NextAvailable: record.LastRob.Add(RobCooldown),
			OnCooldown:    true,
		}, nil
	}

	target := s.ensureRecord(targetID, targetName)
	if target.Balance < 50 {
		return RobResult{}, ErrInsufficientFund
	}

	if record.Balance < 50 {
		return RobResult{}, ErrInsufficientFund
	}

	success := rand.Intn(100) < 40
	var stolen int
	var fine int

	if success {
		stolen = target.Balance * rand.Intn(31) / 100
		if stolen < 10 {
			stolen = 10
		}
		if stolen > target.Balance {
			stolen = target.Balance
		}
		target.Balance -= stolen
		record.Balance += stolen
		record.Earned += stolen
		target.Spent += stolen
		record.RobWins++
	} else {
		fine = record.Balance * rand.Intn(26) / 100
		if fine < 10 {
			fine = 10
		}
		if fine > record.Balance {
			fine = record.Balance
		}
		record.Balance -= fine
		record.Spent += fine
		target.Balance += fine
		target.Earned += fine
		record.RobLosses++
	}

	record.LastRob = now
	record.UpdatedAt = now

	return RobResult{
		Record:        record.copy(),
		Target:        target.copy(),
		Stolen:        stolen,
		Fine:          fine,
		Success:       success,
		NextAvailable: now.Add(RobCooldown),
	}, nil
}

func (s *Store) Shop() []ShopItem {
	return ShopItems
}

func (s *Store) Buy(userID snowflake.ID, name string, itemID string, now time.Time) (BuyResult, error) {
	var item *ShopItem
	for _, shopItem := range ShopItems {
		if shopItem.ID == itemID {
			item = &shopItem
			break
		}
	}
	if item == nil {
		return BuyResult{}, errors.New("item not found")
	}

	s.mu.Lock()
	defer s.mu.Unlock()

	record := s.ensureRecord(userID, name)
	if record.Balance < item.Price {
		return BuyResult{}, ErrInsufficientFund
	}

	record.Balance -= item.Price
	record.Spent += item.Price
	if record.Inventory == nil {
		record.Inventory = make(map[string]*InventoryEntry)
	}
	if entry, ok := record.Inventory[item.ID]; ok {
		entry.Quantity++
	} else {
		record.Inventory[item.ID] = &InventoryEntry{
			ItemID:   item.ID,
			ItemName: item.Name,
			Quantity: 1,
		}
	}
	record.UpdatedAt = now

	return BuyResult{
		Record: record.copy(),
		Item:   *item,
	}, nil
}

func (s *Store) Inventory(userID snowflake.ID, name string) []InventoryEntry {
	s.mu.RLock()
	defer s.mu.RUnlock()

	record := s.ensureRecord(userID, name)
	if record.Inventory == nil {
		return nil
	}

	entries := make([]InventoryEntry, 0, len(record.Inventory))
	for _, entry := range record.Inventory {
		entries = append(entries, *entry)
	}
	return entries
}

func (s *Store) Slots(userID snowflake.ID, name string, wager int, now time.Time) (SlotsResult, error) {
	if wager <= 0 {
		return SlotsResult{}, ErrInvalidAmount
	}

	s.mu.Lock()
	defer s.mu.Unlock()

	record := s.ensureRecord(userID, name)
	if record.Balance < wager {
		return SlotsResult{}, ErrInsufficientFund
	}

	symbols := []string{"🍒", "🍋", "🍊", "🍇", "💎", "7️⃣", "⭐"}
	reels := [3]string{
		symbols[rand.Intn(len(symbols))],
		symbols[rand.Intn(len(symbols))],
		symbols[rand.Intn(len(symbols))],
	}

	var multiplier int
	won := false

	if reels[0] == reels[1] && reels[1] == reels[2] {
		won = true
		switch reels[0] {
		case "7️⃣":
			multiplier = 10
		case "💎":
			multiplier = 7
		case "⭐":
			multiplier = 5
		case "🍇":
			multiplier = 4
		case "🍊":
			multiplier = 3
		default:
			multiplier = 2
		}
	} else if reels[0] == reels[1] || reels[1] == reels[2] || reels[0] == reels[2] {
		won = true
		multiplier = 1
	}

	if won {
		payout := wager * multiplier
		record.Balance += payout
		record.Earned += payout
		record.SlotWins++
	} else {
		record.Balance -= wager
		record.Spent += wager
		record.SlotLosses++
	}
	record.UpdatedAt = now

	return SlotsResult{
		Record:     record.copy(),
		Wager:      wager,
		Won:        won,
		Multiplier: multiplier,
		Symbols:    reels,
	}, nil
}

func (s *Store) ensureRecord(userID snowflake.ID, name string) *Record {
	if record, ok := s.records[userID]; ok {
		if name != "" {
			record.Name = name
		}
		return record
	}

	now := time.Now()
	record := &Record{
		UserID:    userID,
		Name:      name,
		Balance:   StartingBalance,
		Inventory: make(map[string]*InventoryEntry),
		CreatedAt: now,
		UpdatedAt: now,
	}
	s.records[userID] = record
	return record
}

func (r Record) NetWorth() int {
	return r.Balance
}

func (r Record) FlipTotal() int {
	return r.FlipWins + r.FlipLosses
}

func (r Record) copy() Record {
	return r
}
