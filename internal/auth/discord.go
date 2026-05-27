package auth

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"crypto/rand"

	"github.com/watispro5212/PulseKeep/internal/config"
)

var (
	discordAPIURL = "https://discord.com/api/v10"
)

// DiscordUser represents a Discord user
type DiscordUser struct {
	ID            string `json:"id"`
	Username      string `json:"username"`
	Discriminator string `json:"discriminator"`
	Avatar        string `json:"avatar,omitempty"`
	Bot           bool   `json:"bot,omitempty"`
	System        bool   `json:"system,omitempty"`
	MFAEnabled    bool   `json:"mfa_enabled,omitempty"`
	Banner        string `json:"banner,omitempty"`
	AccentColor   int    `json:"accent_color,omitempty"`
	Locale        string `json:"locale,omitempty"`
	Verified      bool   `json:"verified,omitempty"`
	Email         string `json:"email,omitempty"`
	Flags         int    `json:"flags,omitempty"`
	PremiumType   int    `json:"premium_type,omitempty"`
	PublicFlags   int    `json:"public_flags,omitempty"`
}

// DiscordGuild represents a Discord guild (server)
type DiscordGuild struct {
	ID          string `json:"id"`
	Name        string `json:"name"`
	Icon        string `json:"icon,omitempty"`
	IconHash    string `json:"icon_hash,omitempty"`
	Owner       bool   `json:"owner,omitempty"`
	Permissions int    `json:"permissions,omitempty"`
	Features    []string `json:"features,omitempty"`
}

// ExchangeCode exchanges an authorization code for an access token
func ExchangeCode(cfg *config.Config, code string) (string, error) {
	data := url.Values{}
	data.Set("client_id", cfg.DiscordClientID)
	data.Set("client_secret", cfg.DiscordClientSecret)
	data.Set("grant_type", "authorization_code")
	data.Set("code", code)
	data.Set("redirect_uri", cfg.DiscordRedirectURI)
	data.Set("scope", "identify guilds")

	resp, err := http.PostForm(discordAPIURL+"/oauth2/token", data)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("discord token exchange failed: %d", resp.StatusCode)
	}

	var tokenResp struct {
		AccessToken string `json:"access_token"`
		TokenType   string `json:"token_type"`
		ExpiresIn   int    `json:"expires_in"`
		RefreshToken string `json:"refresh_token"`
		Scope       string `json:"scope"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&tokenResp); err != nil {
		return "", err
	}

	return tokenResp.AccessToken, nil
}

// GetUser fetches the authenticated user from Discord
func GetUser(accessToken string) (*DiscordUser, error) {
	req, err := http.NewRequest("GET", discordAPIURL+"/users/@me", nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Authorization", "Bearer "+accessToken)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("failed to fetch user: %d", resp.StatusCode)
	}

	var user DiscordUser
	if err := json.NewDecoder(resp.Body).Decode(&user); err != nil {
		return nil, err
	}

	return &user, nil
}

// GetUserGuilds fetches the guilds the user is a member of
func GetUserGuilds(accessToken string) ([]DiscordGuild, error) {
	req, err := http.NewRequest("GET", discordAPIURL+"/users/@me/guilds", nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Authorization", "Bearer "+accessToken)
	q := req.URL.Query()
	q.Set("limit", "100") // Max allowed by Discord
	req.URL.RawQuery = q.Encode()

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("failed to fetch guilds: %d", resp.StatusCode)
	}

	var guilds []DiscordGuild
	if err := json.NewDecoder(resp.Body).Decode(&guilds); err != nil {
		return nil, err
	}

	return guilds, nil
}

// State generates a random state string for OAuth2 CSRF protection
func State() string {
	b := make([]byte, 32)
	if _, err := rand.Read(b); err != nil {
		return base64.URLEncoding.EncodeToString(b[:16]) // fallback
	}
	return base64.URLEncoding.EncodeToString(b)
}