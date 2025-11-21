package ratelimit

import (
	"context"
	"database/sql"
	"time"

	"github.com/setiadijoe/tibd-rate-limiter/queries"
)

type Service struct {
	DB *sql.DB
}

// Result untuk info ke caller
type Result struct {
	Allowed      bool
	CurrentCount int
	Limit        int
	ResetAt      time.Time
}

func New(db *sql.DB) *Service {
	return &Service{
		DB: db,
	}
}

// CheckAndConsume: cek quota dan sekalian increment counter
func (s *Service) CheckAndConsume(
	ctx context.Context,
	apiKey, scope string,
) (*Result, error) {
	// 1. Ambil config limit (bisa di-cache, tapi untuk awal query saja)
	var limitCount, windowSeconds int
	err := s.DB.QueryRowContext(ctx, queries.GetConfigLimit, apiKey, scope).
		Scan(&limitCount, &windowSeconds)
	if err == sql.ErrNoRows {
		// fallback: misal kasih default limit global
		limitCount = 60
		windowSeconds = 60
	} else if err != nil {
		return nil, err
	}

	now := time.Now().UTC()
	// hitung window_start: floor ke bawah
	windowDur := time.Duration(windowSeconds) * time.Second
	windowStart := now.Truncate(windowDur)
	windowEnd := windowStart.Add(windowDur)

	tx, err := s.DB.BeginTx(ctx, nil)
	if err != nil {
		return nil, err
	}
	defer tx.Rollback() // aman kalau sudah Commit

	// 2. Insert / increment counter
	_, err = tx.ExecContext(ctx, queries.UpsertRateLimitUsage, apiKey, scope, windowStart, windowSeconds)
	if err != nil {
		return nil, err
	}

	// 3. Ambil nilai counter setelah update
	var counter int
	err = tx.QueryRowContext(ctx, queries.GetCounterRateLimitUsage, apiKey, scope, windowStart).Scan(&counter)
	if err != nil {
		return nil, err
	}

	if err := tx.Commit(); err != nil {
		return nil, err
	}

	allowed := counter <= limitCount

	return &Result{
		Allowed:      allowed,
		CurrentCount: counter,
		Limit:        limitCount,
		ResetAt:      windowEnd,
	}, nil
}
