package models

import "time"

type Chat struct {
	ID        int64     `json:"id"`
	Title     string    `json:"title"`
	CreatedAt time.Time `json:"created_at"`
}

type Message struct {
	ID        int64     `json:"id"`
	ChatID    int64     `json:"chat_id"`
	Text      string    `json:"text"`
	CreatedAt time.Time `json:"created_at"`
}

type ChatWithMessages struct {
	Chat
	Messages []Message `json:"messages"`
}
