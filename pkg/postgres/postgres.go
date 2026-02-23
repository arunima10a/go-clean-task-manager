package postgres

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
)

const (
	_defaultMaxPoolSize = 1
	_defaultConnTimeout = time.Second
	_defaultConnAttemps = 10
)

type DB interface {
	Exec(ctx context.Context, sql string, arguments ...any) (pgconn.CommandTag, error)
	Query(ctx context.Context, sql string, args ...any) (pgx.Rows, error)
	QueryRow(ctx context.Context, sql string, args ...any) pgx.Row
}

type ctxTxKey struct{}

type Postgres struct {
	maxPoolSize int
	connTimeout time.Duration
	connAttemps int

	Pool *pgxpool.Pool
}

func New(url string, opts ...Option) (*Postgres, error) {
	pg := &Postgres{
		maxPoolSize: _defaultMaxPoolSize,
		connTimeout: _defaultConnTimeout,
		connAttemps: _defaultConnAttemps,
	}

	for _, opt := range opts {
		opt(pg)
	}

	poolConfig, err := pgxpool.ParseConfig(url)
	if err != nil {
		return nil, fmt.Errorf("postgres - NewPostgres - pgxpool.ParseConfig: %w", err)
	}

	poolConfig.MaxConns = int32(pg.maxPoolSize)

	for pg.connAttemps > 0 {
		pg.Pool, err = pgxpool.NewWithConfig(context.Background(), poolConfig)
		if err == nil {
			break
		}

		log.Printf("Postgres is trying to connect, attemps left: %d", pg.connAttemps)
		time.Sleep(pg.connTimeout)
		pg.connAttemps--
	}

	if err != nil {
		return nil, fmt.Errorf("postgres - NewPostgres - connAttemps == 0: %w", err)
	}
	return pg, nil
}

func (p *Postgres) Close() {
	if p.Pool != nil {
		p.Pool.Close()
	}
}

func (p *Postgres) RunInTx(ctx context.Context, f func(context.Context) error) error {
	tx, err := p.Pool.Begin(ctx)
	if err != nil {
		return err

	}
	ctx = context.WithValue(ctx, ctxTxKey{}, tx)

	err = f(ctx)
	if err != nil {
		_ = tx.Rollback(ctx)
		return err
	}
	return tx.Commit(ctx)
}

func (p *Postgres) GetQueryer(ctx context.Context) DB {
	if tx, ok := ctx.Value(ctxTxKey{}).(pgx.Tx); ok {
		return tx
	}
	return p.Pool
}
