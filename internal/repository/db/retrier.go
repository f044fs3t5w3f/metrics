package db

import (
	"context"
	"database/sql"
	"errors"
	"time"

	"github.com/jackc/pgconn"
	"github.com/jackc/pgerrcode"
)

type retrier struct {
	retryIntervals []time.Duration
	executor       executor
}

var _ executor = &retrier{}

func newRetrier(retryIntervals []time.Duration, executor executor) *retrier {
	return &retrier{
		retryIntervals,
		executor,
	}
}

func isRetriable(err error) bool {
	var pgErr *pgconn.PgError
	if errors.As(err, &pgErr) {
		return pgerrcode.IsConnectionException(pgErr.Code)
	} else {
		return false
	}
}

func (r *retrier) ExecContext(ctx context.Context, query string, args ...any) (sql.Result, error) {
	for _, delay := range r.retryIntervals {
		res, err := r.executor.ExecContext(ctx, query, args...)
		if err == nil || !isRetriable(err) {
			return res, err
		}
		time.Sleep(delay)
	}
	return r.executor.ExecContext(ctx, query, args...)
}

func (r *retrier) QueryContext(ctx context.Context, query string, args ...any) (*sql.Rows, error) {
	for _, delay := range r.retryIntervals {
		res, err := r.executor.QueryContext(ctx, query, args...)
		if err == nil || !isRetriable(err) {
			return res, err
		}
		time.Sleep(delay)
	}
	return r.executor.QueryContext(ctx, query, args...)
}

func (r *retrier) QueryRowContext(ctx context.Context, query string, args ...any) *sql.Row {
	for _, delay := range r.retryIntervals {
		res := r.executor.QueryRowContext(ctx, query, args...)
		err := res.Err()
		if err == nil || !isRetriable(err) {
			return res
		}
		time.Sleep(delay)
	}
	return r.executor.QueryRowContext(ctx, query, args...)
}
