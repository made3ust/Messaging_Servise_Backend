CREATE TABLE messages (
                          id SERIAL PRIMARY KEY,
                          created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
                          updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
                          deleted_at TIMESTAMP WITH TIME ZONE,
                          content TEXT,
                          user_id INTEGER REFERENCES users(id),
                          chat_id INTEGER REFERENCES chats(id)
);
CREATE INDEX idx_messages_deleted_at ON messages(deleted_at);