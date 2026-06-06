package economy

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"log"
	"math/rand"
	"sort"
	"sync"
	"time"

	"github.com/disgoorg/snowflake/v2"
)

const (
	StartingBalance  = 250
	DailyCooldown    = 24 * time.Hour
	WorkCooldown     = 45 * time.Minute
	RobCooldown      = 4 * time.Hour
	FishingCooldown  = 45 * time.Second
	MiningCooldown   = 45 * time.Second
	GambleCooldown   = 10 * time.Second
	InterestRate     = 0.001
	InterestInterval = 6 * time.Hour
	WeeklyCooldown   = 7 * 24 * time.Hour
	DailyBaseReward  = 500
	DailyStreakStep  = 75
	DailyStreakCap   = 750
	WorkMinReward    = 250
	WorkMaxReward    = 650
	WeeklyBaseReward = 2500
	WeeklyStreakStep = 500
	WeeklyStreakCap  = 5000
	MaxWager         = 50000
	MaxTransfer      = 250000
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
	db      *sql.DB
	done    chan struct{}
}

type Record struct {
	UserID            snowflake.ID
	Name              string
	Balance           int
	Earned            int
	Spent             int
	DailyStreak       int
	LastDaily         time.Time
	LastWork          time.Time
	LastRob           time.Time
	LastFish          time.Time
	LastMine          time.Time
	LastGamble        time.Time
	LastInterest      time.Time
	FlipWins          int
	FlipLosses        int
	SlotWins          int
	SlotLosses        int
	RobWins           int
	RobLosses         int
	FishCaught        int
	FishTotal         int
	MineMined         int
	MineTotal         int
	GambleWins        int
	GambleTotal       int
	BlackjackWins     int
	BlackjackLosses   int
	LotteryWins       int
	WeeklyStreak      int
	LastWeekly        time.Time
	XpBoostExpires    time.Time
	TreasureMapActive bool
	GambleBoostActive bool
	Inventory         map[string]*InventoryEntry
	CreatedAt         time.Time
	UpdatedAt         time.Time
}

type DailyResult struct {
	Record        Record
	Reward        int
	StreakBonus   int
	Streak        int
	NextAvailable time.Time
	OnCooldown    bool
}

type WorkResult struct {
	Record        Record
	Reward        int
	BaseReward    int
	Bonus         int
	Job           string
	Boosted       bool
	RareEvent     string
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
	Payout int
	Won    bool
	Boosted bool
}

type RobResult struct {
	Record        Record
	Target        Record
	Stolen        int
	Fine          int
	Success       bool
	Shielded      bool
	NextAvailable time.Time
	OnCooldown    bool
}

type ShopItem struct {
	ID          string
	Name        string
	Description string
	Price       int
	Usable      bool
	Sellable    bool
	OneTimeUse  bool
}

var ShopItems = []ShopItem{
	{ID: "lucky_pickaxe", Name: "Lucky Pickaxe", Description: "+15% coinflip win chance", Price: 5000, Sellable: true},
	{ID: "xp_boost", Name: "XP Boost", Description: "Doubles work earnings for 30 minutes", Price: 10000, Usable: true, Sellable: false},
	{ID: "golden_watch", Name: "Golden Watch", Description: "Reduces daily cooldown by 4 hours", Price: 8000, Usable: true, Sellable: true},
	{ID: "shield_token", Name: "Shield Token", Description: "Protects you from one robbery (consumed on use)", Price: 6000, Usable: true, Sellable: true},
	{ID: "lucky_clover", Name: "Lucky Clover", Description: "Boosts your next gamble win chance by 10%", Price: 5000, Usable: true, Sellable: true},
	{ID: "fishing_rod", Name: "Fishing Rod", Description: "Unlocks the /fish command", Price: 1500, Sellable: true},
	{ID: "iron_pickaxe", Name: "Iron Pickaxe", Description: "Unlocks the /mine command", Price: 2000, Sellable: true},
	{ID: "lottery_ticket_low", Name: "Low Lottery Ticket", Description: "Entry for the standard server lottery draw", Price: 500, Usable: true, Sellable: true},
	{ID: "lottery_ticket_high", Name: "High Lottery Ticket", Description: "Entry for the premium high-roller lottery draw", Price: 5000, Usable: true, Sellable: true},
	{ID: "health_potion", Name: "Health Potion", Description: "Restores 625 Pulses when used (25% refund of purchase price)", Price: 2500, Usable: true, Sellable: true},
	{ID: "treasure_map", Name: "Treasure Map", Description: "Reveals a hidden treasure worth 2500–7500 Pulses", Price: 3000, Usable: true, Sellable: false, OneTimeUse: true},
}

type BuyResult struct {
	Record Record
	Item   ShopItem
}

type SlotsResult struct {
	Record     Record
	Wager      int
	Payout     int
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
	{Name: "Pike", Emoji: "🐟", Weight: "6 lbs", MinPay: 40, MaxPay: 90, Rarity: "Uncommon"},
	{Name: "Salmon", Emoji: "🐟", Weight: "5 lbs", MinPay: 50, MaxPay: 120, Rarity: "Rare"},
	{Name: "Tuna", Emoji: "🐟", Weight: "20 lbs", MinPay: 80, MaxPay: 200, Rarity: "Rare"},
	{Name: "Catfish", Emoji: "🐱", Weight: "8 lbs", MinPay: 60, MaxPay: 150, Rarity: "Rare"},
	{Name: "Mahi-Mahi", Emoji: "🐬", Weight: "15 lbs", MinPay: 120, MaxPay: 250, Rarity: "Epic"},
	{Name: "Swordfish", Emoji: "🗡️", Weight: "50 lbs", MinPay: 150, MaxPay: 350, Rarity: "Epic"},
	{Name: "Great White Shark", Emoji: "🦈", Weight: "1500 lbs", MinPay: 400, MaxPay: 1000, Rarity: "Legendary"},
	{Name: "Golden Koi", Emoji: "✨", Weight: "3 lbs", MinPay: 300, MaxPay: 800, Rarity: "Legendary"},
	{Name: "Kraken Tentacle", Emoji: "🐙", Weight: "300 lbs", MinPay: 800, MaxPay: 2500, Rarity: "Mythic"},
	{Name: "Ancient Coelacanth", Emoji: "🦕", Weight: "80 lbs", MinPay: 500, MaxPay: 1500, Rarity: "Mythic"},
	{Name: "Boot", Emoji: "👢", Weight: "1 lb", MinPay: 0, MaxPay: 5, Rarity: "Junk"},
	{Name: "Old Tire", Emoji: "🔘", Weight: "15 lbs", MinPay: 0, MaxPay: 10, Rarity: "Junk"},
	{Name: "Seaweed", Emoji: "🌿", Weight: "0.2 lbs", MinPay: 2, MaxPay: 8, Rarity: "Junk"},
	{Name: "Rusty Can", Emoji: "🥫", Weight: "0.5 lbs", MinPay: 1, MaxPay: 3, Rarity: "Junk"},
}

