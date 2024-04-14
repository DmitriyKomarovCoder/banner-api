package postgres

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v4/pgxpool"
)

type Postgres struct {
	Pool *pgxpool.Pool
}

func New(connStr string, poolSize int32) (*Postgres, error) {
	var pg Postgres

	poolConfig, err := pgxpool.ParseConfig(connStr)
	if err != nil {
		return nil, fmt.Errorf("postgres - NewPostgres - pgxpool.ParseConfig: %w", err)
	}

	poolConfig.MaxConns = poolSize

	pg.Pool, err = pgxpool.ConnectConfig(context.Background(), poolConfig)

	if err != nil {
		return nil, fmt.Errorf("postgres - NewPostgres - connAttempts == 0: %w", err)
	}

	return &pg, nil
}

func (p *Postgres) Close(ctx context.Context) error {
	if p.Pool != nil {
		p.Pool.Close()
	}
	return nil
}
