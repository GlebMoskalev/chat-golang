-- +goose Up
-- +goose StatementBegin
CREATE TABLE chats (
    id BIGSERIAL PRIMARY KEY,
    title VARCHAR(200) NOT NULL CHECK (length(title) > 0),
    created_at TIMESTAMP NOT NULL DEFAULT now()
);

CREATE TABLE messages (
    id BIGSERIAL PRIMARY KEY,
    chat_id BIGINT NOT NULL,
    text VARCHAR(5000) NOT NULL CHECK (length(text)  > 0),
    created_at TIMESTAMP NOT NULL DEFAULT now(),
    CONSTRAINT fk_messages_chat
                      FOREIGN KEY (chat_id)
                      REFERENCES chats(id)
                      ON DELETE CASCADE
);

CREATE INDEX idx_messages_chat_id ON messages(chat_id);
CREATE INDEX idx_messages_chat_created ON messages(chat_id, created_at DESC);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS messages;
DROP TABLE IF EXISTS chats;
-- +goose StatementEnd
