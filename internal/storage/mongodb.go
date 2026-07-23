package storage

import (
	"context"
	"fmt"
	"time"

	"github.com/voip-app/pkg/models"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

const (
	dbName          = "voip"
	channelsColl    = "channels"
	messagesColl    = "messages"
	usersColl       = "users"
	connectTimeout  = 10 * time.Second
	disconnectTimeout = 5 * time.Second
)

type channelDoc struct {
	ID        string             `bson:"id"`
	Name      string             `bson:"name"`
	Type      string             `bson:"type"`
	IsDefault bool               `bson:"is_default"`
	CreatedAt primitive.DateTime `bson:"created_at"`
	Deleted   bool               `bson:"deleted"`
}

type messageDoc struct {
	ID        string             `bson:"id"`
	ChannelID string             `bson:"channel_id"`
	UserID    string             `bson:"user_id"`
	Username  string             `bson:"username"`
	Content   string             `bson:"content"`
	CreatedAt primitive.DateTime `bson:"created_at"`
}

type userDoc struct {
	ID       string             `bson:"id"`
	Username string             `bson:"username"`
	JoinedAt primitive.DateTime `bson:"joined_at"`
	IsOnline bool               `bson:"is_online"`
}

type MongoDBStorage struct {
	client   *mongo.Client
	database *mongo.Database
}

func NewMongoDBStorage(uri string) (*MongoDBStorage, error) {
	ctx, cancel := context.WithTimeout(context.Background(), connectTimeout)
	defer cancel()

	client, err := mongo.Connect(ctx, options.Client().ApplyURI(uri))
	if err != nil {
		return nil, fmt.Errorf("connect mongo: %w", err)
	}

	if err := client.Ping(ctx, nil); err != nil {
		return nil, fmt.Errorf("ping mongo: %w", err)
	}

	db := client.Database(dbName)

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

	channelsCol := db.Collection(channelsColl)
	_, err := channelsCol.Indexes().CreateOne(ctx, mongo.IndexModel{
		Keys:    bson.D{{Key: "name", Value: 1}},
		Options: options.Index().SetUnique(true),
	})
	if err != nil {
		return fmt.Errorf("channels name index: %w", err)
	}

	messagesCol := db.Collection(messagesColl)
	_, err = messagesCol.Indexes().CreateOne(ctx, mongo.IndexModel{
		Keys: bson.D{{Key: "channel_id", Value: 1}, {Key: "created_at", Value: -1}},
	})
	if err != nil {
		return fmt.Errorf("messages index: %w", err)
	}

	usersCol := db.Collection(usersColl)
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
	ctx, cancel := context.WithTimeout(context.Background(), disconnectTimeout)
	defer cancel()
	return s.client.Disconnect(ctx)
}

func (s *MongoDBStorage) CreateChannel(ctx context.Context, ch *models.Channel) error {
	doc := channelDoc{
		ID:        ch.ID,
		Name:      ch.Name,
		Type:      string(ch.Type),
		IsDefault: ch.IsDefault,
		CreatedAt: primitive.NewDateTimeFromTime(ch.CreatedAt),
		Deleted:   ch.Deleted,
	}
	_, err := s.database.Collection(channelsColl).InsertOne(ctx, doc)
	if err != nil {
		return fmt.Errorf("insert channel: %w", err)
	}
	return nil
}

func (s *MongoDBStorage) GetChannel(ctx context.Context, id string) (*models.Channel, error) {
	var doc channelDoc
	err := s.database.Collection(channelsColl).FindOne(ctx, bson.M{
		"id":      id,
		"deleted": false,
	}).Decode(&doc)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, fmt.Errorf("channel not found: %s", id)
		}
		return nil, fmt.Errorf("find channel: %w", err)
	}

	return &models.Channel{
		ID:        doc.ID,
		Name:      doc.Name,
		Type:      models.ChannelType(doc.Type),
		IsDefault: doc.IsDefault,
		CreatedAt: doc.CreatedAt.Time(),
		Deleted:   doc.Deleted,
	}, nil
}

func (s *MongoDBStorage) ListChannels(ctx context.Context) ([]*models.Channel, error) {
	cursor, err := s.database.Collection(channelsColl).Find(ctx, bson.M{"deleted": false},
		options.Find().SetSort(bson.D{{Key: "created_at", Value: 1}}),
	)
	if err != nil {
		return nil, fmt.Errorf("find channels: %w", err)
	}
	defer cursor.Close(ctx)

	var channels []*models.Channel
	for cursor.Next(ctx) {
		var doc channelDoc
		if err := cursor.Decode(&doc); err != nil {
			return nil, fmt.Errorf("decode channel: %w", err)
		}

		channels = append(channels, &models.Channel{
			ID:        doc.ID,
			Name:      doc.Name,
			Type:      models.ChannelType(doc.Type),
			IsDefault: doc.IsDefault,
			CreatedAt: doc.CreatedAt.Time(),
			Deleted:   doc.Deleted,
		})
	}

	return channels, cursor.Err()
}

