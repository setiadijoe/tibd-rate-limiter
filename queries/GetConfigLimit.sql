SELECT 
    limit_count, 
    window_seconds
FROM 
    rate_limit_config
WHERE 
    api_key_id = ? 
    AND scope = ?