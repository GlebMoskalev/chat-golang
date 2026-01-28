package handler

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/GlebMoskalev/chat-golang/internal/models"
	"github.com/GlebMoskalev/chat-golang/internal/service/mocks"
	"github.com/gorilla/mux"
	"go.uber.org/mock/gomock"
)

func TestCreateChat(t *testing.T) {
	tests := []struct {
		name           string
		requestBody    string
		setupMock      func(*mocks.MockChatServiceInterface)
		expectedStatus int
		expectedBody   string
	}{
		{
			name:        "успешное создание чата",
			requestBody: `{"title":"Тестовый чат"}`,
			setupMock: func(m *mocks.MockChatServiceInterface) {
				m.EXPECT().
					CreateChat(gomock.Any(), "Тестовый чат").
					Return(&models.Chat{
						ID:        1,
						Title:     "Тестовый чат",
						CreatedAt: time.Date(2026, 1, 28, 10, 0, 0, 0, time.UTC),
					}, nil)
			},
			expectedStatus: http.StatusCreated,
		},
		{
			name:        "пустой title",
			requestBody: `{"title":""}`,
			setupMock: func(m *mocks.MockChatServiceInterface) {
				m.EXPECT().
					CreateChat(gomock.Any(), "").
					Return(nil, errors.New("title must be 1-200 characters"))
			},
			expectedStatus: http.StatusBadRequest,
			expectedBody:   "title must be 1-200 characters",
		},
		{
			name:           "невалидный JSON",
			requestBody:    `{invalid}`,
			setupMock:      func(m *mocks.MockChatServiceInterface) {},
			expectedStatus: http.StatusBadRequest,
			expectedBody:   "Invalid JSON",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			mockService := mocks.NewMockChatServiceInterface(ctrl)
			tt.setupMock(mockService)

			handler := &ChatHandler{service: mockService}

			req := httptest.NewRequest(http.MethodPost, "/chats/", bytes.NewBufferString(tt.requestBody))
			req.Header.Set("Content-Type", "application/json")
			w := httptest.NewRecorder()

			handler.CreateChat(w, req)

			if w.Code != tt.expectedStatus {
				t.Errorf("ожидался статус %d, получен %d", tt.expectedStatus, w.Code)
			}

			if tt.expectedBody != "" && !bytes.Contains(w.Body.Bytes(), []byte(tt.expectedBody)) {
				t.Errorf("ожидалось тело ответа содержащее %q, получено %q", tt.expectedBody, w.Body.String())
			}
		})
	}
}

func TestGetChat(t *testing.T) {
	tests := []struct {
		name           string
		chatID         string
		limit          string
		setupMock      func(*mocks.MockChatServiceInterface)
		expectedStatus int
	}{
		{
			name:   "успешное получение чата",
			chatID: "1",
			limit:  "10",
			setupMock: func(m *mocks.MockChatServiceInterface) {
				m.EXPECT().
					GetChatWithMessages(gomock.Any(), int64(1), 10).
					Return(&models.ChatWithMessages{
						Chat: models.Chat{
							ID:        1,
							Title:     "Тестовый чат",
							CreatedAt: time.Date(2026, 1, 28, 10, 0, 0, 0, time.UTC),
						},
						Messages: []models.Message{
							{
								ID:        1,
								ChatID:    1,
								Text:      "Привет",
								CreatedAt: time.Date(2026, 1, 28, 10, 1, 0, 0, time.UTC),
							},
						},
					}, nil)
			},
			expectedStatus: http.StatusOK,
		},
		{
			name:   "чат не найден",
			chatID: "999",
			limit:  "20",
			setupMock: func(m *mocks.MockChatServiceInterface) {
				m.EXPECT().
					GetChatWithMessages(gomock.Any(), int64(999), 20).
					Return(nil, errors.New("chat not found"))
			},
			expectedStatus: http.StatusNotFound,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			mockService := mocks.NewMockChatServiceInterface(ctrl)
			tt.setupMock(mockService)

			handler := &ChatHandler{service: mockService}

			req := httptest.NewRequest(http.MethodGet, "/chats/"+tt.chatID+"?limit="+tt.limit, nil)
			w := httptest.NewRecorder()

			router := mux.NewRouter()
			router.HandleFunc("/chats/{id}", handler.GetChat).Methods("GET")
			router.ServeHTTP(w, req)

			if w.Code != tt.expectedStatus {
				t.Errorf("ожидался статус %d, получен %d", tt.expectedStatus, w.Code)
			}

			if tt.expectedStatus == http.StatusOK {
				var response models.ChatWithMessages
				if err := json.Unmarshal(w.Body.Bytes(), &response); err != nil {
					t.Errorf("ошибка парсинга ответа: %v", err)
				}
			}
		})
	}
}

