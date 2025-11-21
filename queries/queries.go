package queries

import _ "embed"

//go:embed GetConfigLimit.sql
var GetConfigLimit string

//go:embed UpsertRateLimitUsage.sql
var UpsertRateLimitUsage string

//go:embed GetCounterRateLimitUsage.sql
var GetCounterRateLimitUsage string