type FishResult struct {
	Record        Record
	Fish          FishType
	Reward        int
	Bonus         int
	Boosted       bool
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
	{Name: "Amethyst", Emoji: "🔮", MinPay: 60, MaxPay: 130, Rarity: "Uncommon"},
	{Name: "Platinum", Emoji: "💎", MinPay: 80, MaxPay: 180, Rarity: "Rare"},
	{Name: "Ruby", Emoji: "🔴", MinPay: 120, MaxPay: 250, Rarity: "Rare"},
	{Name: "Sapphire", Emoji: "🔵", MinPay: 150, MaxPay: 300, Rarity: "Epic"},
	{Name: "Emerald", Emoji: "🟢", MinPay: 200, MaxPay: 400, Rarity: "Epic"},
	{Name: "Obsidian", Emoji: "⬛", MinPay: 250, MaxPay: 450, Rarity: "Epic"},
	{Name: "Diamond", Emoji: "💠", MinPay: 300, MaxPay: 600, Rarity: "Legendary"},
	{Name: "Dragonstone", Emoji: "🐉", MinPay: 450, MaxPay: 900, Rarity: "Legendary"},
	{Name: "Star Fragment", Emoji: "⭐", MinPay: 800, MaxPay: 2000, Rarity: "Mythic"},
	{Name: "Ancient Relic", Emoji: "🏺", MinPay: 500, MaxPay: 1200, Rarity: "Mythic"},
	{Name: "Stone", Emoji: "🪨", MinPay: 1, MaxPay: 5, Rarity: "Junk"},
	{Name: "Gravel", Emoji: "🪨", MinPay: 1, MaxPay: 3, Rarity: "Junk"},
	{Name: "Fossilized Bone", Emoji: "🦴", MinPay: 5, MaxPay: 15, Rarity: "Junk"},
}

type MineResult struct {
	Record        Record
	Ore           OreType
	Reward        int
	Bonus         int
	Boosted       bool
	NextAvailable time.Time
	OnCooldown    bool
}

type GambleResult struct {
	Record        Record
	Roll          int
	Wager         int
	Payout        int
	Won           bool
	Push          bool
	Multiplier    int
	Boosted       bool
	NextAvailable time.Time
	OnCooldown    bool
}

type SellResult struct {
	Record Record
	Item   InventoryEntry
	Reward int
}

type UseItemResult struct {
	Record      Record
	ItemID      string
	ItemName    string
	Used        bool
	Description string
}

type WeeklyResult struct {
	Record        Record
	Reward        int
	StreakBonus   int
	Streak        int
	NextAvailable time.Time
	OnCooldown    bool
}

type GiftResult struct {
	Sender   Record
	Receiver Record
	Item     InventoryEntry
}

func NewStore(database *sql.DB) *Store {
	s := &Store{
		records: make(map[snowflake.ID]*Record),
		db:      database,
	}
	s.load()
	return s
}

