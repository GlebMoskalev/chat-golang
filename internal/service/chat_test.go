package service

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/GlebMoskalev/chat-golang/internal/models"
	"github.com/GlebMoskalev/chat-golang/internal/repository/mocks"
	"go.uber.org/mock/gomock"
)

func TestCreateChat(t *testing.T) {
	tests := []struct {
		name        string
		title       string
		setupMock   func(*mocks.MockChatRepository)
		expectError bool
		errorMsg    string
	}{
		{
			name:  "успешное создание чата",
			title: "Тестовый чат",
			setupMock: func(m *mocks.MockChatRepository) {
				m.EXPECT().
					Create(gomock.Any(), gomock.Any()).
					DoAndReturn(func(ctx context.Context, chat *models.Chat) error {
						chat.ID = 1
						chat.CreatedAt = time.Now()
						return nil
					})
			},
			expectError: false,
		},
		{
			name:  "тримминг пробелов",
			title: "  Чат с пробелами  ",
			setupMock: func(m *mocks.MockChatRepository) {
				m.EXPECT().
					Create(gomock.Any(), gomock.Any()).
					DoAndReturn(func(ctx context.Context, chat *models.Chat) error {
						if chat.Title != "Чат с пробелами" {
							t.Errorf("ожидался title 'Чат с пробелами', получен '%s'", chat.Title)
						}
						chat.ID = 1
						chat.CreatedAt = time.Now()
						return nil
					})
			},
			expectError: false,
		},
		{
			name:        "пустой title",
			title:       "",
			setupMock:   func(m *mocks.MockChatRepository) {},
			expectError: true,
			errorMsg:    "title cannot be empty",
		},
		{
			name:        "title только из пробелов",
			title:       "   ",
			setupMock:   func(m *mocks.MockChatRepository) {},
			expectError: true,
			errorMsg:    "title cannot be empty",
		},
		{
			name:        "слишком длинный title",
			title:       string(make([]byte, 201)),
			setupMock:   func(m *mocks.MockChatRepository) {},
			expectError: true,
			errorMsg:    "title must be 1-200 characters",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			mockChatRepo := mocks.NewMockChatRepository(ctrl)
			mockMessageRepo := mocks.NewMockMessageRepository(ctrl)

			tt.setupMock(mockChatRepo)

			service := NewChatService(mockChatRepo, mockMessageRepo)

			chat, err := service.CreateChat(context.Background(), tt.title)

			if tt.expectError {
				if err == nil {
					t.Error("ожидалась ошибка, но её не было")
				}
				if tt.errorMsg != "" && err.Error() != tt.errorMsg {
					t.Errorf("ожидалась ошибка %q, получена %q", tt.errorMsg, err.Error())
				}
			} else {
				if err != nil {
					t.Errorf("неожиданная ошибка: %v", err)
				}
				if chat == nil {
					t.Error("чат не должен быть nil")
				}
			}
		})
	}
}

func TestGetChatWithMessages(t *testing.T) {
	tests := []struct {
		name          string
		chatID        int64
		limit         int
		setupMock     func(*mocks.MockChatRepository, *mocks.MockMessageRepository)
		expectError   bool
		expectedLimit int
	}{
		{
			name:   "успешное получение чата с сообщениями",
			chatID: 1,
			limit:  10,
			setupMock: func(cr *mocks.MockChatRepository, mr *mocks.MockMessageRepository) {
				cr.EXPECT().
					GetByID(gomock.Any(), int64(1)).
					Return(&models.Chat{
						ID:        1,
						Title:     "Тест",
						CreatedAt: time.Now(),
					}, nil)

				mr.EXPECT().
					GetByChatID(gomock.Any(), int64(1), 10).
					Return([]models.Message{
						{ID: 1, ChatID: 1, Text: "Привет", CreatedAt: time.Now()},
					}, nil)
			},
			expectError:   false,
			expectedLimit: 10,
		},
		{
			name:   "чат не найден",
			chatID: 999,
			limit:  20,
			setupMock: func(cr *mocks.MockChatRepository, mr *mocks.MockMessageRepository) {
				cr.EXPECT().
					GetByID(gomock.Any(), int64(999)).
					Return(nil, nil)
			},
			expectError: true,
		},
		{
			name:   "limit по умолчанию",
			chatID: 1,
			limit:  0,
			setupMock: func(cr *mocks.MockChatRepository, mr *mocks.MockMessageRepository) {
				cr.EXPECT().
					GetByID(gomock.Any(), int64(1)).
					Return(&models.Chat{
						ID:        1,
						Title:     "Тест",
						CreatedAt: time.Now(),
					}, nil)

				mr.EXPECT().
					GetByChatID(gomock.Any(), int64(1), 20).
					Return([]models.Message{}, nil)
			},
			expectError:   false,
			expectedLimit: 20,
		},
		{
			name:   "limit больше максимума",
			chatID: 1,
			limit:  150,
			setupMock: func(cr *mocks.MockChatRepository, mr *mocks.MockMessageRepository) {
				cr.EXPECT().
					GetByID(gomock.Any(), int64(1)).
					Return(&models.Chat{
						ID:        1,
						Title:     "Тест",
						CreatedAt: time.Now(),
					}, nil)

				mr.EXPECT().
					GetByChatID(gomock.Any(), int64(1), 100).
					Return([]models.Message{}, nil)
			},
			expectError:   false,
			expectedLimit: 100,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			mockChatRepo := mocks.NewMockChatRepository(ctrl)
			mockMessageRepo := mocks.NewMockMessageRepository(ctrl)

			tt.setupMock(mockChatRepo, mockMessageRepo)

			service := NewChatService(mockChatRepo, mockMessageRepo)

			result, err := service.GetChatWithMessages(context.Background(), tt.chatID, tt.limit)

			if tt.expectError {
				if err == nil {
					t.Error("ожидалась ошибка, но её не было")
				}
			} else {
				if err != nil {
					t.Errorf("неожиданная ошибка: %v", err)
				}
				if result == nil {
					t.Error("результат не должен быть nil")
				}
			}
		})
	}
}

