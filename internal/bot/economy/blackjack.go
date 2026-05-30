package economy

import (
	"math/rand"
	"time"

	"github.com/disgoorg/snowflake/v2"
)

type BlackjackDifficulty int

const (
	BlackjackEasy   BlackjackDifficulty = 14
	BlackjackNormal BlackjackDifficulty = 17
	BlackjackHard   BlackjackDifficulty = 19
	BlackjackExpert BlackjackDifficulty = 20
)

type BlackjackHand struct {
	Cards    []int
	Value    int
	Bust     bool
	CardStr  string
}

type BlackjackResult struct {
	Record  Record
	Wager   int
	Won     bool
	Payout  int
	Player  BlackjackHand
	Dealer  BlackjackHand
	Natural bool
	Push    bool
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

func (s *Store) Blackjack(userID snowflake.ID, name string, wager int, difficulty BlackjackDifficulty, now time.Time) (BlackjackResult, error) {
	if wager <= 0 {
		return BlackjackResult{}, ErrInvalidAmount
	}

	s.mu.Lock()
	defer s.mu.Unlock()

	record := s.ensureRecord(userID, name)
	if record.Balance < wager {
		return BlackjackResult{}, ErrInsufficientFund
	}

	playerCards := []int{drawCard(), drawCard()}
	playerValue := handValue(playerCards)

	dealerCards := []int{drawCard(), drawCard()}
	dealerValue := handValue(dealerCards)

	won := false
	push := false
	natural := false
	netPayout := 0

	// check for natural 21s
	if playerValue == 21 && dealerValue == 21 {
		push = true
	} else if playerValue == 21 {
		won = true
		natural = true
	} else if dealerValue == 21 {
		// dealer natural, player loses
		netPayout = -wager
	} else {
		standValue := int(difficulty)
		for dealerValue < standValue {
			dealerCards = append(dealerCards, drawCard())
			dealerValue = handValue(dealerCards)
		}

		if dealerValue > 21 {
			won = true
			netPayout = wager
		} else if playerValue > dealerValue {
			won = true
			netPayout = wager
		} else if playerValue == dealerValue {
			push = true
		} else {
			netPayout = -wager
		}
	}

	if natural {
		netPayout = wager * 3 / 2
	}

	if netPayout > 0 {
		record.Balance += netPayout
		record.Earned += netPayout
		record.BlackjackWins++
	} else if netPayout < 0 {
		record.Balance += netPayout
		record.Spent += -netPayout
		record.BlackjackLosses++
	}
	record.UpdatedAt = now
	s.save(record)

	return BlackjackResult{
		Record:  record.copy(),
		Wager:   wager,
		Won:     won,
		Payout:  netPayout,
		Player:  BlackjackHand{Cards: playerCards, Value: playerValue, CardStr: cardsStr(playerCards)},
		Dealer:  BlackjackHand{Cards: dealerCards, Value: dealerValue, CardStr: cardsStr(dealerCards)},
		Natural: natural,
		Push:    push,
	}, nil
}