func (s *Store) load() {
	if s.db == nil {
		return
	}

	ctx := context.Background()
	rows, err := s.db.QueryContext(ctx, `
		SELECT user_id, name, balance, total_earned, total_spent,
		       daily_streak, last_daily, last_work, last_rob,
		       last_fish, last_mine, last_gamble, last_interest,
		       flip_wins, flip_losses, slot_wins, slot_losses,
		       rob_wins, rob_losses, fish_caught, fish_total,
		       mine_mined, mine_total, gamble_wins, gamble_total,
		       blackjack_wins, blackjack_losses,
		       lottery_wins, weekly_streak, last_weekly,
		       xp_boost_expires, treasure_map_active, gamble_boost_active,
		       created_at, updated_at
		FROM user_economy`)
	if err != nil {
		log.Printf("Failed to load economy records: %v", err)
		return
	}
	defer rows.Close()

	for rows.Next() {
		var (
			userIDStr, name                                          string
			balance, earned, spent, dailyStreak, weeklyStreak        int
			flipWins, flipLosses, slotWins, slotLosses               int
			robWins, robLosses, fishCaught, fishTotal                int
			mineMined, mineTotal, gambleWins, gambleTotal            int
			blackjackWins, blackjackLosses, lotteryWins              int
			lastDaily, lastWork, lastRob                             *time.Time
			lastFish, lastMine, lastGamble, lastInterest, lastWeekly *time.Time
			xpBoostExpires                                           *time.Time
			treasureMapActive                                        bool
			gambleBoostActive                                        bool
			createdAt, updatedAt                                     time.Time
		)
		err := rows.Scan(&userIDStr, &name, &balance, &earned, &spent,
			&dailyStreak,
			&lastDaily, &lastWork, &lastRob,
			&lastFish, &lastMine, &lastGamble, &lastInterest,
			&flipWins, &flipLosses, &slotWins, &slotLosses,
			&robWins, &robLosses, &fishCaught, &fishTotal,
			&mineMined, &mineTotal, &gambleWins, &gambleTotal,
			&blackjackWins, &blackjackLosses,
			&lotteryWins, &weeklyStreak, &lastWeekly,
			&xpBoostExpires, &treasureMapActive, &gambleBoostActive, &createdAt, &updatedAt)
		if err != nil {
			log.Printf("Failed to scan economy record: %v", err)
			continue
		}

		userID, err := snowflake.Parse(userIDStr)
		if err != nil {
			log.Printf("Failed to parse user ID %s: %v", userIDStr, err)
			continue
		}

		record := &Record{
			UserID:            userID,
			Name:              name,
			Balance:           balance,
			Earned:            earned,
			Spent:             spent,
			DailyStreak:       dailyStreak,
			WeeklyStreak:      weeklyStreak,
			FlipWins:          flipWins,
			FlipLosses:        flipLosses,
			SlotWins:          slotWins,
			SlotLosses:        slotLosses,
			RobWins:           robWins,
			RobLosses:         robLosses,
			FishCaught:        fishCaught,
			FishTotal:         fishTotal,
			MineMined:         mineMined,
			MineTotal:         mineTotal,
			GambleWins:        gambleWins,
			GambleTotal:       gambleTotal,
			BlackjackWins:     blackjackWins,
			BlackjackLosses:   blackjackLosses,
			LotteryWins:       lotteryWins,
			TreasureMapActive: treasureMapActive,
			GambleBoostActive: gambleBoostActive,
			Inventory:         make(map[string]*InventoryEntry),
			CreatedAt:         createdAt,
			UpdatedAt:         updatedAt,
		}
		if lastDaily != nil {
			record.LastDaily = *lastDaily
		}
		if lastWork != nil {
			record.LastWork = *lastWork
		}
		if lastRob != nil {
			record.LastRob = *lastRob
		}
		if lastFish != nil {
			record.LastFish = *lastFish
		}
		if lastMine != nil {
			record.LastMine = *lastMine
		}
		if lastGamble != nil {
			record.LastGamble = *lastGamble
		}
		if lastInterest != nil {
			record.LastInterest = *lastInterest
		}
		if lastWeekly != nil {
			record.LastWeekly = *lastWeekly
		}
		if xpBoostExpires != nil {
			record.XpBoostExpires = *xpBoostExpires
		}
		s.records[userID] = record
	}
	if err := rows.Err(); err != nil {
		log.Printf("Error iterating economy records: %v", err)
	}

	invRows, err := s.db.QueryContext(ctx, `
		SELECT user_id, item_id, item_name, quantity FROM user_inventory WHERE quantity > 0
		ORDER BY user_id, item_id`)
	if err != nil {
		log.Printf("Failed to load inventory: %v", err)
		return
	}
	defer invRows.Close()

	for invRows.Next() {
		var userIDStr, itemID, itemName string
		var quantity int
		if err := invRows.Scan(&userIDStr, &itemID, &itemName, &quantity); err != nil {
			log.Printf("Failed to scan inventory item: %v", err)
			continue
		}
		userID, err := snowflake.Parse(userIDStr)
		if err != nil {
			log.Printf("Failed to parse user ID %s in inventory: %v", userIDStr, err)
			continue
		}
		if record, ok := s.records[userID]; ok {
			record.Inventory[itemID] = &InventoryEntry{
				ItemID:   itemID,
				ItemName: itemName,
				Quantity: quantity,
			}
		}
	}
	if err := invRows.Err(); err != nil {
		log.Printf("Error iterating inventory records: %v", err)
	}
	log.Printf("Loaded %d economy records from database", len(s.records))
}

func (s *Store) save(record *Record) {
	if s.db == nil {
		return
	}

	ctx := context.Background()
	_, err := s.db.ExecContext(ctx, `
		INSERT INTO user_economy (
			user_id, name, balance, total_earned, total_spent,
			daily_streak, last_daily, last_work, last_rob,
			last_fish, last_mine, last_gamble, last_interest,
			flip_wins, flip_losses, slot_wins, slot_losses,
			rob_wins, rob_losses, fish_caught, fish_total,
			mine_mined, mine_total, gamble_wins, gamble_total,
			blackjack_wins, blackjack_losses,
			lottery_wins, weekly_streak, last_weekly,
			xp_boost_expires, treasure_map_active, gamble_boost_active,
			created_at, updated_at
		) VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11,$12,$13,
		          $14,$15,$16,$17,$18,$19,$20,$21,$22,$23,$24,$25,$26,$27,$28,$29,
		          $30,$31,$32,$33,$34,$35)
		ON CONFLICT (user_id) DO UPDATE SET
			name = EXCLUDED.name,
			balance = EXCLUDED.balance,
			total_earned = EXCLUDED.total_earned,
			total_spent = EXCLUDED.total_spent,
			daily_streak = EXCLUDED.daily_streak,
			last_daily = EXCLUDED.last_daily,
			last_work = EXCLUDED.last_work,
			last_rob = EXCLUDED.last_rob,
			last_fish = EXCLUDED.last_fish,
			last_mine = EXCLUDED.last_mine,
			last_gamble = EXCLUDED.last_gamble,
			last_interest = EXCLUDED.last_interest,
			flip_wins = EXCLUDED.flip_wins,
			flip_losses = EXCLUDED.flip_losses,
			slot_wins = EXCLUDED.slot_wins,
			slot_losses = EXCLUDED.slot_losses,
			rob_wins = EXCLUDED.rob_wins,
			rob_losses = EXCLUDED.rob_losses,
			fish_caught = EXCLUDED.fish_caught,
			fish_total = EXCLUDED.fish_total,
			mine_mined = EXCLUDED.mine_mined,
			mine_total = EXCLUDED.mine_total,
			gamble_wins = EXCLUDED.gamble_wins,
			gamble_total = EXCLUDED.gamble_total,
			blackjack_wins = EXCLUDED.blackjack_wins,
			blackjack_losses = EXCLUDED.blackjack_losses,
			lottery_wins = EXCLUDED.lottery_wins,
			weekly_streak = EXCLUDED.weekly_streak,
			last_weekly = EXCLUDED.last_weekly,
			xp_boost_expires = EXCLUDED.xp_boost_expires,
			treasure_map_active = EXCLUDED.treasure_map_active,
			gamble_boost_active = EXCLUDED.gamble_boost_active,
			updated_at = EXCLUDED.updated_at
	`, record.UserID.String(), record.Name, record.Balance, record.Earned, record.Spent,
		record.DailyStreak,
		nullTime(record.LastDaily), nullTime(record.LastWork), nullTime(record.LastRob),
		nullTime(record.LastFish), nullTime(record.LastMine), nullTime(record.LastGamble),
		nullTime(record.LastInterest),
		record.FlipWins, record.FlipLosses, record.SlotWins, record.SlotLosses,
		record.RobWins, record.RobLosses, record.FishCaught, record.FishTotal,
		record.MineMined, record.MineTotal, record.GambleWins, record.GambleTotal,
		record.BlackjackWins, record.BlackjackLosses,
		record.LotteryWins, record.WeeklyStreak, nullTime(record.LastWeekly),
		nullTime(record.XpBoostExpires), record.TreasureMapActive, record.GambleBoostActive,
		record.CreatedAt, record.UpdatedAt)
	if err != nil {
		log.Printf("Failed to save economy record for user %s: %v", record.UserID, err)
	}

	for _, entry := range record.Inventory {
		_, err = s.db.ExecContext(ctx, `
			INSERT INTO user_inventory (user_id, item_id, item_name, quantity)
			VALUES ($1, $2, $3, $4)
			ON CONFLICT (user_id, item_id) DO UPDATE SET
				item_name = EXCLUDED.item_name,
				quantity = EXCLUDED.quantity
		`, record.UserID.String(), entry.ItemID, entry.ItemName, entry.Quantity)
		if err != nil {
			log.Printf("Failed to save inventory item %s for user %s: %v", entry.ItemID, record.UserID, err)
		}
	}
}

