-- name: GetSender :one
SELECT user_chat_id, user_message_id FROM suggestions
WHERE bot_message_id = $1;

-- name: CreateSuggestion :exec
INSERT INTO suggestions (user_chat_id, user_message_id, bot_message_id)
VALUES ($1, $2, $3);
