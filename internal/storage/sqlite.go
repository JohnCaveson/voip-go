package storage

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/voip-app/pkg/models"

	_ "modernc.org/sqlite"
)

type SQLiteStorage struct {
	db *sql.DB
}

func NewSQLiteStorage(dbPath string) (*SQLiteStorage, error) {
	db, err := sql.Open("sqlite", dbPath)
	if err != nil {
		return nil, fmt.Errorf("open sqlite: %w", err)
	}

	db.SetMaxOpenConns(1)

	if err := migrate(db); err != nil {
		return nil, fmt.Errorf("migrate: %w", err)
	}

	return &SQLiteStorage{db: db}, nil
}

func migrate(db *sql.DB) error {
	queries := []string{
		`CREATE TABLE IF NOT EXISTS channels (
			id TEXT PRIMARY KEY,
			name TEXT NOT NULL,
			type TEXT NOT NULL CHECK(type IN ('text','voice')),
			is_default INTEGER NOT NULL DEFAULT 0,
			created_at TEXT NOT NULL,
			deleted INTEGER NOT NULL DEFAULT 0
		)`,
		`CREATE TABLE IF NOT EXISTS messages (
			id TEXT PRIMARY KEY,
			channel_id TEXT NOT NULL REFERENCES channels(id),
			user_id TEXT NOT NULL,
			username TEXT NOT NULL,
			content TEXT NOT NULL,
			created_at TEXT NOT NULL
		)`,
		`CREATE TABLE IF NOT EXISTS users (
			id TEXT PRIMARY KEY,
			username TEXT NOT NULL,
			joined_at TEXT NOT NULL,
			is_online INTEGER NOT NULL DEFAULT 0
		)`,
	}

	for _, q := range queries {
		if _, err := db.Exec(q); err != nil {
			return err
		}
	}

	return nil
}

func (s *SQLiteStorage) Close() error {
	return s.db.Close()
}

func (s *SQLiteStorage) CreateChannel(ctx context.Context, ch *models.Channel) error {
	_, err := s.db.ExecContext(ctx,
		`INSERT INTO channels (id, name, type, is_default, created_at, deleted) VALUES (?, ?, ?, ?, ?, ?)`,
		ch.ID, ch.Name, string(ch.Type), boolToInt(ch.IsDefault), ch.CreatedAt.Format(time.RFC3339), boolToInt(ch.Deleted),
	)
	return err
}

func (s *SQLiteStorage) GetChannel(ctx context.Context, id string) (*models.Channel, error) {
	row := s.db.QueryRowContext(ctx,
		`SELECT id, name, type, is_default, created_at, deleted FROM channels WHERE id = ? AND deleted = 0`, id,
	)

	ch := &models.Channel{}
	var createdAt string
	var deletedInt int
	err := row.Scan(&ch.ID, &ch.Name, (*string)(&ch.Type), &ch.IsDefault, &createdAt, &deletedInt)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("channel not found: %s", id)
		}
		return nil, err
	}

	ch.CreatedAt, err = time.Parse(time.RFC3339, createdAt)
	if err != nil {
		return nil, fmt.Errorf("parse channel created_at: %w", err)
	}
	ch.Deleted = intToBool(deletedInt)
	return ch, nil
}