func TestCreateMessage(t *testing.T) {
	tests := []struct {
		name        string
		chatID      int64
		text        string
		setupMock   func(*mocks.MockChatRepository, *mocks.MockMessageRepository)
		expectError bool
		errorMsg    string
	}{
		{
			name:   "успешное создание сообщения",
			chatID: 1,
			text:   "Привет!",
			setupMock: func(cr *mocks.MockChatRepository, mr *mocks.MockMessageRepository) {
				cr.EXPECT().
					Exists(gomock.Any(), int64(1)).
					Return(true, nil)

				mr.EXPECT().
					Create(gomock.Any(), gomock.Any()).
					DoAndReturn(func(ctx context.Context, message *models.Message) error {
						message.ID = 1
						message.CreatedAt = time.Now()
						return nil
					})
			},
			expectError: false,
		},
		{
			name:   "тримминг пробелов",
			chatID: 1,
			text:   "  Сообщение с пробелами  ",
			setupMock: func(cr *mocks.MockChatRepository, mr *mocks.MockMessageRepository) {
				cr.EXPECT().
					Exists(gomock.Any(), int64(1)).
					Return(true, nil)

				mr.EXPECT().
					Create(gomock.Any(), gomock.Any()).
					DoAndReturn(func(ctx context.Context, message *models.Message) error {
						if message.Text != "Сообщение с пробелами" {
							t.Errorf("ожидался text 'Сообщение с пробелами', получен '%s'", message.Text)
						}
						message.ID = 1
						message.CreatedAt = time.Now()
						return nil
					})
			},
			expectError: false,
		},
		{
			name:   "чат не существует",
			chatID: 999,
			text:   "Привет!",
			setupMock: func(cr *mocks.MockChatRepository, mr *mocks.MockMessageRepository) {
				cr.EXPECT().
					Exists(gomock.Any(), int64(999)).
					Return(false, nil)
			},
			expectError: true,
			errorMsg:    "chat not found",
		},
		{
			name:   "пустой текст",
			chatID: 1,
			text:   "",
			setupMock: func(cr *mocks.MockChatRepository, mr *mocks.MockMessageRepository) {
				cr.EXPECT().
					Exists(gomock.Any(), int64(1)).
					Return(true, nil)
			},
			expectError: true,
			errorMsg:    "text cannot be empty",
		},
		{
			name:   "текст только из пробелов",
			chatID: 1,
			text:   "   ",
			setupMock: func(cr *mocks.MockChatRepository, mr *mocks.MockMessageRepository) {
				cr.EXPECT().
					Exists(gomock.Any(), int64(1)).
					Return(true, nil)
			},
			expectError: true,
			errorMsg:    "text cannot be empty",
		},
		{
			name:   "слишком длинный текст",
			chatID: 1,
			text:   string(make([]byte, 5001)),
			setupMock: func(cr *mocks.MockChatRepository, mr *mocks.MockMessageRepository) {
				cr.EXPECT().
					Exists(gomock.Any(), int64(1)).
					Return(true, nil)
			},
			expectError: true,
			errorMsg:    "text must be 1-5000 characters",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			mockChatRepo := mocks.NewMockChatRepository(ctrl)
			mockMessageRepo := mocks.NewMockMessageRepository(ctrl)

			tt.setupMock(mockChatRepo, mockMessageRepo)

			service := NewChatService(mockChatRepo, mockMessageRepo)

			message, err := service.CreateMessage(context.Background(), tt.chatID, tt.text)

			if tt.expectError {
				if err == nil {
					t.Error("ожидалась ошибка, но её не было")
				}
				if tt.errorMsg != "" && err.Error() != tt.errorMsg {
					t.Errorf("ожидалась ошибка %q, получена %q", tt.errorMsg, err.Error())
				}
			} else {
				if err != nil {
					t.Errorf("неожиданная ошибка: %v", err)
				}
				if message == nil {
					t.Error("сообщение не должно быть nil")
				}
			}
		})
	}
}

func TestDeleteChat(t *testing.T) {
	tests := []struct {
		name        string
		chatID      int64
		setupMock   func(*mocks.MockChatRepository)
		expectError bool
	}{
		{
			name:   "успешное удаление",
			chatID: 1,
			setupMock: func(m *mocks.MockChatRepository) {
				m.EXPECT().
					Delete(gomock.Any(), int64(1)).
					Return(nil)
			},
			expectError: false,
		},
		{
			name:   "ошибка при удалении",
			chatID: 999,
			setupMock: func(m *mocks.MockChatRepository) {
				m.EXPECT().
					Delete(gomock.Any(), int64(999)).
					Return(errors.New("chat not found"))
			},
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			mockChatRepo := mocks.NewMockChatRepository(ctrl)
			mockMessageRepo := mocks.NewMockMessageRepository(ctrl)

			tt.setupMock(mockChatRepo)

			service := NewChatService(mockChatRepo, mockMessageRepo)

			err := service.DeleteChat(context.Background(), tt.chatID)

			if tt.expectError {
				if err == nil {
					t.Error("ожидалась ошибка, но её не было")
				}
			} else {
				if err != nil {
					t.Errorf("неожиданная ошибка: %v", err)
				}
			}
		})
	}
}
