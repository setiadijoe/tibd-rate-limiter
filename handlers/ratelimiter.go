package handlers

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/setiadijoe/tibd-rate-limiter/ratelimit"
)

type Handler struct {
	Service *ratelimit.Service
}

func New(rl *ratelimit.Service) *Handler {
	return &Handler{
		Service: rl,
	}
}

func (h *Handler) RateLimitMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		apiKey := r.Header.Get("X-API-Key")
		if apiKey == "" {
			http.Error(w, "missing API Key", http.StatusUnauthorized)
			return
		}

		scope := r.Method + ":" + r.URL.Path

		res, err := h.Service.CheckAndConsume(context.Background(), apiKey, scope)
		if err != nil {
			log.Println("rate limit error:", err)
			http.Error(w, "internal error", http.StatusInternalServerError)
			return
		}

		// Set header info limit
		w.Header().Set("X-RateLimit-Limit", fmt.Sprint(res.Limit))
		w.Header().Set("X-RateLimit-Remaining", fmt.Sprint(res.Limit-res.CurrentCount))
		w.Header().Set("X-RateLimit-Reset", res.ResetAt.Format(time.RFC3339))

		if !res.Allowed {
			w.WriteHeader(http.StatusTooManyRequests)
			_ = json.NewEncoder(w).Encode(map[string]any{
				"error":   "rate_limit_exceeded",
				"message": "Too many requests",
			})
			return
		}

		next.ServeHTTP(w, r)
	})
}
