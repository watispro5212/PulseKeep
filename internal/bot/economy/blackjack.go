package economy

import (
	"fmt"
	"math/rand"
	"sync"
	"time"

	"github.com/disgoorg/snowflake/v2"
)

type BlackjackDifficulty int

const (
	BlackjackEasy   BlackjackDifficulty = 14
	BlackjackNormal BlackjackDifficulty = 17
	BlackjackHard   BlackjackDifficulty = 19
	BlackjackExpert BlackjackDifficulty = 20

	BlackjackHitCustomID   = "pulsekeep:bj:hit"
	BlackjackStandCustomID = "pulsekeep:bj:stand"
)

const blackjackTimeout = 90 * time.Second

type BlackjackHand struct {
	Cards   []int
	Value   int
	Bust    bool
	CardStr string
}

type PendingBJ struct {
	mu         sync.Mutex
	UserID     snowflake.ID
	Name       string
	Player     []int
	Dealer     []int
	Wager      int
	Difficulty BlackjackDifficulty
	StartedAt  time.Time
}

type BlackjackStartResult struct {
	Record Record
	Player BlackjackHand
	Dealer BlackjackHand
	Wager  int
	Won    bool
	Payout int
	Push   bool
	Natural bool
	GameOver bool
}

type BlackjackTurnResult struct {
	Record  Record
	Player  BlackjackHand
	Dealer  BlackjackHand
	Wager   int
	Won     bool
	Payout  int
	Push    bool
	GameOver bool
}

var (
	bjGamesMu sync.RWMutex
	bjGames   = make(map[snowflake.ID]*PendingBJ)
)

func init() {
	go bjCleanupLoop()
}

func bjCleanupLoop() {
	for {
		time.Sleep(30 * time.Second)
		bjGamesMu.Lock()
		now := time.Now()
		for id, g := range bjGames {
			g.mu.Lock()
			if now.Sub(g.StartedAt) > blackjackTimeout {
				g.mu.Unlock()
				delete(bjGames, id)
			} else {
				g.mu.Unlock()
			}
		}
		bjGamesMu.Unlock()
	}
}

func (s *Store) BlackjackStart(userID snowflake.ID, name string, wager int, difficulty BlackjackDifficulty, now time.Time) (BlackjackStartResult, error) {
	if wager <= 0 {
		return BlackjackStartResult{}, ErrInvalidAmount
	}

	s.mu.Lock()
	record := s.ensureRecord(userID, name)
	if record.Balance < wager {
		s.mu.Unlock()
		return BlackjackStartResult{}, ErrInsufficientFund
	}

	// Check for existing pending game
	bjGamesMu.RLock()
	existing := bjGames[userID]
	bjGamesMu.RUnlock()
	if existing != nil {
		s.mu.Unlock()
		return BlackjackStartResult{}, fmt.Errorf("you already have an active blackjack game")
	}

	record.Balance -= wager
	record.Spent += wager
	record.UpdatedAt = now
	s.save(record)
	s.mu.Unlock()

	pCards := []int{drawCard(), drawCard()}
	dCards := []int{drawCard(), drawCard()}
	pVal := handValue(pCards)
	dVal := handValue(dCards)

	// Check for naturals
	if pVal == 21 && dVal == 21 {
		record.Balance += wager
		s.mu.Lock()
		s.save(record)
		s.mu.Unlock()
		return BlackjackStartResult{
			Record:   record.copy(),
			Player:   BlackjackHand{Cards: pCards, Value: pVal, CardStr: cardsStr(pCards)},
			Dealer:   BlackjackHand{Cards: dCards, Value: dVal, CardStr: cardsStr(dCards)},
			Wager:    wager,
			Push:     true,
			GameOver: true,
		}, nil
	}

	if pVal == 21 {
		payout := wager * 3 / 2
		record.Balance += payout
		record.Earned += payout
		record.BlackjackWins++
		s.mu.Lock()
		s.save(record)
		s.mu.Unlock()
		return BlackjackStartResult{
			Record:   record.copy(),
			Player:   BlackjackHand{Cards: pCards, Value: pVal, CardStr: cardsStr(pCards)},
			Dealer:   BlackjackHand{Cards: dCards, Value: dVal, CardStr: cardsStr(dCards)},
			Wager:    wager,
			Won:      true,
			Payout:   payout,
			Natural:  true,
			GameOver: true,
		}, nil
	}

	if dVal == 21 {
		record.BlackjackLosses++
		s.mu.Lock()
		s.save(record)
		s.mu.Unlock()
		return BlackjackStartResult{
			Record:   record.copy(),
			Player:   BlackjackHand{Cards: pCards, Value: pVal, CardStr: cardsStr(pCards)},
			Dealer:   BlackjackHand{Cards: dCards, Value: dVal, CardStr: cardsStr(dCards)},
			Wager:    wager,
			Payout:   -wager,
			GameOver: true,
		}, nil
	}

	// Store pending game
	pending := &PendingBJ{
		UserID:     userID,
		Name:       name,
		Player:     pCards,
		Dealer:     dCards,
		Wager:      wager,
		Difficulty: difficulty,
		StartedAt:  now,
	}

	bjGamesMu.Lock()
	bjGames[userID] = pending
	bjGamesMu.Unlock()

	return BlackjackStartResult{
		Record:   record.copy(),
		Player:   BlackjackHand{Cards: pCards, Value: pVal, CardStr: cardsStr(pCards)},
		Dealer:   BlackjackHand{Cards: dCards[:1], Value: handValue(dCards[:1]), CardStr: cardsStr(dCards[:1]) + " ??"},
		Wager:    wager,
		GameOver: false,
	}, nil
}