func (s *Store) FlushAll() {
	s.mu.Lock()
	defer s.mu.Unlock()
	for _, record := range s.records {
		s.save(record)
	}
}

func nullTime(t time.Time) *time.Time {
	if t.IsZero() {
		return nil
	}
	return &t
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
	streakBonus := min(record.DailyStreak*DailyStreakStep, DailyStreakCap)
	reward := DailyBaseReward + streakBonus
	record.Balance += reward
	record.Earned += reward
	record.LastDaily = now
	record.UpdatedAt = now
	s.save(record)

	return DailyResult{
		Record:        record.copy(),
		Reward:        reward,
		StreakBonus:   streakBonus,
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
		"coded a new PulseKeep feature",
		"debugged a tricky race condition",
		"optimized slow database queries",
		"wrote unit tests for the economy",
		"refactored legacy command handlers",
		"updated the slash command registry",
		"monitored server performance metrics",
		"responded to community feedback",
		"balanced gambling win rates",
		"designed new fish and ore types",
		"rewrote the ticket dispatch system",
		"audited permission checks across guilds",
		"improved embed formatting across the board",
		"reviewed and merged open pull requests",
		"ran a security audit on command inputs",
		"drafted the PulseKeep changelog update",
	}
	baseReward := WorkMinReward + rand.Intn(WorkMaxReward-WorkMinReward+1)
	reward := baseReward
	job := jobs[rand.Intn(len(jobs))]
	var rareEvent string
	bonus := 0
	if rand.Intn(100) < 12 {
		rareEvents := []string{
			"Quality streak bonus",
			"Emergency support payout",
			"Automation cleanup bounty",
			"Community thank-you tip",
			"Late-shift performance bonus",
		}
		rareEvent = rareEvents[rand.Intn(len(rareEvents))]
		bonus = 125 + rand.Intn(376)
		reward += bonus
	}

	// XP Boost — double work earnings
	boosted := false
	if !record.XpBoostExpires.IsZero() && now.Before(record.XpBoostExpires) {
		reward *= 2
		boosted = true
	} else {
		record.XpBoostExpires = time.Time{} // clear expired boost
	}

	record.Balance += reward
	record.Earned += reward
	record.LastWork = now
	record.UpdatedAt = now
	s.save(record)

	return WorkResult{
		Record:        record.copy(),
		Reward:        reward,
		BaseReward:    baseReward,
		Bonus:         bonus,
		Job:           job,
		Boosted:       boosted,
		RareEvent:     rareEvent,
		NextAvailable: now.Add(WorkCooldown),
	}
}

