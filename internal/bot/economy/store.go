package economy

import (
	"errors"
	"fmt"
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
	FishingCooldown = 45 * time.Second
	MiningCooldown  = 45 * time.Second
	GambleCooldown  = 10 * time.Second
	InterestRate    = 0.001
	InterestInterval = 6 * time.Hour
)

var (
	ErrCooldown         = errors.New("command is on cooldown")
	ErrInvalidAmount    = errors.New("amount must be greater than zero")
	ErrInsufficientFund = errors.New("insufficient balance")
	ErrSelfPayment      = errors.New("you cannot pay yourself")
	ErrItemNotFound     = errors.New("item not found")
	ErrNotOwned         = errors.New("you do not own this item")
	ErrCannotUse        = errors.New("this item cannot be used")
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
	LastFish    time.Time
	LastMine    time.Time
	LastGamble  time.Time
	LastInterest time.Time
	FlipWins    int
	FlipLosses  int
	SlotWins    int
	SlotLosses  int
	RobWins     int
	RobLosses   int
	FishCaught  int
	FishTotal   int
	MineMined   int
	MineTotal   int
	GambleWins  int
	GambleTotal int
	LotteryWins int
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
	{ID: "shield_token", Name: "Shield Token", Description: "Protects you from one robbery (consumed on use)", Price: 6000},
	{ID: "lucky_clover", Name: "Lucky Clover", Description: "+1 slot reel position", Price: 4000},
	{ID: "fishing_rod", Name: "Fishing Rod", Description: "Unlocks the /fish command", Price: 1500},
	{ID: "iron_pickaxe", Name: "Iron Pickaxe", Description: "Unlocks the /mine command", Price: 2000},
	{ID: "lottery_ticket", Name: "Lottery Ticket", Description: "Entry for the server lottery draw", Price: 500},
	{ID: "health_potion", Name: "Health Potion", Description: "Restores 25% of lost rob fines", Price: 2500},
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

type FishType struct {
	Name   string
	Emoji  string
	Weight string
	MinPay int
	MaxPay int
	Rarity string
}

var FishTable = []FishType{
	{Name: "Minnow", Emoji: "🐟", Weight: "0.1 lbs", MinPay: 5, MaxPay: 15, Rarity: "Common"},
	{Name: "Perch", Emoji: "🐟", Weight: "0.5 lbs", MinPay: 10, MaxPay: 25, Rarity: "Common"},
	{Name: "Bass", Emoji: "🐠", Weight: "2 lbs", MinPay: 20, MaxPay: 50, Rarity: "Uncommon"},
	{Name: "Trout", Emoji: "🐟", Weight: "3 lbs", MinPay: 30, MaxPay: 70, Rarity: "Uncommon"},
	{Name: "Salmon", Emoji: "🐟", Weight: "5 lbs", MinPay: 50, MaxPay: 120, Rarity: "Rare"},
	{Name: "Tuna", Emoji: "🐟", Weight: "20 lbs", MinPay: 80, MaxPay: 200, Rarity: "Rare"},
	{Name: "Catfish", Emoji: "🐱", Weight: "8 lbs", MinPay: 60, MaxPay: 150, Rarity: "Rare"},
	{Name: "Swordfish", Emoji: "🗡️", Weight: "50 lbs", MinPay: 150, MaxPay: 350, Rarity: "Epic"},
	{Name: "Golden Koi", Emoji: "✨", Weight: "3 lbs", MinPay: 300, MaxPay: 800, Rarity: "Legendary"},
	{Name: "Ancient Coelacanth", Emoji: "🦕", Weight: "80 lbs", MinPay: 500, MaxPay: 1500, Rarity: "Mythic"},
	{Name: "Boot", Emoji: "👢", Weight: "1 lb", MinPay: 0, MaxPay: 5, Rarity: "Junk"},
	{Name: "Old Tire", Emoji: "🔘", Weight: "15 lbs", MinPay: 0, MaxPay: 10, Rarity: "Junk"},
	{Name: "Seaweed", Emoji: "🌿", Weight: "0.2 lbs", MinPay: 2, MaxPay: 8, Rarity: "Junk"},
}

type FishResult struct {
	Record        Record
	Fish          FishType
	Reward        int
	NextAvailable time.Time
	OnCooldown    bool
}

type OreType struct {
	Name   string
	Emoji  string
	MinPay int
	MaxPay int
	Rarity string
}

var OreTable = []OreType{
	{Name: "Coal", Emoji: "⚫", MinPay: 5, MaxPay: 15, Rarity: "Common"},
	{Name: "Copper", Emoji: "🟤", MinPay: 10, MaxPay: 25, Rarity: "Common"},
	{Name: "Iron", Emoji: "🔩", MinPay: 15, MaxPay: 35, Rarity: "Common"},
	{Name: "Silver", Emoji: "🥈", MinPay: 30, MaxPay: 60, Rarity: "Uncommon"},
	{Name: "Gold", Emoji: "🥇", MinPay: 50, MaxPay: 100, Rarity: "Uncommon"},
	{Name: "Platinum", Emoji: "💎", MinPay: 80, MaxPay: 180, Rarity: "Rare"},
	{Name: "Ruby", Emoji: "🔴", MinPay: 120, MaxPay: 250, Rarity: "Rare"},
	{Name: "Sapphire", Emoji: "🔵", MinPay: 150, MaxPay: 300, Rarity: "Epic"},
	{Name: "Emerald", Emoji: "🟢", MinPay: 200, MaxPay: 400, Rarity: "Epic"},
	{Name: "Diamond", Emoji: "💠", MinPay: 300, MaxPay: 600, Rarity: "Legendary"},
	{Name: "Ancient Relic", Emoji: "🏺", MinPay: 500, MaxPay: 1200, Rarity: "Mythic"},
	{Name: "Stone", Emoji: "🪨", MinPay: 1, MaxPay: 5, Rarity: "Junk"},
	{Name: "Gravel", Emoji: "🪨", MinPay: 1, MaxPay: 3, Rarity: "Junk"},
}

type MineResult struct {
	Record        Record
	Ore           OreType
	Reward        int
	NextAvailable time.Time
	OnCooldown    bool
}

type GambleResult struct {
	Record        Record
	Roll          int
	Wager         int
	Won           bool
	Multiplier    int
	NextAvailable time.Time
	OnCooldown    bool
}

type SellResult struct {
	Record  Record
	Item    InventoryEntry
	Reward  int
}

type UseItemResult struct {
	Record      Record
	ItemID      string
	ItemName    string
	Used        bool
	Description string
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

	// check if target has shield token
	shieldActive := false
	if target.Inventory != nil {
		if entry, ok := target.Inventory["shield_token"]; ok && entry.Quantity > 0 {
			shieldActive = true
			entry.Quantity--
			if entry.Quantity <= 0 {
				delete(target.Inventory, "shield_token")
			}
		}
	}

	if shieldActive {
		fine := record.Balance * 15 / 100
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
		record.LastRob = now
		record.UpdatedAt = now
		return RobResult{
			Record:        record.copy(),
			Target:        target.copy(),
			Fine:          fine,
			NextAvailable: now.Add(RobCooldown),
		}, nil
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

func (s *Store) Fish(userID snowflake.ID, name string, now time.Time) FishResult {
	s.mu.Lock()
	defer s.mu.Unlock()

	record := s.ensureRecord(userID, name)
	if !record.LastFish.IsZero() && now.Sub(record.LastFish) < FishingCooldown {
		return FishResult{
			Record:        record.copy(),
			NextAvailable: record.LastFish.Add(FishingCooldown),
			OnCooldown:    true,
		}
	}

	hasRod := false
	if record.Inventory != nil {
		if _, ok := record.Inventory["fishing_rod"]; ok {
			hasRod = true
		}
	}

	if !hasRod {
		fish := FishTable[len(FishTable)-1]
		record.FishTotal++
		record.LastFish = now
		record.UpdatedAt = now
		return FishResult{
			Record:        record.copy(),
			Fish:          fish,
			Reward:        0,
			NextAvailable: now.Add(FishingCooldown),
		}
	}

	idx := rand.Intn(len(FishTable))
	fish := FishTable[idx]
	reward := fish.MinPay + rand.Intn(fish.MaxPay-fish.MinPay+1)

	if reward > 0 {
		record.Balance += reward
		record.Earned += reward
		record.FishCaught++
	}
	record.FishTotal++
	record.LastFish = now
	record.UpdatedAt = now

	return FishResult{
		Record:        record.copy(),
		Fish:          fish,
		Reward:        reward,
		NextAvailable: now.Add(FishingCooldown),
	}
}

func (s *Store) Mine(userID snowflake.ID, name string, now time.Time) MineResult {
	s.mu.Lock()
	defer s.mu.Unlock()

	record := s.ensureRecord(userID, name)
	if !record.LastMine.IsZero() && now.Sub(record.LastMine) < MiningCooldown {
		return MineResult{
			Record:        record.copy(),
			NextAvailable: record.LastMine.Add(MiningCooldown),
			OnCooldown:    true,
		}
	}

	hasPickaxe := false
	if record.Inventory != nil {
		if _, ok := record.Inventory["iron_pickaxe"]; ok {
			hasPickaxe = true
		}
	}

	if !hasPickaxe {
		ore := OreTable[len(OreTable)-1]
		record.MineTotal++
		record.LastMine = now
		record.UpdatedAt = now
		return MineResult{
			Record:        record.copy(),
			Ore:           ore,
			Reward:        0,
			NextAvailable: now.Add(MiningCooldown),
		}
	}

	idx := rand.Intn(len(OreTable))
	ore := OreTable[idx]
	reward := ore.MinPay + rand.Intn(ore.MaxPay-ore.MinPay+1)

	if reward > 0 {
		record.Balance += reward
		record.Earned += reward
		record.MineMined++
	}
	record.MineTotal++
	record.LastMine = now
	record.UpdatedAt = now

	return MineResult{
		Record:        record.copy(),
		Ore:           ore,
		Reward:        reward,
		NextAvailable: now.Add(MiningCooldown),
	}
}

func (s *Store) Gamble(userID snowflake.ID, name string, wager int, now time.Time) (GambleResult, error) {
	if wager <= 0 {
		return GambleResult{}, ErrInvalidAmount
	}

	s.mu.Lock()
	defer s.mu.Unlock()

	record := s.ensureRecord(userID, name)
	if record.Balance < wager {
		return GambleResult{}, ErrInsufficientFund
	}
	if !record.LastGamble.IsZero() && now.Sub(record.LastGamble) < GambleCooldown {
		return GambleResult{
			Record:        record.copy(),
			OnCooldown:    true,
			NextAvailable: record.LastGamble.Add(GambleCooldown),
		}, nil
	}

	roll := 1 + rand.Intn(100)
	var won bool
	multiplier := 0

	switch {
	case roll == 100:
		won = true
		multiplier = 10
	case roll >= 95:
		won = true
		multiplier = 4
	case roll >= 85:
		won = true
		multiplier = 2
	case roll >= 60:
		won = true
		multiplier = 1
	default:
		won = false
	}

	if won {
		payout := wager * multiplier
		record.Balance += payout
		record.Earned += payout
		record.GambleWins++
	} else {
		record.Balance -= wager
		record.Spent += wager
	}
	record.GambleTotal++
	record.LastGamble = now
	record.UpdatedAt = now

	return GambleResult{
		Record:        record.copy(),
		Roll:          roll,
		Wager:         wager,
		Won:           won,
		Multiplier:    multiplier,
		NextAvailable: now.Add(GambleCooldown),
	}, nil
}

func (s *Store) Sell(userID snowflake.ID, name string, itemID string) (SellResult, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	record := s.ensureRecord(userID, name)
	if record.Inventory == nil {
		return SellResult{}, ErrNotOwned
	}

	entry, ok := record.Inventory[itemID]
	if !ok {
		return SellResult{}, ErrNotOwned
	}

	// find the shop item price
	var shopPrice int
	for _, shopItem := range ShopItems {
		if shopItem.ID == itemID {
			shopPrice = shopItem.Price
			break
		}
	}

	refund := shopPrice * 60 / 100
	if refund < 1 {
		refund = 1
	}

	entry.Quantity--
	if entry.Quantity <= 0 {
		delete(record.Inventory, itemID)
	}

	record.Balance += refund
	record.Earned += refund
	record.UpdatedAt = time.Now()

	return SellResult{
		Record: record.copy(),
		Item:   *entry,
		Reward: refund,
	}, nil
}

func (s *Store) UseItem(userID snowflake.ID, name string, itemID string, now time.Time) (UseItemResult, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	record := s.ensureRecord(userID, name)
	if record.Inventory == nil {
		return UseItemResult{}, ErrNotOwned
	}

	entry, ok := record.Inventory[itemID]
	if !ok {
		return UseItemResult{}, ErrNotOwned
	}

	entry.Quantity--
	if entry.Quantity <= 0 {
		delete(record.Inventory, itemID)
	}
	record.UpdatedAt = now

	switch itemID {
	case "shield_token":
		return UseItemResult{
			Record:      record.copy(),
			ItemID:      itemID,
			ItemName:    "Shield Token",
			Used:        true,
			Description: "You activate the Shield Token. You are now protected from the next robbery attempt!",
		}, nil
	case "health_potion":
		refund := 0
		for _, entry := range record.Inventory {
			refund += entry.Quantity
		}
		heal := 25 + refund*5
		record.Balance += heal
		record.Earned += heal
		return UseItemResult{
			Record:      record.copy(),
			ItemID:      itemID,
			ItemName:    "Health Potion",
			Used:        true,
			Description: fmt.Sprintf("You drink the potion and recover **%d Pulses**.", heal),
		}, nil
	default:
		// restore the item since it can't be used
		if entry.Quantity <= 0 {
			record.Inventory[itemID] = entry
			entry.Quantity++
		} else {
			entry.Quantity++
		}
		return UseItemResult{}, ErrCannotUse
	}
}

func (s *Store) ApplyInterest(userID snowflake.ID, name string, now time.Time) int {
	s.mu.Lock()
	defer s.mu.Unlock()

	record := s.ensureRecord(userID, name)
	if !record.LastInterest.IsZero() && now.Sub(record.LastInterest) < InterestInterval {
		return 0
	}

	interest := int(float64(record.Balance) * InterestRate)
	if interest < 1 {
		interest = 0
	}

	if interest > 0 {
		record.Balance += interest
		record.Earned += interest
	}
	record.LastInterest = now
	record.UpdatedAt = now
	return interest
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
		UserID:      userID,
		Name:        name,
		Balance:     StartingBalance,
		LastInterest: now,
		Inventory:   make(map[string]*InventoryEntry),
		CreatedAt:   now,
		UpdatedAt:   now,
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
