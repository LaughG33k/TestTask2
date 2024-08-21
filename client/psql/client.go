package psql

import (
	"context"
	"time"

	"github.com/LaughG33k/TestTask2/iternal"
	"github.com/LaughG33k/TestTask2/pkg"
	"github.com/jackc/pgx"
)

func NewPool(ctx context.Context, retryAttepmts int, retryTimeSleep time.Duration, cfg iternal.DbConfig) (*pgx.ConnPool, error) {

	var pool *pgx.ConnPool

	err := pkg.Retry(func() error {

		p, err := pgx.NewConnPool(pgx.ConnPoolConfig{MaxConnections: cfg.MaxPoolSize, ConnConfig: pgx.ConnConfig{
			Host:      cfg.Host,
			Port:      cfg.Port,
			User:      cfg.User,
			Database:  cfg.DB,
			Password:  cfg.Password,
			TLSConfig: nil,
		}})

		if err != nil {
			return err
		}

		pool = p

		return nil

	}, retryAttepmts, retryTimeSleep)

	if err != nil {
		return nil, err
	}

	return pool, nil

}
