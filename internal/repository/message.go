package repository

import (
	"context"

	"gorm.io/gorm"

	"github.com/GlebMoskalev/chat-golang/internal/models"
)

//go:generate mockgen -destination=mocks/mock_message_repository.go -package=mocks github.com/GlebMoskalev/chat-golang/internal/repository MessageRepository

type MessageRepository interface {
	Create(ctx context.Context, message *models.Message) error
	GetByChatID(ctx context.Context, chatID int64, limit int) ([]models.Message, error)
}

type messageRepository struct {
	db *gorm.DB
}

func NewMessageRepository(db *gorm.DB) MessageRepository {
	return &messageRepository{db: db}
}

// Create создаёт сообщение
func (r *messageRepository) Create(ctx context.Context, message *models.Message) error {
	return r.db.WithContext(ctx).Create(&message).Error
}

// GetByChatID получает последние N сообщений чата
func (r *messageRepository) GetByChatID(ctx context.Context, chatID int64, limit int) ([]models.Message, error) {
	var messages []models.Message

	err := r.db.WithContext(ctx).
		Where("chat_id = ?", chatID).
		Order("created_at DESC").
		Limit(limit).
		Find(&messages).Error

	return messages, err
}