func (s *MongoDBStorage) UpdateChannel(ctx context.Context, ch *models.Channel) error {
	_, err := s.database.Collection(channelsColl).UpdateOne(ctx,
		bson.M{"id": ch.ID},
		bson.M{"$set": bson.M{
			"name":       ch.Name,
			"type":       string(ch.Type),
			"is_default": ch.IsDefault,
			"deleted":    ch.Deleted,
		}},
	)
	if err != nil {
		return fmt.Errorf("update channel: %w", err)
	}
	return nil
}

func (s *MongoDBStorage) DeleteChannel(ctx context.Context, id string) error {
	_, err := s.database.Collection(channelsColl).UpdateOne(ctx,
		bson.M{"id": id},
		bson.M{"$set": bson.M{"deleted": true}},
	)
	if err != nil {
		return fmt.Errorf("delete channel: %w", err)
	}
	return nil
}

func (s *MongoDBStorage) SendMessage(ctx context.Context, msg *models.Message) error {
	doc := messageDoc{
		ID:        msg.ID,
		ChannelID: msg.ChannelID,
		UserID:    msg.UserID,
		Username:  msg.Username,
		Content:   msg.Content,
		CreatedAt: primitive.NewDateTimeFromTime(msg.CreatedAt),
	}
	_, err := s.database.Collection(messagesColl).InsertOne(ctx, doc)
	if err != nil {
		return fmt.Errorf("insert message: %w", err)
	}
	return nil
}

func (s *MongoDBStorage) ListMessages(ctx context.Context, channelID string, limit, offset int) ([]*models.Message, error) {
	if limit <= 0 {
		limit = 50
	}

	opts := options.Find().
		SetSort(bson.D{{Key: "created_at", Value: -1}}).
		SetLimit(int64(limit)).
		SetSkip(int64(offset))

	cursor, err := s.database.Collection(messagesColl).Find(ctx, bson.M{"channel_id": channelID}, opts)
	if err != nil {
		return nil, fmt.Errorf("find messages: %w", err)
	}
	defer cursor.Close(ctx)

	var messages []*models.Message
	for cursor.Next(ctx) {
		var doc messageDoc
		if err := cursor.Decode(&doc); err != nil {
			return nil, fmt.Errorf("decode message: %w", err)
		}

		messages = append(messages, &models.Message{
			ID:        doc.ID,
			ChannelID: doc.ChannelID,
			UserID:    doc.UserID,
			Username:  doc.Username,
			Content:   doc.Content,
			CreatedAt: doc.CreatedAt.Time(),
		})
	}

	return messages, cursor.Err()
}

func (s *MongoDBStorage) DeleteMessage(ctx context.Context, id string) error {
	_, err := s.database.Collection(messagesColl).DeleteOne(ctx, bson.M{"id": id})
	if err != nil {
		return fmt.Errorf("delete message: %w", err)
	}
	return nil
}

func (s *MongoDBStorage) AddUser(ctx context.Context, user *models.User) error {
	doc := userDoc{
		ID:       user.ID,
		Username: user.Username,
		JoinedAt: primitive.NewDateTimeFromTime(user.JoinedAt),
		IsOnline: user.IsOnline,
	}
	_, err := s.database.Collection(usersColl).InsertOne(ctx, doc)
	if err != nil {
		return fmt.Errorf("insert user: %w", err)
	}
	return nil
}

func (s *MongoDBStorage) GetUser(ctx context.Context, id string) (*models.User, error) {
	var doc userDoc
	err := s.database.Collection(usersColl).FindOne(ctx, bson.M{"id": id}).Decode(&doc)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, fmt.Errorf("user not found: %s", id)
		}
		return nil, fmt.Errorf("find user: %w", err)
	}

	return &models.User{
		ID:       doc.ID,
		Username: doc.Username,
		JoinedAt: doc.JoinedAt.Time(),
		IsOnline: doc.IsOnline,
	}, nil
}

func (s *MongoDBStorage) ListUsers(ctx context.Context) ([]*models.User, error) {
	cursor, err := s.database.Collection(usersColl).Find(ctx, bson.M{},
		options.Find().SetSort(bson.D{{Key: "joined_at", Value: 1}}),
	)
	if err != nil {
		return nil, fmt.Errorf("find users: %w", err)
	}
	defer cursor.Close(ctx)

	var users []*models.User
	for cursor.Next(ctx) {
		var doc userDoc
		if err := cursor.Decode(&doc); err != nil {
			return nil, fmt.Errorf("decode user: %w", err)
		}

		users = append(users, &models.User{
			ID:       doc.ID,
			Username: doc.Username,
			JoinedAt: doc.JoinedAt.Time(),
			IsOnline: doc.IsOnline,
		})
	}

	return users, cursor.Err()
}

func (s *MongoDBStorage) SetUserOnline(ctx context.Context, id string, online bool) error {
	_, err := s.database.Collection(usersColl).UpdateOne(ctx,
		bson.M{"id": id},
		bson.M{"$set": bson.M{"is_online": online}},
	)
	if err != nil {
		return fmt.Errorf("set user online: %w", err)
	}
	return nil
}

var _ Storage = (*MongoDBStorage)(nil)
