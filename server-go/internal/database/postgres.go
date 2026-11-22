package database

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)

// Connect initialise le pool PostgreSQL exclusivement depuis le serveur Go (Yaya / Mengue).
func Connect(ctx context.Context, databaseURL string) (*pgxpool.Pool, error) {
	cfg, err := pgxpool.ParseConfig(databaseURL)
	if err != nil {
		return nil, fmt.Errorf("invalid database url: %w", err)
	}
	cfg.MaxConns = 25
	cfg.MinConns = 2
	cfg.MaxConnLifetime = 30 * time.Minute

	pool, err := pgxpool.NewWithConfig(ctx, cfg)
	if err != nil {
		return nil, fmt.Errorf("pgxpool init failed: %w", err)
	}
	if err := pool.Ping(ctx); err != nil {
		pool.Close()
		return nil, fmt.Errorf("postgres ping failed: %w", err)
	}
	return pool, nil
}

// RunMigrations exécute les migrations SQL S1 à S6 dans l'ordre si le dossier existe.
func RunMigrations(ctx context.Context, pool *pgxpool.Pool, migrationsDir string) error {
	entries, err := os.ReadDir(migrationsDir)
	if err != nil {
		return fmt.Errorf("cannot read migrations dir %s: %w", migrationsDir, err)
	}
	var files []string
	for _, e := range entries {
		if !e.IsDir() && filepath.Ext(e.Name()) == ".sql" {
			files = append(files, e.Name())
		}
	}
	sort.Strings(files)

	for _, f := range files {
		sqlBytes, err := os.ReadFile(filepath.Join(migrationsDir, f))
		if err != nil {
			return fmt.Errorf("read migration %s: %w", f, err)
		}
		if _, err := pool.Exec(ctx, string(sqlBytes)); err != nil {
			return fmt.Errorf("exec migration %s: %w", f, err)
		}
	}
	return nil
}
