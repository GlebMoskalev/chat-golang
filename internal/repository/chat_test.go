package repository

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"

	"github.com/GlebMoskalev/chat-golang/internal/models"
)

func setupTestDB(t *testing.T) *gorm.DB {
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	require.NoError(t, err)

	err = db.AutoMigrate(&models.Chat{}, &models.Message{})
	require.NoError(t, err)

	return db
}

func TestChatRepository_Create(t *testing.T) {
	db := setupTestDB(t)
	repo := NewChatRepository(db)
	ctx := context.Background()

	chat := &models.Chat{
		Title:     "Test Chat",
		CreatedAt: time.Now(),
	}

	err := repo.Create(ctx, chat)
	assert.NoError(t, err)
	assert.NotZero(t, chat.ID)
}

func TestChatRepository_GetByID(t *testing.T) {
	db := setupTestDB(t)
	repo := NewChatRepository(db)
	ctx := context.Background()

	chat := &models.Chat{
		Title:     "Test Chat",
		CreatedAt: time.Now(),
	}
	err := repo.Create(ctx, chat)
	require.NoError(t, err)

	found, err := repo.GetByID(ctx, chat.ID)
	assert.NoError(t, err)
	assert.NotNil(t, found)
	assert.Equal(t, chat.Title, found.Title)

	notFound, err := repo.GetByID(ctx, 999)
	assert.NoError(t, err)
	assert.Nil(t, notFound)
}

func TestChatRepository_Exists(t *testing.T) {
	db := setupTestDB(t)
	repo := NewChatRepository(db)
	ctx := context.Background()

	chat := &models.Chat{
		Title:     "Test Chat",
		CreatedAt: time.Now(),
	}
	err := repo.Create(ctx, chat)
	require.NoError(t, err)

	exists, err := repo.Exists(ctx, chat.ID)
	assert.NoError(t, err)
	assert.True(t, exists)

	exists, err = repo.Exists(ctx, 999)
	assert.NoError(t, err)
	assert.False(t, exists)
}

func TestChatRepository_Delete(t *testing.T) {
	db := setupTestDB(t)
	repo := NewChatRepository(db)
	ctx := context.Background()

	chat := &models.Chat{
		Title:     "Test Chat",
		CreatedAt: time.Now(),
	}
	err := repo.Create(ctx, chat)
	require.NoError(t, err)

	err = repo.Delete(ctx, chat.ID)
	assert.NoError(t, err)

	exists, err := repo.Exists(ctx, chat.ID)
	assert.NoError(t, err)
	assert.False(t, exists)

	err = repo.Delete(ctx, 999)
	assert.ErrorIs(t, err, ErrChatNotFound)
}
