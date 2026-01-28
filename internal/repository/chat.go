package repository

import (
	"context"
	"errors"

	"gorm.io/gorm"

	"github.com/GlebMoskalev/chat-golang/internal/models"
)

//go:generate mockgen -destination=mocks/mock_chat_repository.go -package=mocks github.com/GlebMoskalev/chat-golang/internal/repository ChatRepository

var (
	ErrChatNotFound = errors.New("chat not found")
)

type ChatRepository interface {
	Create(ctx context.Context, chat *models.Chat) error
	Delete(ctx context.Context, id int64) error
	Exists(ctx context.Context, id int64) (bool, error)
	GetByID(ctx context.Context, id int64) (*models.Chat, error)
}

type chatRepository struct {
	db *gorm.DB
}

func NewChatRepository(db *gorm.DB) ChatRepository {
	return &chatRepository{db: db}
}

// Create создаёт новый чат
func (r *chatRepository) Create(ctx context.Context, chat *models.Chat) error {
	return r.db.WithContext(ctx).Create(chat).Error
}

// Delete удаляет чат (сообщения удалятся каскадом)
func (r *chatRepository) Delete(ctx context.Context, id int64) error {
	result := r.db.WithContext(ctx).Delete(&models.Chat{}, id)

	if result.Error != nil {
		return result.Error
	}

	if result.RowsAffected == 0 {
		return ErrChatNotFound
	}

	return nil
}

// GetByID получает чат по ID
func (r *chatRepository) GetByID(ctx context.Context, id int64) (*models.Chat, error) {
	var chat models.Chat
	err := r.db.WithContext(ctx).First(&chat, id).Error

	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil // Возвращаем nil, nil для "не найдено"
	}

	if err != nil {
		return nil, err
	}

	return &chat, nil
}

// Exists проверяет существование чата
func (r *chatRepository) Exists(ctx context.Context, id int64) (bool, error) {
	var count int64
	err := r.db.
		WithContext(ctx).
		Model(&models.Chat{}).
		Where("id = ?", id).
		Count(&count).Error
	return count > 0, err
}
