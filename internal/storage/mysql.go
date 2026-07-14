package storage

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/voip-app/internal/models"

	_ "github.com/go-sql-driver/mysql"
)

type MySQLStorage struct {
	db *sql.DB
}

func NewMySQLStorage(dsn string) (*MySQLStorage, error) {
	db, err := sql.Open("mysql", dsn)
	if err != nil {
		return nil, fmt.Errorf("open mysql: %w", err)
	}

	db.SetMaxOpenConns(25)
	db.SetMaxIdleConns(5)

	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("ping mysql: %w", err)
	}

	if err := migrateMySQL(db); err != nil {
		return nil, fmt.Errorf("migrate: %w", err)
	}

	return &MySQLStorage{db: db}, nil
}

func migrateMySQL(db *sql.DB) error {
	queries := []string{
		`CREATE TABLE IF NOT EXISTS channels (
			id VARCHAR(64) PRIMARY KEY,
			name VARCHAR(255) NOT NULL,
			type VARCHAR(10) NOT NULL CHECK(type IN ('text','voice')),
			is_default BOOLEAN NOT NULL DEFAULT FALSE,
			created_at DATETIME(3) NOT NULL,
			deleted BOOLEAN NOT NULL DEFAULT FALSE
		)`,
		`CREATE TABLE IF NOT EXISTS messages (
			id VARCHAR(64) PRIMARY KEY,
			channel_id VARCHAR(64) NOT NULL,
			user_id VARCHAR(64) NOT NULL,
			username VARCHAR(255) NOT NULL,
			content TEXT NOT NULL,
			created_at DATETIME(3) NOT NULL,
			FOREIGN KEY (channel_id) REFERENCES channels(id)
		)`,
		`CREATE TABLE IF NOT EXISTS users (
			id VARCHAR(64) PRIMARY KEY,
			username VARCHAR(255) NOT NULL,
			joined_at DATETIME(3) NOT NULL,
			is_online BOOLEAN NOT NULL DEFAULT FALSE
		)`,
	}

	for _, q := range queries {
		if _, err := db.Exec(q); err != nil {
			return err
		}
	}

	return nil
}

func (s *MySQLStorage) Close() error {
	return s.db.Close()
}

func (s *MySQLStorage) CreateChannel(ctx context.Context, ch *models.Channel) error {
	_, err := s.db.ExecContext(ctx,
		`INSERT INTO channels (id, name, type, is_default, created_at, deleted) VALUES (?, ?, ?, ?, ?, ?)`,
		ch.ID, ch.Name, string(ch.Type), ch.IsDefault, ch.CreatedAt, ch.Deleted,
	)
	return err
}

func (s *MySQLStorage) GetChannel(ctx context.Context, id string) (*models.Channel, error) {
	row := s.db.QueryRowContext(ctx,
		`SELECT id, name, type, is_default, created_at, deleted FROM channels WHERE id = ? AND deleted = FALSE`, id,
	)

	ch := &models.Channel{}
	err := row.Scan(&ch.ID, &ch.Name, &ch.Type, &ch.IsDefault, &ch.CreatedAt, &ch.Deleted)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("channel not found: %s", id)
		}
		return nil, err
	}

	return ch, nil
}

func (s *MySQLStorage) ListChannels(ctx context.Context) ([]*models.Channel, error) {
	rows, err := s.db.QueryContext(ctx,
		`SELECT id, name, type, is_default, created_at, deleted FROM channels WHERE deleted = FALSE ORDER BY created_at`,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var channels []*models.Channel
	for rows.Next() {
		ch := &models.Channel{}
		if err := rows.Scan(&ch.ID, &ch.Name, &ch.Type, &ch.IsDefault, &ch.CreatedAt, &ch.Deleted); err != nil {
			return nil, err
		}
		channels = append(channels, ch)
	}

	return channels, rows.Err()
}

func (s *MySQLStorage) UpdateChannel(ctx context.Context, ch *models.Channel) error {
	_, err := s.db.ExecContext(ctx,
		`UPDATE channels SET name = ?, type = ?, is_default = ?, deleted = ? WHERE id = ?`,
		ch.Name, string(ch.Type), ch.IsDefault, ch.Deleted, ch.ID,
	)
	return err
}

func (s *MySQLStorage) DeleteChannel(ctx context.Context, id string) error {
	_, err := s.db.ExecContext(ctx,
		`UPDATE channels SET deleted = TRUE WHERE id = ?`, id,
	)
	return err
}

func (s *MySQLStorage) SendMessage(ctx context.Context, msg *models.Message) error {
	_, err := s.db.ExecContext(ctx,
		`INSERT INTO messages (id, channel_id, user_id, username, content, created_at) VALUES (?, ?, ?, ?, ?, ?)`,
		msg.ID, msg.ChannelID, msg.UserID, msg.Username, msg.Content, msg.CreatedAt,
	)
	return err
}

func (s *MySQLStorage) ListMessages(ctx context.Context, channelID string, limit, offset int) ([]*models.Message, error) {
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
		if err := rows.Scan(&msg.ID, &msg.ChannelID, &msg.UserID, &msg.Username, &msg.Content, &msg.CreatedAt); err != nil {
			return nil, err
		}
		messages = append(messages, msg)
	}

	return messages, rows.Err()
}

func (s *MySQLStorage) DeleteMessage(ctx context.Context, id string) error {
	_, err := s.db.ExecContext(ctx,
		`DELETE FROM messages WHERE id = ?`, id,
	)
	return err
}

func (s *MySQLStorage) AddUser(ctx context.Context, user *models.User) error {
	_, err := s.db.ExecContext(ctx,
		`INSERT INTO users (id, username, joined_at, is_online) VALUES (?, ?, ?, ?)`,
		user.ID, user.Username, user.JoinedAt, user.IsOnline,
	)
	return err
}

func (s *MySQLStorage) GetUser(ctx context.Context, id string) (*models.User, error) {
	row := s.db.QueryRowContext(ctx,
		`SELECT id, username, joined_at, is_online FROM users WHERE id = ?`, id,
	)

	user := &models.User{}
	err := row.Scan(&user.ID, &user.Username, &user.JoinedAt, &user.IsOnline)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("user not found: %s", id)
		}
		return nil, err
	}

	return user, nil
}

func (s *MySQLStorage) ListUsers(ctx context.Context) ([]*models.User, error) {
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
		if err := rows.Scan(&user.ID, &user.Username, &user.JoinedAt, &user.IsOnline); err != nil {
			return nil, err
		}
		users = append(users, user)
	}

	return users, rows.Err()
}

func (s *MySQLStorage) SetUserOnline(ctx context.Context, id string, online bool) error {
	_, err := s.db.ExecContext(ctx,
		`UPDATE users SET is_online = ? WHERE id = ?`, online, id,
	)
	return err
}

var _ Storage = (*MySQLStorage)(nil)
