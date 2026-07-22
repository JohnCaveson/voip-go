package storage

import (
	"context"
	"fmt"
	"time"

	"github.com/voip-app/internal/models"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type MongoDBStorage struct {
	client   *mongo.Client
	database *mongo.Database
}

func NewMongoDBStorage(uri string) (*MongoDBStorage, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	client, err := mongo.Connect(ctx, options.Client().ApplyURI(uri))
	if err != nil {
		return nil, fmt.Errorf("connect mongo: %w", err)
	}

	if err := client.Ping(ctx, nil); err != nil {
		return nil, fmt.Errorf("ping mongo: %w", err)
	}

	db := client.Database("voip")

	if err := ensureIndexes(db); err != nil {
		return nil, fmt.Errorf("ensure indexes: %w", err)
	}

	return &MongoDBStorage{
		client:   client,
		database: db,
	}, nil
}

func ensureIndexes(db *mongo.Database) error {
	ctx := context.Background()

	channelsCol := db.Collection("channels")
	_, err := channelsCol.Indexes().CreateOne(ctx, mongo.IndexModel{
		Keys:    bson.D{{Key: "name", Value: 1}},
		Options: options.Index().SetUnique(true),
	})
	if err != nil {
		return fmt.Errorf("channels name index: %w", err)
	}

	messagesCol := db.Collection("messages")
	_, err = messagesCol.Indexes().CreateOne(ctx, mongo.IndexModel{
		Keys: bson.D{{Key: "channel_id", Value: 1}, {Key: "created_at", Value: -1}},
	})
	if err != nil {
		return fmt.Errorf("messages index: %w", err)
	}

	usersCol := db.Collection("users")
	_, err = usersCol.Indexes().CreateOne(ctx, mongo.IndexModel{
		Keys:    bson.D{{Key: "username", Value: 1}},
		Options: options.Index().SetUnique(true),
	})
	if err != nil {
		return fmt.Errorf("users username index: %w", err)
	}

	return nil
}

func (s *MongoDBStorage) Close() error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	return s.client.Disconnect(ctx)
}

func (s *MongoDBStorage) CreateChannel(ctx context.Context, ch *models.Channel) error {
	doc := bson.M{
		"id":         ch.ID,
		"name":       ch.Name,
		"type":       string(ch.Type),
		"is_default": ch.IsDefault,
		"created_at": ch.CreatedAt,
		"deleted":    ch.Deleted,
	}
	_, err := s.database.Collection("channels").InsertOne(ctx, doc)
	return err
}

func (s *MongoDBStorage) GetChannel(ctx context.Context, id string) (*models.Channel, error) {
	var doc bson.M
	err := s.database.Collection("channels").FindOne(ctx, bson.M{
		"id":      id,
		"deleted": false,
	}).Decode(&doc)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, fmt.Errorf("channel not found: %s", id)
		}
		return nil, err
	}

	ch := &models.Channel{
		ID:        doc["id"].(string),
		Name:      doc["name"].(string),
		Type:      models.ChannelType(doc["type"].(string)),
		IsDefault: doc["is_default"].(bool),
		Deleted:   doc["deleted"].(bool),
	}
	if t, ok := doc["created_at"].(primitive.DateTime); ok {
		ch.CreatedAt = t.Time()
	}

	return ch, nil
}

func (s *MongoDBStorage) ListChannels(ctx context.Context) ([]*models.Channel, error) {
	cursor, err := s.database.Collection("channels").Find(ctx, bson.M{"deleted": false},
		options.Find().SetSort(bson.D{{Key: "created_at", Value: 1}}),
	)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var channels []*models.Channel
	for cursor.Next(ctx) {
		var doc bson.M
		if err := cursor.Decode(&doc); err != nil {
			return nil, err
		}

		ch := &models.Channel{
			ID:        doc["id"].(string),
			Name:      doc["name"].(string),
			Type:      models.ChannelType(doc["type"].(string)),
			IsDefault: doc["is_default"].(bool),
			Deleted:   doc["deleted"].(bool),
		}
		if t, ok := doc["created_at"].(primitive.DateTime); ok {
			ch.CreatedAt = t.Time()
		}
		channels = append(channels, ch)
	}

	return channels, cursor.Err()
}

