INSERT INTO rate_limit_usage(api_key_id, scope, window_start, window_seconds, counter)
VALUES (?, ?, ?, ?, 1)
ON DUPLICATE KEY UPDATE counter = counter + 1