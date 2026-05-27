package automod

import (
	"regexp"
	"strings"
	"sync"
	"time"
	"unicode"
)

var urlRegex = regexp.MustCompile(`https?://[^\s]+`)

type spamTracker struct {
	mu     sync.Mutex
	times  map[string][]time.Time
}

type ActionType string

const (
	ActionNone   ActionType = "none"
	ActionWarn   ActionType = "warn"
	ActionDelete ActionType = "delete"
	ActionTimeout ActionType = "timeout"
)

type RuleResult struct {
	Action      ActionType
	Reason      string
	DeleteMsg   bool
}

type Engine struct {
	spam   *spamTracker
	config *ConfigStore
}

func NewEngine(cfgStore *ConfigStore) *Engine {
	return &Engine{
		spam:   &spamTracker{times: make(map[string][]time.Time)},
		config: cfgStore,
	}
}

func (e *Engine) CheckMessage(guildID, userID, content string) RuleResult {
	cfg := e.config.Get(guildID)
	if cfg == nil || !cfg.AutomodEnabled {
		return RuleResult{}
	}

	// 1. Spam detection
	if cfg.SpamEnabled {
		if result := e.checkSpam(guildID, userID, cfg); result.Action != ActionNone {
			return result
		}
	}

	// 2. Mass mention protection
	if cfg.MentionEnabled && countMentions(content) > cfg.MentionMax {
		action := parseAction(cfg.MentionAction)
		return RuleResult{Action: action, Reason: "Mass mention detected", DeleteMsg: action == ActionDelete}
	}

	// 3. Keyword filter
	if cfg.BannedWords != "" {
		if word := checkBannedWords(content, cfg.BannedWords); word != "" {
			return RuleResult{Action: ActionDelete, Reason: "Message contained banned word: " + word, DeleteMsg: true}
		}
	}

	// 4. Link protection
	if cfg.LinksEnabled && urlRegex.MatchString(content) {
		action := parseAction(cfg.LinksAction)
		return RuleResult{Action: action, Reason: "Links are not allowed", DeleteMsg: action == ActionDelete}
	}

	// 5. Caps lock protection
	if cfg.CapsEnabled && len(content) >= cfg.CapsMinLength && isMostlyCaps(content, cfg.CapsPercent) {
		action := parseAction(cfg.CapsAction)
		return RuleResult{Action: action, Reason: "Excessive caps lock", DeleteMsg: action == ActionDelete}
	}

	return RuleResult{}
}

func (e *Engine) checkSpam(guildID, userID string, cfg *GuildConfig) RuleResult {
	key := guildID + ":" + userID
	now := time.Now()

	e.spam.mu.Lock()
	times := e.spam.times[key]

	// prune old entries
	cutoff := now.Add(-time.Duration(cfg.SpamWindowSecs) * time.Second)
	var recent []time.Time
	for _, t := range times {
		if t.After(cutoff) {
			recent = append(recent, t)
		}
	}
	recent = append(recent, now)
	e.spam.times[key] = recent
	e.spam.mu.Unlock()

	if len(recent) > cfg.SpamMaxMessages {
		action := parseAction(cfg.SpamAction)
		if action == ActionWarn {
			action = ActionTimeout
		}
		return RuleResult{Action: action, Reason: "Spam detected", DeleteMsg: true}
	}

	return RuleResult{}
}

func parseAction(s string) ActionType {
	switch strings.ToLower(s) {
	case "warn":
		return ActionWarn
	case "delete":
		return ActionDelete
	case "timeout":
		return ActionTimeout
	default:
		return ActionWarn
	}
}

func countMentions(s string) int {
	count := 0
	for i := 0; i < len(s); i++ {
		if s[i] == '@' {
			count++
		}
	}
	return count
}

func checkBannedWords(content, bannedStr string) string {
	if bannedStr == "" {
		return ""
	}
	lower := strings.ToLower(content)
	for _, word := range strings.Split(bannedStr, ",") {
		word = strings.TrimSpace(word)
		if word == "" {
			continue
		}
		if strings.Contains(lower, strings.ToLower(word)) {
			return word
		}
	}
	return ""
}

func isMostlyCaps(s string, percent int) bool {
	upper := 0
	letters := 0
	for _, r := range s {
		if unicode.IsLetter(r) {
			letters++
			if unicode.IsUpper(r) {
				upper++
			}
		}
	}
	if letters == 0 {
		return false
	}
	return upper*100/letters >= percent
}