func (s *MongoDBStorage) UpdateChannel(ctx context.Context, ch *models.Channel) error {
	_, err := s.database.Collection("channels").UpdateOne(ctx,
		bson.M{"id": ch.ID},
		bson.M{"$set": bson.M{
			"name":       ch.Name,
			"type":       string(ch.Type),
			"is_default": ch.IsDefault,
			"deleted":    ch.Deleted,
		}},
	)
	return err
}

func (s *MongoDBStorage) DeleteChannel(ctx context.Context, id string) error {
	_, err := s.database.Collection("channels").UpdateOne(ctx,
		bson.M{"id": id},
		bson.M{"$set": bson.M{"deleted": true}},
	)
	return err
}

func (s *MongoDBStorage) SendMessage(ctx context.Context, msg *models.Message) error {
	doc := bson.M{
		"id":         msg.ID,
		"channel_id": msg.ChannelID,
		"user_id":    msg.UserID,
		"username":   msg.Username,
		"content":    msg.Content,
		"created_at": msg.CreatedAt,
	}
	_, err := s.database.Collection("messages").InsertOne(ctx, doc)
	return err
}

func (s *MongoDBStorage) ListMessages(ctx context.Context, channelID string, limit, offset int) ([]*models.Message, error) {
	if limit <= 0 {
		limit = 50
	}

	opts := options.Find().
		SetSort(bson.D{{Key: "created_at", Value: -1}}).
		SetLimit(int64(limit)).
		SetSkip(int64(offset))

	cursor, err := s.database.Collection("messages").Find(ctx, bson.M{"channel_id": channelID}, opts)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var messages []*models.Message
	for cursor.Next(ctx) {
		var doc bson.M
		if err := cursor.Decode(&doc); err != nil {
			return nil, err
		}

		msg := &models.Message{
			ID:        doc["id"].(string),
			ChannelID: doc["channel_id"].(string),
			UserID:    doc["user_id"].(string),
			Username:  doc["username"].(string),
			Content:   doc["content"].(string),
		}
		if t, ok := doc["created_at"].(primitive.DateTime); ok {
			msg.CreatedAt = t.Time()
		}
		messages = append(messages, msg)
	}

	return messages, cursor.Err()
}

func (s *MongoDBStorage) DeleteMessage(ctx context.Context, id string) error {
	_, err := s.database.Collection("messages").DeleteOne(ctx, bson.M{"id": id})
	return err
}

func (s *MongoDBStorage) AddUser(ctx context.Context, user *models.User) error {
	doc := bson.M{
		"id":        user.ID,
		"username":  user.Username,
		"joined_at": user.JoinedAt,
		"is_online": user.IsOnline,
	}
	_, err := s.database.Collection("users").InsertOne(ctx, doc)
	return err
}

func (s *MongoDBStorage) GetUser(ctx context.Context, id string) (*models.User, error) {
	var doc bson.M
	err := s.database.Collection("users").FindOne(ctx, bson.M{"id": id}).Decode(&doc)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, fmt.Errorf("user not found: %s", id)
		}
		return nil, err
	}

	user := &models.User{
		ID:       doc["id"].(string),
		Username: doc["username"].(string),
		IsOnline: doc["is_online"].(bool),
	}
	if t, ok := doc["joined_at"].(primitive.DateTime); ok {
		user.JoinedAt = t.Time()
	}

	return user, nil
}

func (s *MongoDBStorage) ListUsers(ctx context.Context) ([]*models.User, error) {
	cursor, err := s.database.Collection("users").Find(ctx, bson.M{},
		options.Find().SetSort(bson.D{{Key: "joined_at", Value: 1}}),
	)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var users []*models.User
	for cursor.Next(ctx) {
		var doc bson.M
		if err := cursor.Decode(&doc); err != nil {
			return nil, err
		}

		user := &models.User{
			ID:       doc["id"].(string),
			Username: doc["username"].(string),
			IsOnline: doc["is_online"].(bool),
		}
		if t, ok := doc["joined_at"].(primitive.DateTime); ok {
			user.JoinedAt = t.Time()
		}
		users = append(users, user)
	}

	return users, cursor.Err()
}

func (s *MongoDBStorage) SetUserOnline(ctx context.Context, id string, online bool) error {
	_, err := s.database.Collection("users").UpdateOne(ctx,
		bson.M{"id": id},
		bson.M{"$set": bson.M{"is_online": online}},
	)
	return err
}

var _ Storage = (*MongoDBStorage)(nil)
