CREATE TABLE chats (
                       id SERIAL PRIMARY KEY,
                       created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
                       updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
                       deleted_at TIMESTAMP WITH TIME ZONE,
                       name VARCHAR(255),
                       is_group BOOLEAN DEFAULT FALSE
);
CREATE INDEX idx_chats_deleted_at ON chats(deleted_at);