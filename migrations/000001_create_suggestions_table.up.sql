CREATE TABLE suggestions (
    id uuid PRIMARY KEY DEFAULT gen_random_uuid(),
    user_chat_id BIGINT NOT NULL,
    user_message_id BIGINT NOT NULL,
    bot_message_id BIGINT NOT NULL,
    is_answered BOOLEAN DEFAULT false
);