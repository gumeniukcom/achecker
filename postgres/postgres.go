package postgres

import (
	"context"
	"github.com/gumeniukcom/achecker/configs"
	"github.com/jackc/pgx/v4"
	"time"

	"github.com/jackc/pgx/v4/pgxpool"
)

// DB container fot pg pool
type DB struct {
	*pgxpool.Pool
}

// Stop is for close pg pool
func (v *DB) Stop() {
	v.Close()
}

// New create new instance pg DB
func New(cfg configs.PostgresqlConf) (*DB, error) {
	pgcfg, err := pgxpool.ParseConfig(cfg.ConnectString)
	if err != nil {
		return nil, err
	}

	pgcfg.MaxConns = 8
	pgcfg.ConnConfig.LogLevel = pgx.LogLevelDebug
	pgcfg.MaxConnLifetime = 1 * time.Minute
	pgcfg.HealthCheckPeriod = 10 * time.Second

	conn, err := pgxpool.ConnectConfig(context.Background(), pgcfg)
	if err != nil {
		return nil, err
	}

	return &DB{conn}, nil
}
