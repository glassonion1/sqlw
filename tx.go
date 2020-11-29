package sqlw

import (
	"context"
	"database/sql"
)

// TxHandlerFunc is for executing SQL on a transaction.
// To make the SQL to be executed a transition target, must execute it via the type sqlw.Tx.
type TxHandlerFunc func(*Tx) error

// Tx is a wrapper around sql.Tx
type Tx struct {
	parent *sql.Tx
}

// Query executes a query that returns rows, typically a SELECT.
func (tx *Tx) Query(query SQLQuery, args ...interface{}) (*sql.Rows, error) {
	if err := query.Validate(); err != nil {
		return nil, err
	}
	return tx.parent.Query(query.String(), args...)
}

// QueryContext executes a query that returns rows, typically a SELECT.
func (tx *Tx) QueryContext(ctx context.Context, query SQLQuery, args ...interface{}) (*sql.Rows, error) {
	if err := query.Validate(); err != nil {
		return nil, err
	}
	return tx.parent.QueryContext(ctx, query.String(), args...)
}

// QueryRow executes a query that is expected to return at most one row. QueryRow always returns a non-nil value. Errors are deferred until Row's Scan method is called. If the query selects no rows, the *Row's Scan will return ErrNoRows. Otherwise, the *Row's Scan scans the first selected row and discards the rest.
func (tx *Tx) QueryRow(query SQLQuery, args ...interface{}) *sql.Row {
	if err := query.Validate(); err != nil {
		return nil
	}
	return tx.parent.QueryRow(query.String(), args...)
}

// QueryRowContext executes a query that is expected to return at most one row. QueryRow always returns a non-nil value. Errors are deferred until Row's Scan method is called. If the query selects no rows, the *Row's Scan will return ErrNoRows. Otherwise, the *Row's Scan scans the first selected row and discards the rest.
func (tx *Tx) QueryRowContext(ctx context.Context, query SQLQuery, args ...interface{}) *sql.Row {
	if err := query.Validate(); err != nil {
		return nil
	}
	return tx.parent.QueryRowContext(ctx, query.String(), args...)
}

// Exec executes a query without returning any rows. The args are for any placeholder parameters in the query.
func (tx *Tx) Exec(query SQLMutation, args ...interface{}) (sql.Result, error) {
	if err := query.Validate(); err != nil {
		return nil, err
	}
	return tx.parent.Exec(query.String(), args...)
}

// ExecContext executes a query without returning any rows. The args are for any placeholder parameters in the query.
func (tx *Tx) ExecContext(ctx context.Context, query SQLMutation, args ...interface{}) (sql.Result, error) {
	if err := query.Validate(); err != nil {
		return nil, err
	}
	return tx.parent.ExecContext(ctx, query.String(), args...)
}
