INSERT INTO messages (
    uuid,
    session_id,
    type,
    content,
    send_id,
    send_name,
    send_avatar,
    receive_id,
    status,
    created_at,
    send_at
) VALUES (
    '123e4567-e89b-12d3-a456-426614174000',  -- uuid
    'sess_001',                              -- session_id
    1,                                       -- type (1 = text message)
    'Hello, commie! This is a test message.', -- content
    'sharve',                                -- send_id
    'sharve',                                -- send_name
    'https://example.com/avatar/sharve.jpg', -- send_avatar
    'commie',                                -- receive_id
    1,                                       -- status (1 = sent)
    NOW(3),                                  -- created_at (current time with milliseconds)
    NOW(3)                                   -- send_at (same as created_at)
);