func TestDeleteChat(t *testing.T) {
	tests := []struct {
		name           string
		chatID         string
		setupMock      func(*mocks.MockChatServiceInterface)
		expectedStatus int
	}{
		{
			name:   "успешное удаление",
			chatID: "1",
			setupMock: func(m *mocks.MockChatServiceInterface) {
				m.EXPECT().
					DeleteChat(gomock.Any(), int64(1)).
					Return(nil)
			},
			expectedStatus: http.StatusNoContent,
		},
		{
			name:   "чат не найден",
			chatID: "999",
			setupMock: func(m *mocks.MockChatServiceInterface) {
				m.EXPECT().
					DeleteChat(gomock.Any(), int64(999)).
					Return(errors.New("chat not found"))
			},
			expectedStatus: http.StatusNotFound,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			mockService := mocks.NewMockChatServiceInterface(ctrl)
			tt.setupMock(mockService)

			handler := &ChatHandler{service: mockService}

			req := httptest.NewRequest(http.MethodDelete, "/chats/"+tt.chatID, nil)
			w := httptest.NewRecorder()

			router := mux.NewRouter()
			router.HandleFunc("/chats/{id}", handler.DeleteChat).Methods("DELETE")
			router.ServeHTTP(w, req)

			if w.Code != tt.expectedStatus {
				t.Errorf("ожидался статус %d, получен %d", tt.expectedStatus, w.Code)
			}
		})
	}
}

func TestCreateMessage(t *testing.T) {
	tests := []struct {
		name           string
		chatID         string
		requestBody    string
		setupMock      func(*mocks.MockChatServiceInterface)
		expectedStatus int
		expectedBody   string
	}{
		{
			name:        "успешное создание сообщения",
			chatID:      "1",
			requestBody: `{"text":"Привет!"}`,
			setupMock: func(m *mocks.MockChatServiceInterface) {
				m.EXPECT().
					CreateMessage(gomock.Any(), int64(1), "Привет!").
					Return(&models.Message{
						ID:        1,
						ChatID:    1,
						Text:      "Привет!",
						CreatedAt: time.Date(2026, 1, 28, 10, 0, 0, 0, time.UTC),
					}, nil)
			},
			expectedStatus: http.StatusCreated,
		},
		{
			name:        "чат не найден",
			chatID:      "999",
			requestBody: `{"text":"Привет!"}`,
			setupMock: func(m *mocks.MockChatServiceInterface) {
				m.EXPECT().
					CreateMessage(gomock.Any(), int64(999), "Привет!").
					Return(nil, errors.New("chat not found"))
			},
			expectedStatus: http.StatusNotFound,
			expectedBody:   "chat not found",
		},
		{
			name:        "пустой текст",
			chatID:      "1",
			requestBody: `{"text":""}`,
			setupMock: func(m *mocks.MockChatServiceInterface) {
				m.EXPECT().
					CreateMessage(gomock.Any(), int64(1), "").
					Return(nil, errors.New("text must be 1-5000 characters"))
			},
			expectedStatus: http.StatusBadRequest,
			expectedBody:   "text must be 1-5000 characters",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			mockService := mocks.NewMockChatServiceInterface(ctrl)
			tt.setupMock(mockService)

			handler := &ChatHandler{service: mockService}

			req := httptest.NewRequest(http.MethodPost, "/chats/"+tt.chatID+"/messages/", bytes.NewBufferString(tt.requestBody))
			req.Header.Set("Content-Type", "application/json")
			w := httptest.NewRecorder()

			router := mux.NewRouter()
			router.HandleFunc("/chats/{id}/messages/", handler.CreateMessage).Methods("POST")
			router.ServeHTTP(w, req)

			if w.Code != tt.expectedStatus {
				t.Errorf("ожидался статус %d, получен %d", tt.expectedStatus, w.Code)
			}

			if tt.expectedBody != "" && !bytes.Contains(w.Body.Bytes(), []byte(tt.expectedBody)) {
				t.Errorf("ожидалось тело ответа содержащее %q, получено %q", tt.expectedBody, w.Body.String())
			}
		})
	}
}
