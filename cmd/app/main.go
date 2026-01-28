package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/gorilla/mux"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"

	"github.com/GlebMoskalev/chat-golang/internal/handler"
	"github.com/GlebMoskalev/chat-golang/internal/repository"
	"github.com/GlebMoskalev/chat-golang/internal/service"
)

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

func main() {
	host := getEnv("DB_HOST", "localhost")
	user := getEnv("DB_USER", "postgres")
	password := getEnv("DB_PASSWORD", "postgres")
	dbname := getEnv("DB_NAME", "chat")
	port := getEnv("DB_PORT", "5432")

	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=disable",
		host, user, password, dbname, port)

	log.Println("Connecting to database...")
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal("Failed to connect to database:", err)
	}

	chatRepo := repository.NewChatRepository(db)
	messageRepo := repository.NewMessageRepository(db)
	chatService := service.NewChatService(chatRepo, messageRepo)
	chatHandler := handler.NewChatHandler(chatService)

	r := mux.NewRouter()
	r.HandleFunc("/chats/", chatHandler.CreateChat).Methods("POST")
	r.HandleFunc("/chats/{id}", chatHandler.GetChat).Methods("GET")
	r.HandleFunc("/chats/{id}", chatHandler.DeleteChat).Methods("DELETE")
	r.HandleFunc("/chats/{id}/messages/", chatHandler.CreateMessage).Methods("POST")

	log.Println("Server started on :8080")
	log.Fatal(http.ListenAndServe(":8080", r))
}
