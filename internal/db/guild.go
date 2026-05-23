package db

import (
	"context"
	"database/sql"
)

type GuildConfig struct {
	GuildID      string
	LogChannelID *string
}

func (db *Database) GetGuildConfig(ctx context.Context, guildID string) (*GuildConfig, error) {
	var cfg GuildConfig
	err := db.Conn.QueryRowContext(ctx, "SELECT guild_id, log_channel_id FROM guild_configs WHERE guild_id = $1", guildID).Scan(&cfg.GuildID, &cfg.LogChannelID)
	if err == sql.ErrNoRows {
		_, err = db.Conn.ExecContext(ctx, "INSERT INTO guild_configs (guild_id) VALUES ($1)", guildID)
		if err != nil {
			return nil, err
		}
		return &GuildConfig{GuildID: guildID}, nil
	}
	return &cfg, err
}

func (db *Database) SetLogChannel(ctx context.Context, guildID string, channelID string) error {
	_, err := db.Conn.ExecContext(ctx, `
		INSERT INTO guild_configs (guild_id, log_channel_id) VALUES ($1, $2)
		ON CONFLICT (guild_id) DO UPDATE SET log_channel_id = $2`, guildID, channelID)
	return err
}

// GetCommandsRunCount returns the total number of commands run from the command_logs table
func (db *Database) GetCommandsRunCount(ctx context.Context) (int, error) {
	var count int
	err := db.Conn.QueryRowContext(ctx, "SELECT COUNT(*) FROM command_logs").Scan(&count)
	if err != nil {
		return 0, err
	}
	return count, nil
}

// IncrementCommandRun logs a command execution for tracking
func (db *Database) IncrementCommandRun(ctx context.Context, guildID, userID, commandName string) error {
	_, err := db.Conn.ExecContext(ctx, `
		INSERT INTO command_logs (guild_id, user_id, command_name) VALUES ($1, $2, $3)`, guildID, userID, commandName)
	return err
}
