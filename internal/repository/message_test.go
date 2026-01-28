package repository

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/GlebMoskalev/chat-golang/internal/models"
)

func TestMessageRepository_Create(t *testing.T) {
	db := setupTestDB(t)
	chatRepo := NewChatRepository(db)
	msgRepo := NewMessageRepository(db)
	ctx := context.Background()

	chat := &models.Chat{
		Title:     "Test Chat",
		CreatedAt: time.Now(),
	}
	err := chatRepo.Create(ctx, chat)
	require.NoError(t, err)

	message := &models.Message{
		ChatID:    chat.ID,
		Text:      "Test message",
		CreatedAt: time.Now(),
	}

	err = msgRepo.Create(ctx, message)
	assert.NoError(t, err)
	assert.NotZero(t, message.ID)
}

func TestMessageRepository_GetByChatID(t *testing.T) {
	db := setupTestDB(t)
	chatRepo := NewChatRepository(db)
	msgRepo := NewMessageRepository(db)
	ctx := context.Background()

	chat := &models.Chat{
		Title:     "Test Chat",
		CreatedAt: time.Now(),
	}
	err := chatRepo.Create(ctx, chat)
	require.NoError(t, err)

	messages := []models.Message{
		{ChatID: chat.ID, Text: "Message 1", CreatedAt: time.Now()},
		{ChatID: chat.ID, Text: "Message 2", CreatedAt: time.Now().Add(1 * time.Second)},
		{ChatID: chat.ID, Text: "Message 3", CreatedAt: time.Now().Add(2 * time.Second)},
	}

	for i := range messages {
		err := msgRepo.Create(ctx, &messages[i])
		require.NoError(t, err)
	}

	found, err := msgRepo.GetByChatID(ctx, chat.ID, 10)
	assert.NoError(t, err)
	assert.Len(t, found, 3)

	assert.Equal(t, "Message 3", found[0].Text)
	assert.Equal(t, "Message 2", found[1].Text)
	assert.Equal(t, "Message 1", found[2].Text)
}

func TestMessageRepository_GetByChatID_WithLimit(t *testing.T) {
	db := setupTestDB(t)
	chatRepo := NewChatRepository(db)
	msgRepo := NewMessageRepository(db)
	ctx := context.Background()

	chat := &models.Chat{
		Title:     "Test Chat",
		CreatedAt: time.Now(),
	}
	err := chatRepo.Create(ctx, chat)
	require.NoError(t, err)

	for i := 1; i <= 5; i++ {
		msg := &models.Message{
			ChatID:    chat.ID,
			Text:      "Message",
			CreatedAt: time.Now().Add(time.Duration(i) * time.Second),
		}
		err := msgRepo.Create(ctx, msg)
		require.NoError(t, err)
	}

	found, err := msgRepo.GetByChatID(ctx, chat.ID, 2)
	assert.NoError(t, err)
	assert.Len(t, found, 2)
}

func TestMessageRepository_GetByChatID_EmptyChat(t *testing.T) {
	db := setupTestDB(t)
	chatRepo := NewChatRepository(db)
	msgRepo := NewMessageRepository(db)
	ctx := context.Background()

	chat := &models.Chat{
		Title:     "Empty Chat",
		CreatedAt: time.Now(),
	}
	err := chatRepo.Create(ctx, chat)
	require.NoError(t, err)

	found, err := msgRepo.GetByChatID(ctx, chat.ID, 10)
	assert.NoError(t, err)
	assert.Empty(t, found)
}
