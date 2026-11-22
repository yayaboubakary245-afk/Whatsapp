package security

import (
	"context"
	"encoding/json"
	"sync"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)

// RateLimiter protège l'API REST contre les abus (Ndue).
type RateLimiter struct {
	mu      sync.Mutex
	buckets map[string][]time.Time
	limit   int
	window  time.Duration
}

func NewRateLimiter(limit int, window time.Duration) *RateLimiter {
	return &RateLimiter{
		buckets: make(map[string][]time.Time),
		limit:   limit,
		window:  window,
	}
}

func (rl *RateLimiter) Allow(key string) bool {
	rl.mu.Lock()
	defer rl.mu.Unlock()
	now := time.Now()
	cutoff := now.Add(-rl.window)
	var recent []time.Time
	for _, t := range rl.buckets[key] {
		if t.After(cutoff) {
			recent = append(recent, t)
		}
	}
	if len(recent) >= rl.limit {
		rl.buckets[key] = recent
		return false
	}
	rl.buckets[key] = append(recent, now)
	return true
}

func LogAudit(ctx context.Context, db *pgxpool.Pool, actorID, action, targetType, targetID, ip string, metadata map[string]any) {
	metaJSON, _ := json.Marshal(metadata)
	_, _ = db.Exec(ctx, `
		INSERT INTO audit_logs (actor_id, action, target_type, target_id, metadata, ip_address)
		VALUES (NULLIF($1, '')::uuid, $2, $3, $4, $5::jsonb, $6)
	`, actorID, action, targetType, targetID, string(metaJSON), ip)
}