func (s *SQLiteStorage) ListChannels(ctx context.Context) ([]*models.Channel, error) {
	rows, err := s.db.QueryContext(ctx,
		`SELECT id, name, type, is_default, created_at, deleted FROM channels WHERE deleted = 0 ORDER BY created_at`,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var channels []*models.Channel
	for rows.Next() {
		ch := &models.Channel{}
		var createdAt string
		var deletedInt int
		if err := rows.Scan(&ch.ID, &ch.Name, (*string)(&ch.Type), &ch.IsDefault, &createdAt, &deletedInt); err != nil {
			return nil, err
		}
		ch.CreatedAt, err = time.Parse(time.RFC3339, createdAt)
		if err != nil {
			return nil, fmt.Errorf("parse channel created_at: %w", err)
		}
		ch.Deleted = intToBool(deletedInt)
		channels = append(channels, ch)
	}

	return channels, rows.Err()
}

func (s *SQLiteStorage) UpdateChannel(ctx context.Context, ch *models.Channel) error {
	_, err := s.db.ExecContext(ctx,
		`UPDATE channels SET name = ?, type = ?, is_default = ?, deleted = ? WHERE id = ?`,
		ch.Name, string(ch.Type), boolToInt(ch.IsDefault), boolToInt(ch.Deleted), ch.ID,
	)
	return err
}

func (s *SQLiteStorage) DeleteChannel(ctx context.Context, id string) error {
	_, err := s.db.ExecContext(ctx,
		`UPDATE channels SET deleted = 1 WHERE id = ?`, id,
	)
	return err
}

func (s *SQLiteStorage) SendMessage(ctx context.Context, msg *models.Message) error {
	_, err := s.db.ExecContext(ctx,
		`INSERT INTO messages (id, channel_id, user_id, username, content, created_at) VALUES (?, ?, ?, ?, ?, ?)`,
		msg.ID, msg.ChannelID, msg.UserID, msg.Username, msg.Content, msg.CreatedAt.Format(time.RFC3339),
	)
	return err
}

func (s *SQLiteStorage) ListMessages(ctx context.Context, channelID string, limit, offset int) ([]*models.Message, error) {
	if limit <= 0 {
		limit = 50
	}

	rows, err := s.db.QueryContext(ctx,
		`SELECT id, channel_id, user_id, username, content, created_at FROM messages WHERE channel_id = ? ORDER BY created_at DESC LIMIT ? OFFSET ?`,
		channelID, limit, offset,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var messages []*models.Message
	for rows.Next() {
		msg := &models.Message{}
		var createdAt string
		if err := rows.Scan(&msg.ID, &msg.ChannelID, &msg.UserID, &msg.Username, &msg.Content, &createdAt); err != nil {
			return nil, err
		}
		msg.CreatedAt, err = time.Parse(time.RFC3339, createdAt)
		if err != nil {
			return nil, fmt.Errorf("parse message created_at: %w", err)
		}
		messages = append(messages, msg)
	}

	return messages, rows.Err()
}

func (s *SQLiteStorage) DeleteMessage(ctx context.Context, id string) error {
	_, err := s.db.ExecContext(ctx,
		`DELETE FROM messages WHERE id = ?`, id,
	)
	return err
}

func (s *SQLiteStorage) AddUser(ctx context.Context, user *models.User) error {
	_, err := s.db.ExecContext(ctx,
		`INSERT INTO users (id, username, joined_at, is_online) VALUES (?, ?, ?, ?)`,
		user.ID, user.Username, user.JoinedAt.Format(time.RFC3339), boolToInt(user.IsOnline),
	)
	return err
}

func (s *SQLiteStorage) GetUser(ctx context.Context, id string) (*models.User, error) {
	row := s.db.QueryRowContext(ctx,
		`SELECT id, username, joined_at, is_online FROM users WHERE id = ?`, id,
	)

	user := &models.User{}
	var joinedAt string
	err := row.Scan(&user.ID, &user.Username, &joinedAt, &user.IsOnline)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("user not found: %s", id)
		}
		return nil, err
	}

	user.JoinedAt, err = time.Parse(time.RFC3339, joinedAt)
	if err != nil {
		return nil, fmt.Errorf("parse user joined_at: %w", err)
	}
	return user, nil
}

func (s *SQLiteStorage) ListUsers(ctx context.Context) ([]*models.User, error) {
	rows, err := s.db.QueryContext(ctx,
		`SELECT id, username, joined_at, is_online FROM users ORDER BY joined_at`,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var users []*models.User
	for rows.Next() {
		user := &models.User{}
		var joinedAt string
		if err := rows.Scan(&user.ID, &user.Username, &joinedAt, &user.IsOnline); err != nil {
			return nil, err
		}
		user.JoinedAt, err = time.Parse(time.RFC3339, joinedAt)
		if err != nil {
			return nil, fmt.Errorf("parse user joined_at: %w", err)
		}
		users = append(users, user)
	}

	return users, rows.Err()
}

func (s *SQLiteStorage) SetUserOnline(ctx context.Context, id string, online bool) error {
	_, err := s.db.ExecContext(ctx,
		`UPDATE users SET is_online = ? WHERE id = ?`, boolToInt(online), id,
	)
	return err
}

func boolToInt(b bool) int {
	if b {
		return 1
	}
	return 0
}

func intToBool(i int) bool {
	return i == 1
}
