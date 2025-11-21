SELECT counter
FROM rate_limit_usage
WHERE api_key_id = ? AND scope = ? AND window_start = ?
FOR UPDATE