func (s *Store) Pay(senderID snowflake.ID, senderName string, receiverID snowflake.ID, receiverName string, amount int, now time.Time) (TransferResult, error) {
	if amount <= 0 || amount > MaxTransfer {
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
	s.save(sender)
	s.save(receiver)

	return TransferResult{
		Sender:    sender.copy(),
		Receiver:  receiver.copy(),
		Amount:    amount,
		Recipient: receiver.Name,
	}, nil
}

func (s *Store) Coinflip(userID snowflake.ID, name string, side string, wager int, now time.Time) (FlipResult, error) {
	if wager <= 0 || wager > MaxWager {
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

	// check for Lucky Pickaxe — +15% win chance
	hasPickaxe := false
	if record.Inventory != nil {
		if _, ok := record.Inventory["lucky_pickaxe"]; ok {
			hasPickaxe = true
		}
	}

	result := "heads"
	if rand.Intn(2) == 1 {
		result = "tails"
	}
	won := side == result

	// Lucky Pickaxe: on loss, reroll once (effective ~57.5% win rate)
	if !won && hasPickaxe {
		result = "heads"
		if rand.Intn(2) == 1 {
			result = "tails"
		}
		won = side == result
	}

	// consume lucky pickaxe if used
	if hasPickaxe {
		entry := record.Inventory["lucky_pickaxe"]
		entry.Quantity--
		if entry.Quantity <= 0 {
			delete(record.Inventory, "lucky_pickaxe")
		}
	}
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
	s.save(record)

	payout := 0
	if won {
		payout = wager
	}
	return FlipResult{
		Record: record.copy(),
		Side:   side,
		Result: result,
		Wager:  wager,
		Payout: payout,
		Won:    won,
		Boosted: hasPickaxe,
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
		record.RobLosses++
		record.LastRob = now
		record.UpdatedAt = now
		s.save(record)
		s.save(target)
		return RobResult{
			Record:        record.copy(),
			Target:        target.copy(),
			Fine:          fine,
			Shielded:      true,
			NextAvailable: now.Add(RobCooldown),
		}, nil
	}

	success := rand.Intn(100) < 45
	var stolen int
	var fine int

	if success {
		stolen = target.Balance * (10 + rand.Intn(21)) / 100
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
		fine = record.Balance * (8 + rand.Intn(18)) / 100
		if fine < 10 {
			fine = 10
		}
		if fine > record.Balance {
			fine = record.Balance
		}
		record.Balance -= fine
		record.Spent += fine
		record.RobLosses++
	}

	record.LastRob = now
	record.UpdatedAt = now
	s.save(record)
	s.save(target)

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
		return BuyResult{}, ErrItemNotFound
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
	s.save(record)

	return BuyResult{
		Record: record.copy(),
		Item:   *item,
	}, nil
}

func (s *Store) Inventory(userID snowflake.ID, name string) []InventoryEntry {
	s.mu.Lock()
	defer s.mu.Unlock()

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
	if wager <= 0 || wager > MaxWager {
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
	} else if rand.Intn(100) < 10 {
		// Small chance to nudge into a pair
		reels[2] = reels[1]
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
	s.save(record)

	return SlotsResult{
		Record:     record.copy(),
		Wager:      wager,
		Payout:     wager * multiplier,
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
		s.save(record)
		return FishResult{
			Record:        record.copy(),
			Fish:          fish,
			Reward:        0,
			NextAvailable: now.Add(FishingCooldown),
		}
	}

	idx := rand.Intn(len(FishTable))
	// 15% chance to reroll if junk — slightly better rare rates
	if idx >= 10 && rand.Intn(100) < 15 {
		idx = rand.Intn(len(FishTable) - 3)
	}
	fish := FishTable[idx]
	reward := fish.MinPay + rand.Intn(fish.MaxPay-fish.MinPay+1)
	bonus := 0
	if reward > 0 {
		switch fish.Rarity {
		case "Legendary":
			bonus = 250 + rand.Intn(501)
		case "Mythic":
			bonus = 750 + rand.Intn(1001)
		}
		reward += bonus
	}
	boosted := false

	if reward > 0 {
		// Treasure Map — double payout, then consume
		if record.TreasureMapActive {
			reward *= 2
			boosted = true
			record.TreasureMapActive = false
		}
		record.Balance += reward
		record.Earned += reward
		record.FishCaught++
	}
	record.FishTotal++
	record.LastFish = now
	record.UpdatedAt = now
	s.save(record)

	return FishResult{
		Record:        record.copy(),
		Fish:          fish,
		Reward:        reward,
		Bonus:         bonus,
		Boosted:       boosted,
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
		s.save(record)
		return MineResult{
			Record:        record.copy(),
			Ore:           ore,
			Reward:        0,
			NextAvailable: now.Add(MiningCooldown),
		}
	}

	idx := rand.Intn(len(OreTable))
	// 15% chance to reroll if junk — slightly better rare rates
	if idx >= 10 && rand.Intn(100) < 15 {
		idx = rand.Intn(len(OreTable) - 3)
	}
	ore := OreTable[idx]
	reward := ore.MinPay + rand.Intn(ore.MaxPay-ore.MinPay+1)
	bonus := 0
	if reward > 0 {
		switch ore.Rarity {
		case "Legendary":
			bonus = 250 + rand.Intn(501)
		case "Mythic":
			bonus = 750 + rand.Intn(1001)
		}
		reward += bonus
	}
	boosted := false

	if reward > 0 {
		// Treasure Map — double payout, then consume
		if record.TreasureMapActive {
			reward *= 2
			boosted = true
			record.TreasureMapActive = false
		}
		record.Balance += reward
		record.Earned += reward
		record.MineMined++
	}
	record.MineTotal++
	record.LastMine = now
	record.UpdatedAt = now
	s.save(record)

	return MineResult{
		Record:        record.copy(),
		Ore:           ore,
		Reward:        reward,
		Bonus:         bonus,
		Boosted:       boosted,
		NextAvailable: now.Add(MiningCooldown),
	}
}

func (s *Store) Gamble(userID snowflake.ID, name string, wager int, now time.Time) (GambleResult, error) {
	if wager <= 0 || wager > MaxWager {
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
	var push bool
	multiplier := 0

	// Lucky Clover — boost win chance by 10%
	hasCloverBoost := record.GambleBoostActive
	if hasCloverBoost {
		record.GambleBoostActive = false
	}

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
	case roll >= 46:
		won = true
		multiplier = 1
	case roll >= 36:
		push = true
		multiplier = 0
	default:
		// lose on ≤35 — Lucky Clover gives a second chance
		if hasCloverBoost {
			roll = 1 + rand.Intn(100)
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
			case roll >= 46:
				won = true
				multiplier = 1
			case roll >= 36:
				push = true
				multiplier = 0
			}
		}
	}

	if won {
		payout := wager * multiplier
		record.Balance += payout
		record.Earned += payout
		record.GambleWins++
	} else if push {
		// wager returned, nothing changes
	} else {
		record.Balance -= wager
		record.Spent += wager
	}
	record.GambleTotal++
	record.LastGamble = now
	record.UpdatedAt = now
	s.save(record)

	return GambleResult{
		Record:        record.copy(),
		Roll:          roll,
		Wager:         wager,
		Payout:        wager * multiplier,
		Won:           won,
		Push:          push,
		Multiplier:    multiplier,
		Boosted:       hasCloverBoost,
		NextAvailable: now.Add(GambleCooldown),
	}, nil
}

func (s *Store) Sell(userID snowflake.ID, name string, itemID string, now time.Time) (SellResult, error) {
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
	var sellable bool
	for _, shopItem := range ShopItems {
		if shopItem.ID == itemID {
			shopPrice = shopItem.Price
			sellable = shopItem.Sellable
			break
		}
	}

	if !sellable {
		return SellResult{}, ErrCannotUse
	}

	refund := shopPrice * 60 / 100
	if refund < 1 {
		refund = 1
	}

	soldEntry := *entry
	entry.Quantity--
	if entry.Quantity <= 0 {
		delete(record.Inventory, itemID)
	}

	record.Balance += refund
	record.UpdatedAt = now
	s.save(record)

	return SellResult{
		Record: record.copy(),
		Item:   soldEntry,
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
	record.UpdatedAt = now

	switch itemID {
	case "shield_token":
		entry.Quantity--
		if entry.Quantity <= 0 {
			delete(record.Inventory, itemID)
		}
		s.save(record)
		return UseItemResult{
			Record:      record.copy(),
			ItemID:      itemID,
			ItemName:    "Shield Token",
			Used:        true,
			Description: "You activate the Shield Token. You are now protected from the next robbery attempt!",
		}, nil
	case "health_potion":
		entry.Quantity--
		if entry.Quantity <= 0 {
			delete(record.Inventory, itemID)
		}
		heal := 625
		record.Balance += heal
		record.Earned += heal
		s.save(record)
		return UseItemResult{
			Record:      record.copy(),
			ItemID:      itemID,
			ItemName:    "Health Potion",
			Used:        true,
			Description: fmt.Sprintf("You drink the potion and recover **%d Pulses** (25%% refund).", heal),
		}, nil
	case "xp_boost":
		entry.Quantity--
		if entry.Quantity <= 0 {
			delete(record.Inventory, itemID)
		}
		record.XpBoostExpires = now.Add(30 * time.Minute)
		s.save(record)
		return UseItemResult{
			Record:      record.copy(),
			ItemID:      itemID,
			ItemName:    "XP Boost",
			Used:        true,
			Description: "You activate the XP Boost! Your work earnings are doubled for the next 30 minutes.",
		}, nil
	case "golden_watch":
		if record.LastDaily.IsZero() {
			return UseItemResult{}, ErrCannotUse
		}
		entry.Quantity--
		if entry.Quantity <= 0 {
			delete(record.Inventory, itemID)
		}
		adjusted := record.LastDaily.Add(-4 * time.Hour)
		if adjusted.After(time.Time{}) {
			record.LastDaily = adjusted
		}
		s.save(record)
		return UseItemResult{
			Record:      record.copy(),
			ItemID:      itemID,
			ItemName:    "Golden Watch",
			Used:        true,
			Description: "You wind the Golden Watch! Your `/daily` cooldown has been reduced by 4 hours.",
		}, nil
	case "treasure_map":
		entry.Quantity--
		if entry.Quantity <= 0 {
			delete(record.Inventory, itemID)
		}
		treasure := 2500 + rand.Intn(5001)
		record.Balance += treasure
		record.Earned += treasure
		s.save(record)
		return UseItemResult{
			Record:      record.copy(),
			ItemID:      itemID,
			ItemName:    "Treasure Map",
			Used:        true,
			Description: fmt.Sprintf("You follow the Treasure Map and uncover **%d Pulses** buried underground!", treasure),
		}, nil
	case "lucky_clover":
		entry.Quantity--
		if entry.Quantity <= 0 {
			delete(record.Inventory, itemID)
		}
		record.GambleBoostActive = true
		s.save(record)
		return UseItemResult{
			Record:      record.copy(),
			ItemID:      itemID,
			ItemName:    "Lucky Clover",
			Used:        true,
			Description: "You clutch the Lucky Clover! Your next `/gamble` loss will be rerolled.",
		}, nil
	case "lottery_ticket_low", "lottery_ticket_high":
		entry.Quantity--
		if entry.Quantity <= 0 {
			delete(record.Inventory, itemID)
		}
		tier := "low"
		tierName := "Low Tier"
		if itemID == "lottery_ticket_high" {
			tier = "high"
			tierName = "High Tier"
		}
		s.LotteryBuyTicket(userID, tier, now)
		s.save(record)
		return UseItemResult{
			Record:      record.copy(),
			ItemID:      itemID,
			ItemName:    fmt.Sprintf("%s Lottery Ticket", tierName),
			Used:        true,
			Description: fmt.Sprintf("You enter the weekly %s lottery draw! You will be automatically entered into the next drawing.", tierName),
		}, nil
	default:
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
	s.save(record)
	return interest
}

func (s *Store) Weekly(userID snowflake.ID, name string, now time.Time) WeeklyResult {
	s.mu.Lock()
	defer s.mu.Unlock()

	record := s.ensureRecord(userID, name)
	if !record.LastWeekly.IsZero() && now.Sub(record.LastWeekly) < WeeklyCooldown {
		return WeeklyResult{
			Record:        record.copy(),
			NextAvailable: record.LastWeekly.Add(WeeklyCooldown),
			OnCooldown:    true,
		}
	}

	if record.LastWeekly.IsZero() || now.Sub(record.LastWeekly) > 14*24*time.Hour {
		record.WeeklyStreak = 0
	}

	record.WeeklyStreak++
	streakBonus := min(record.WeeklyStreak*WeeklyStreakStep, WeeklyStreakCap)
	reward := WeeklyBaseReward + streakBonus
	record.Balance += reward
	record.Earned += reward
	record.LastWeekly = now
	record.UpdatedAt = now
	s.save(record)

	return WeeklyResult{
		Record:        record.copy(),
		Reward:        reward,
		StreakBonus:   streakBonus,
		Streak:        record.WeeklyStreak,
		NextAvailable: now.Add(WeeklyCooldown),
	}
}

type LotteryStatus struct {
	TotalEntriesLow  int
	JackpotLow       int
	LastWinnerIDLow  string
	LastDrawTimeLow  time.Time
	AutoDrawLow      bool
	
	TotalEntriesHigh int
	JackpotHigh      int
	LastWinnerIDHigh string
	LastDrawTimeHigh time.Time
	AutoDrawHigh     bool
	
	WeekStart    time.Time
}

func (s *Store) LotteryStatus(now time.Time) LotteryStatus {
	weekStart := weekStartTime(now)

	var totalEntriesLow, totalEntriesHigh int
	_ = s.db.QueryRow(`SELECT COUNT(*) FROM lottery_entries WHERE week_start = $1 AND claimed = false AND tier = 'low'`, weekStart).Scan(&totalEntriesLow)
	_ = s.db.QueryRow(`SELECT COUNT(*) FROM lottery_entries WHERE week_start = $1 AND claimed = false AND tier = 'high'`, weekStart).Scan(&totalEntriesHigh)

	jackpotLow := totalEntriesLow * 400 // 500 * 0.8
	jackpotHigh := totalEntriesHigh * 4000 // 5000 * 0.8

	var lastWinnerIDLow, lastWinnerIDHigh string
	var lastDrawTimeLow, lastDrawTimeHigh time.Time
	_ = s.db.QueryRow(`SELECT COALESCE(last_winner_id,''), COALESCE(last_draw_time, 'epoch'::timestamptz), COALESCE(high_last_winner_id,''), COALESCE(high_last_draw_time, 'epoch'::timestamptz) FROM lottery_config WHERE guild_id = 'global'`).Scan(&lastWinnerIDLow, &lastDrawTimeLow, &lastWinnerIDHigh, &lastDrawTimeHigh)

	autoDrawLow := true
	autoDrawHigh := true
	_ = s.db.QueryRow(`SELECT auto_draw, high_auto_draw FROM lottery_config WHERE guild_id = 'global'`).Scan(&autoDrawLow, &autoDrawHigh)

	return LotteryStatus{
		TotalEntriesLow:  totalEntriesLow,
		JackpotLow:       jackpotLow,
		LastWinnerIDLow:  lastWinnerIDLow,
		LastDrawTimeLow:  lastDrawTimeLow,
		AutoDrawLow:      autoDrawLow,
		
		TotalEntriesHigh: totalEntriesHigh,
		JackpotHigh:      jackpotHigh,
		LastWinnerIDHigh: lastWinnerIDHigh,
		LastDrawTimeHigh: lastDrawTimeHigh,
		AutoDrawHigh:     autoDrawHigh,
		
		WeekStart:    weekStart,
	}
}

func (s *Store) LotteryBuyTicket(userID snowflake.ID, tier string, now time.Time) {
	weekStart := weekStartTime(now)
	_, _ = s.db.Exec(`INSERT INTO lottery_entries (user_id, week_start, tier) VALUES ($1, $2, $3)`, userID.String(), weekStart, tier)
}

func (s *Store) LotteryClaim(userID snowflake.ID, name string, tier string, now time.Time) (UseItemResult, error) {
	weekStart := weekStartTime(now)

	// Determine column names based on tier
	drawTimeCol := "last_draw_time"
	autoDrawCol := "auto_draw"
	winnerCol := "last_winner_id"
	ticketPrice := 500
	itemName := "Low Lottery Ticket"
	if tier == "high" {
		drawTimeCol = "high_last_draw_time"
		autoDrawCol = "high_auto_draw"
		winnerCol = "high_last_winner_id"
		ticketPrice = 5000
		itemName = "High Lottery Ticket"
	}

	var lastDrawTime time.Time
	_ = s.db.QueryRow(fmt.Sprintf(`SELECT COALESCE(%s, 'epoch'::timestamptz) FROM lottery_config WHERE guild_id = 'global'`, drawTimeCol)).Scan(&lastDrawTime)

	var autoDraw bool
	_ = s.db.QueryRow(fmt.Sprintf(`SELECT COALESCE(%s, true) FROM lottery_config WHERE guild_id = 'global'`, autoDrawCol)).Scan(&autoDraw)

	if lastDrawTime.Before(weekStart) {
		if autoDraw {
			s.lotteryRunDraw(weekStart, tier, now)
		} else {
			return UseItemResult{}, errors.New("the weekly draw has not started yet. wait for the auto-draw or ask an admin to trigger one")
		}
	}

	var entryID int
	err := s.db.QueryRow(`SELECT id FROM lottery_entries WHERE week_start = $1 AND user_id = $2 AND claimed = false AND tier = $3`, weekStart, userID.String(), tier).Scan(&entryID)
	if err != nil {
		if err == sql.ErrNoRows {
			return UseItemResult{}, fmt.Errorf("you did not win the %s lottery this week. better luck next time!", tier)
		}
		return UseItemResult{}, errors.New("could not check lottery status")
	}

	var totalEntries int
	_ = s.db.QueryRow(`SELECT COUNT(*) FROM lottery_entries WHERE week_start = $1 AND claimed = false AND tier = $2`, weekStart, tier).Scan(&totalEntries)
	jackpot := int(float64(totalEntries) * float64(ticketPrice) * 0.8)

	s.mu.Lock()
	record := s.ensureRecord(userID, name)
	record.Balance += jackpot
	record.Earned += jackpot
	record.LotteryWins++
	record.UpdatedAt = now
	s.save(record)
	s.mu.Unlock()

	_, _ = s.db.Exec(`UPDATE lottery_entries SET claimed = true WHERE id = $1`, entryID)
	_, _ = s.db.Exec(fmt.Sprintf(`UPDATE lottery_config SET %s = $1 WHERE guild_id = 'global'`, winnerCol), userID.String())

	return UseItemResult{
		Record:      record.copy(),
		ItemID:      "lottery_ticket_" + tier,
		ItemName:    itemName,
		Used:        true,
		Description: fmt.Sprintf("You won the weekly %s lottery! **%d Pulses** have been added to your balance.", tier, jackpot),
	}, nil
}

func (s *Store) LotteryToggleAutoDraw(tier string, enable bool) error {
	col := "auto_draw"
	if tier == "high" {
		col = "high_auto_draw"
	}
	_, err := s.db.Exec(fmt.Sprintf(`
		INSERT INTO lottery_config (guild_id, %s)
		VALUES ('global', $1)
		ON CONFLICT (guild_id) DO UPDATE SET %s = $1
	`, col, col), enable)
	return err
}

func weekStartTime(now time.Time) time.Time {
	// Calculate start of current week (Monday 00:00 UTC)
	daysSinceMonday := int(now.Weekday()) - 1
	if daysSinceMonday < 0 {
		daysSinceMonday = 6
	}
	start := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, time.UTC)
	start = start.AddDate(0, 0, -daysSinceMonday)
	return start
}

func (s *Store) lotteryRunDraw(weekStart time.Time, tier string, now time.Time) {
	var totalEntries int
	_ = s.db.QueryRow(`SELECT COUNT(*) FROM lottery_entries WHERE week_start = $1 AND claimed = false AND tier = $2`, weekStart, tier).Scan(&totalEntries)
	if totalEntries == 0 {
		return
	}

	var userIDStr string
	_ = s.db.QueryRow(`SELECT user_id FROM lottery_entries WHERE week_start = $1 AND claimed = false AND tier = $3 OFFSET floor(random() * $2) LIMIT 1`, weekStart, totalEntries, tier).Scan(&userIDStr)

	drawTimeCol := "last_draw_time"
	winnerCol := "last_winner_id"
	if tier == "high" {
		drawTimeCol = "high_last_draw_time"
		winnerCol = "high_last_winner_id"
	}

	_, _ = s.db.Exec(fmt.Sprintf(`UPDATE lottery_config SET %s = $1, %s = $2 WHERE guild_id = 'global'`, drawTimeCol, winnerCol), now, userIDStr)
}

// StartLotteryAutoDraw starts a goroutine that checks weekly for auto-draw
func (s *Store) StartLotteryAutoDraw(ctx context.Context) {
	go func() {
		ticker := time.NewTicker(1 * time.Hour)
		defer ticker.Stop()
		for {
			select {
			case <-ticker.C:
				var autoDrawLow, autoDrawHigh bool
				_ = s.db.QueryRow(`SELECT COALESCE(auto_draw, true), COALESCE(high_auto_draw, true) FROM lottery_config WHERE guild_id = 'global'`).Scan(&autoDrawLow, &autoDrawHigh)
				
				now := time.Now()
				weekStart := weekStartTime(now)
				
				var lastDrawTimeLow, lastDrawTimeHigh time.Time
				_ = s.db.QueryRow(`SELECT COALESCE(last_draw_time, 'epoch'::timestamptz), COALESCE(high_last_draw_time, 'epoch'::timestamptz) FROM lottery_config WHERE guild_id = 'global'`).Scan(&lastDrawTimeLow, &lastDrawTimeHigh)
				
				if autoDrawLow && lastDrawTimeLow.Before(weekStart) {
					s.lotteryRunDraw(weekStart, "low", now)
				}
				if autoDrawHigh && lastDrawTimeHigh.Before(weekStart) {
					s.lotteryRunDraw(weekStart, "high", now)
				}
			case <-ctx.Done():
				return
			}
		}
	}()
}

func (s *Store) GiftItem(senderID snowflake.ID, senderName string, receiverID snowflake.ID, receiverName string, itemID string, now time.Time) (GiftResult, error) {
	if senderID == receiverID {
		return GiftResult{}, ErrSelfPayment
	}

	s.mu.Lock()
	defer s.mu.Unlock()

	sender := s.ensureRecord(senderID, senderName)
	if sender.Inventory == nil {
		return GiftResult{}, ErrNotOwned
	}

	entry, ok := sender.Inventory[itemID]
	if !ok {
		return GiftResult{}, ErrNotOwned
	}

	receiver := s.ensureRecord(receiverID, receiverName)

	entry.Quantity--
	if entry.Quantity <= 0 {
		delete(sender.Inventory, itemID)
	}

	if receiver.Inventory == nil {
		receiver.Inventory = make(map[string]*InventoryEntry)
	}
	if existing, ok := receiver.Inventory[itemID]; ok {
		existing.Quantity++
	} else {
		receiver.Inventory[itemID] = &InventoryEntry{
			ItemID:   entry.ItemID,
			ItemName: entry.ItemName,
			Quantity: 1,
		}
	}

	sender.UpdatedAt = now
	receiver.UpdatedAt = now
	s.save(sender)
	s.save(receiver)

	return GiftResult{
		Sender:   sender.copy(),
		Receiver: receiver.copy(),
		Item:     *entry,
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
		UserID:       userID,
		Name:         name,
		Balance:      StartingBalance,
		LastInterest: now,
		Inventory:    make(map[string]*InventoryEntry),
		CreatedAt:    now,
		UpdatedAt:    now,
	}
	s.records[userID] = record
	s.save(record)
	return record
}

func (r Record) NetWorth() int {
	val := r.Balance
	for _, entry := range r.Inventory {
		for _, shopItem := range ShopItems {
			if shopItem.ID == entry.ItemID {
				val += shopItem.Price * entry.Quantity
				break
			}
		}
	}
	return val
}

func (r Record) FlipTotal() int {
	return r.FlipWins + r.FlipLosses
}

func (r Record) copy() Record {
	c := r
	if r.Inventory != nil {
		c.Inventory = make(map[string]*InventoryEntry, len(r.Inventory))
		for k, v := range r.Inventory {
			entry := *v
			c.Inventory[k] = &entry
		}
	}
	return c
}
