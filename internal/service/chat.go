package service

import (
	"context"
	"errors"
	"strings"

	"github.com/GlebMoskalev/chat-golang/internal/models"
	"github.com/GlebMoskalev/chat-golang/internal/repository"
)

//go:generate mockgen -destination=mocks/mock_chat_service.go -package=mocks github.com/GlebMoskalev/chat-golang/internal/service ChatServiceInterface

type ChatServiceInterface interface {
	CreateChat(ctx context.Context, title string) (*models.Chat, error)
	GetChatWithMessages(ctx context.Context, chatID int64, limit int) (*models.ChatWithMessages, error)
	DeleteChat(ctx context.Context, chatID int64) error
	CreateMessage(ctx context.Context, chatID int64, text string) (*models.Message, error)
}

type ChatService struct {
	chatRepo    repository.ChatRepository
	messageRepo repository.MessageRepository
}

func NewChatService(chatRepo repository.ChatRepository, messageRepo repository.MessageRepository) *ChatService {
	return &ChatService{
		chatRepo:    chatRepo,
		messageRepo: messageRepo,
	}
}

// CreateChat создаёт новый чат
func (s *ChatService) CreateChat(ctx context.Context, title string) (*models.Chat, error) {
	title = strings.TrimSpace(title)
	if title == "" {
		return nil, errors.New("title cannot be empty")
	}
	if len(title) > 200 {
		return nil, errors.New("title must be 1-200 characters")
	}

	chat := &models.Chat{
		Title: title,
	}

	if err := s.chatRepo.Create(ctx, chat); err != nil {
		return nil, err
	}

	return chat, nil
}

// GetChatWithMessages получает чат с сообщениями
func (s *ChatService) GetChatWithMessages(ctx context.Context, chatID int64, limit int) (*models.ChatWithMessages, error) {
	if limit <= 0 {
		limit = 20
	}
	if limit > 100 {
		limit = 100
	}

	chat, err := s.chatRepo.GetByID(ctx, chatID)
	if err != nil {
		return nil, err
	}
	if chat == nil {
		return nil, errors.New("chat not found")
	}

	messages, err := s.messageRepo.GetByChatID(ctx, chatID, limit)
	if err != nil {
		return nil, err
	}

	return &models.ChatWithMessages{
		Chat:     *chat,
		Messages: messages,
	}, nil
}

// DeleteChat удаляет чат
func (s *ChatService) DeleteChat(ctx context.Context, chatID int64) error {
	return s.chatRepo.Delete(ctx, chatID)
}

// CreateMessage создаёт сообщение
func (s *ChatService) CreateMessage(ctx context.Context, chatID int64, text string) (*models.Message, error) {
	exists, err := s.chatRepo.Exists(ctx, chatID)
	if err != nil {
		return nil, err
	}
	if !exists {
		return nil, errors.New("chat not found")
	}

	text = strings.TrimSpace(text)
	if text == "" {
		return nil, errors.New("text cannot be empty")
	}
	if len(text) > 5000 {
		return nil, errors.New("text must be 1-5000 characters")
	}

	message := &models.Message{
		ChatID: chatID,
		Text:   text,
	}

	if err := s.messageRepo.Create(ctx, message); err != nil {
		return nil, err
	}

	return message, nil
}