func (s *Store) BlackjackHit(userID snowflake.ID, now time.Time) (BlackjackTurnResult, error) {
	bjGamesMu.Lock()
	pending := bjGames[userID]
	if pending == nil {
		bjGamesMu.Unlock()
		return BlackjackTurnResult{}, fmt.Errorf("you don't have an active blackjack game")
	}
	delete(bjGames, userID)
	bjGamesMu.Unlock()

	pending.mu.Lock()
	defer pending.mu.Unlock()

	pending.Player = append(pending.Player, drawCard())
	pVal := handValue(pending.Player)

	if pVal > 21 {
		// Bust
		s.mu.Lock()
		record := s.ensureRecord(userID, pending.Name)
		record.BlackjackLosses++
		record.UpdatedAt = now
		s.save(record)
		s.mu.Unlock()

		return BlackjackTurnResult{
			Record:   record.copy(),
			Player:   BlackjackHand{Cards: pending.Player, Value: pVal, Bust: true, CardStr: cardsStr(pending.Player)},
			Dealer:   BlackjackHand{Cards: pending.Dealer, Value: handValue(pending.Dealer), CardStr: cardsStr(pending.Dealer)},
			Wager:    pending.Wager,
			Payout:   -pending.Wager,
			GameOver: true,
		}, nil
	}

	// Store back for next move
	bjGamesMu.Lock()
	bjGames[userID] = pending
	bjGamesMu.Unlock()

	return BlackjackTurnResult{
		Record:   s.getRecord(userID),
		Player:   BlackjackHand{Cards: pending.Player, Value: pVal, CardStr: cardsStr(pending.Player)},
		Dealer:   BlackjackHand{Cards: pending.Dealer[:1], Value: handValue(pending.Dealer[:1]), CardStr: cardsStr(pending.Dealer[:1]) + " ??"},
		Wager:    pending.Wager,
		GameOver: false,
	}, nil
}

func (s *Store) BlackjackStand(userID snowflake.ID, now time.Time) (BlackjackTurnResult, error) {
	bjGamesMu.Lock()
	pending := bjGames[userID]
	if pending == nil {
		bjGamesMu.Unlock()
		return BlackjackTurnResult{}, fmt.Errorf("you don't have an active blackjack game")
	}
	delete(bjGames, userID)
	bjGamesMu.Unlock()

	pending.mu.Lock()
	defer pending.mu.Unlock()

	// Dealer plays
	standVal := int(pending.Difficulty)
	for handValue(pending.Dealer) < standVal {
		pending.Dealer = append(pending.Dealer, drawCard())
	}

	pVal := handValue(pending.Player)
	dVal := handValue(pending.Dealer)

	s.mu.Lock()
	record := s.ensureRecord(userID, pending.Name)
	var payout int
	won := false
	push := false

	if dVal > 21 {
		payout = pending.Wager
		won = true
		record.Balance += payout
		record.Earned += payout
		record.BlackjackWins++
	} else if pVal > dVal {
		payout = pending.Wager
		won = true
		record.Balance += payout
		record.Earned += payout
		record.BlackjackWins++
	} else if pVal == dVal {
		payout = pending.Wager
		record.Balance += payout
		push = true
	} else {
		payout = -pending.Wager
		record.BlackjackLosses++
	}
	record.UpdatedAt = now
	s.save(record)
	s.mu.Unlock()

	return BlackjackTurnResult{
		Record:   record.copy(),
		Player:   BlackjackHand{Cards: pending.Player, Value: pVal, CardStr: cardsStr(pending.Player)},
		Dealer:   BlackjackHand{Cards: pending.Dealer, Value: dVal, CardStr: cardsStr(pending.Dealer)},
		Wager:    pending.Wager,
		Won:      won,
		Payout:   payout,
		Push:     push,
		GameOver: true,
	}, nil
}

func (s *Store) getRecord(userID snowflake.ID) Record {
	s.mu.RLock()
	defer s.mu.RUnlock()
	record := s.records[userID]
	if record != nil {
		return *record
	}
	return Record{}
}

func drawCard() int {
	return 1 + rand.Intn(13)
}

func handValue(cards []int) int {
	total := 0
	aces := 0
	for _, c := range cards {
		if c == 1 {
			aces++
			total += 11
		} else if c > 10 {
			total += 10
		} else {
			total += c
		}
	}
	for aces > 0 && total > 21 {
		total -= 10
		aces--
	}
	return total
}

func cardStr(card int) string {
	switch card {
	case 1:
		return "A"
	case 11:
		return "J"
	case 12:
		return "Q"
	case 13:
		return "K"
	default:
		if card == 10 {
			return "10"
		}
		return string(rune('0' + card))
	}
}

func cardsStr(cards []int) string {
	s := ""
	for i, c := range cards {
		if i > 0 {
			s += " "
		}
		s += cardStr(c)
	}
	return s
}
