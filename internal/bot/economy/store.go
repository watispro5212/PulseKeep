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
	FlipWins    int
	FlipLosses  int
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
