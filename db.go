package sqlw

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"math/rand"
	"strings"
	"sync"
	"time"

	_ "github.com/go-sql-driver/mysql"
)

// DB is a wrapper around sql.DB
type DB struct {
	master       *sql.DB
	readreplicas []*sql.DB
	mu           sync.Mutex
}

// NewMySQLDB returns a new sqlx DB wrapper for a pre-existing *sql.DB
//
// This function should be used outside of Goroutine.
func NewMySQLDB(masterConf Config, replicaConfs ...Config) (*DB, error) {
	master, err := sql.Open("mysql", masterConf.mysqlStr())
	if err != nil {
		return nil, err
	}

	replicas := []*sql.DB{}
	for _, conf := range replicaConfs {
		r, err := sql.Open("mysql", conf.mysqlStr())
		if err != nil {
			continue
		}
		if err := r.Ping(); err != nil {
			continue
		}
		replicas = append(replicas, r)
	}
	return NewDB(master, replicas...), nil
}

// NewPostgresDB returns a new sqlx DB wrapper for a pre-existing *sql.DB
//
// This function should be used outside of Goroutine.
func NewPostgresDB(masterConf Config, replicaConfs ...Config) (*DB, error) {
	master, err := sql.Open("postgres", masterConf.postgresStr())
	if err != nil {
		return nil, err
	}

	replicas := []*sql.DB{}
	for _, conf := range replicaConfs {
		r, err := sql.Open("postgres", conf.postgresStr())
		if err != nil {
			continue
		}
		if err := r.Ping(); err != nil {
			continue
		}
		replicas = append(replicas, r)
	}
	return NewDB(master, replicas...), nil
}

// NewDB returns a new sqlx DB wrapper for a pre-existing *sql.DB
//
// This function should be used outside of Goroutine.
func NewDB(master *sql.DB, readreplicas ...*sql.DB) *DB {
	rand.Seed(time.Now().UnixNano())

	list := []*sql.DB{}
	for _, r := range readreplicas {
		if r != nil {
			list = append(list, r)
		}
	}
	return &DB{
		master:       master,
		readreplicas: list,
	}
}

func (db *DB) getReplica() *sql.DB {
	if len(db.readreplicas) == 0 {
		return db.master
	}
	return db.readreplicas[rand.Intn(len(db.readreplicas))]
}

// Close closes all databases.
func (db *DB) Close() error {
	errList := []string{}
	if err := db.master.Close(); err != nil {
		errList = append(errList, err.Error())
	}

	for _, r := range db.readreplicas {
		if rerr := r.Close(); rerr != nil {
			errList = append(errList, rerr.Error())
		}
	}
	if len(errList) > 0 {
		str := strings.Join(errList, ",")
		return errors.New(str)
	}
	return nil
}

// SetConnMaxLifetime sets the maximum amount of time a connection may be reused.
func (db *DB) SetConnMaxLifetime(d time.Duration) {
	db.master.SetConnMaxLifetime(d)
	for _, r := range db.readreplicas {
		r.SetConnMaxLifetime(d)
	}
}

// SetMaxIdleConns sets the maximum number of connections in the idle connection pool.
func (db *DB) SetMaxIdleConns(n int) {
	db.master.SetMaxIdleConns(n)
	for _, r := range db.readreplicas {
		r.SetMaxIdleConns(n)
	}
}

// SetMaxOpenConns sets the maximum number of open connections to the database.
func (db *DB) SetMaxOpenConns(n int) {
	db.master.SetMaxOpenConns(n)
	for _, r := range db.readreplicas {
		r.SetMaxOpenConns(n)
	}
}

// Readable checks if the database can be readable.
func (db *DB) Readable() error {
	errList := []string{}

	if err := db.master.Ping(); err != nil {
		errList = append(errList, fmt.Sprintf("failed to ping master: %v", err))
	}

	for i, r := range db.readreplicas {
		if err := r.Ping(); err != nil {
			str := fmt.Sprintf("failed to ping replica%d: %v", i, err)
			errList = append(errList, str)
		}
	}

	if len(errList) > 0 {
		str := strings.Join(errList, ",")
		return errors.New(str)
	}

	return nil
}

// Writable checks if the database is writable.
func (db *DB) Writable() error {
	return db.master.Ping()
}

// Query executes a query that returns rows, typically a SELECT.
// This method is executed on the read replica.
func (db *DB) Query(query SQLQuery, args ...interface{}) (*sql.Rows, error) {
	if err := query.Validate(); err != nil {
		return nil, err
	}
	return db.getReplica().Query(query.String(), args...)
}

// QueryForMaster executes a query that returns rows, typically a SELECT.
// This method is executed on the master.
//
// It is used to refer to the data immediately after executing the mutation query(INSERT/UPDATE/DELETE).
func (db *DB) QueryForMaster(query SQLQuery, args ...interface{}) (*sql.Rows, error) {
	if err := query.Validate(); err != nil {
		return nil, err
	}
	return db.master.Query(query.String(), args...)
}

// QueryContext executes a query that returns rows, typically a SELECT.
// This method is executed on the read replica.
func (db *DB) QueryContext(ctx context.Context, query SQLQuery, args ...interface{}) (*sql.Rows, error) {
	if err := query.Validate(); err != nil {
		return nil, err
	}
	return db.getReplica().QueryContext(ctx, query.String(), args...)
}

// QueryContextForMaster executes a query that returns rows, typically a SELECT.
// This method is executed on the master.
//
// It is used to refer to the data immediately after executing the mutation query(INSERT/UPDATE/DELETE).
func (db *DB) QueryContextForMaster(ctx context.Context, query SQLQuery, args ...interface{}) (*sql.Rows, error) {
	if err := query.Validate(); err != nil {
		return nil, err
	}
	return db.master.QueryContext(ctx, query.String(), args...)
}

// QueryRow executes a query that is expected to return at most one row. QueryRow always returns a non-nil value. Errors are deferred until Row's Scan method is called. If the query selects no rows, the *Row's Scan will return ErrNoRows. Otherwise, the *Row's Scan scans the first selected row and discards the rest.
// This method is executed on the read replica.
func (db *DB) QueryRow(query SQLQuery, args ...interface{}) *sql.Row {
	if err := query.Validate(); err != nil {
		return nil
	}
	return db.getReplica().QueryRow(query.String(), args...)
}

// QueryRowForMaster executes a query that is expected to return at most one row. QueryRow always returns a non-nil value. Errors are deferred until Row's Scan method is called. If the query selects no rows, the *Row's Scan will return ErrNoRows. Otherwise, the *Row's Scan scans the first selected row and discards the rest.
// This method is executed on the master.
func (db *DB) QueryRowForMaster(query SQLQuery, args ...interface{}) *sql.Row {
	if err := query.Validate(); err != nil {
		return nil
	}
	return db.master.QueryRow(query.String(), args...)
}

// QueryRowContext executes a query that is expected to return at most one row. QueryRow always returns a non-nil value. Errors are deferred until Row's Scan method is called. If the query selects no rows, the *Row's Scan will return ErrNoRows. Otherwise, the *Row's Scan scans the first selected row and discards the rest.
// This method is executed on the read replica.
func (db *DB) QueryRowContext(ctx context.Context, query SQLQuery, args ...interface{}) *sql.Row {
	if err := query.Validate(); err != nil {
		return nil
	}
	return db.getReplica().QueryRowContext(ctx, query.String(), args...)
}

// QueryRowContextForMaster executes a query that is expected to return at most one row. QueryRow always returns a non-nil value. Errors are deferred until Row's Scan method is called. If the query selects no rows, the *Row's Scan will return ErrNoRows. Otherwise, the *Row's Scan scans the first selected row and discards the rest.
// This method is executed on the master.
func (db *DB) QueryRowContextForMaster(ctx context.Context, query SQLQuery, args ...interface{}) *sql.Row {
	if err := query.Validate(); err != nil {
		return nil
	}
	return db.master.QueryRowContext(ctx, query.String(), args...)
}

// PrepareQuery creates a prepared statement for later queries.The caller must call the statement's Close method when the statement is no longer needed.
// This method is executed on the read replica and can use for SELECT statements only.
func (db *DB) PrepareQuery(query SQLQuery) (*sql.Stmt, error) {
	if err := query.Validate(); err != nil {
		return nil, err
	}
	return db.getReplica().Prepare(query.String())
}

// PrepareQueryForMaster creates a prepared statement for later queries(SELECT).The caller must call the statement's Close method when the statement is no longer needed.
// This method is executed on the master and can use for SELECT statements only.
func (db *DB) PrepareQueryForMaster(query SQLQuery) (*sql.Stmt, error) {
	if err := query.Validate(); err != nil {
		return nil, err
	}
	return db.master.Prepare(query.String())
}

// PrepareQueryContext creates a prepared statement for later queries.The caller must call the statement's Close method when the statement is no longer needed.
// This method is executed on the read replica and can use for SELECT statements only.
func (db *DB) PrepareQueryContext(ctx context.Context, query SQLQuery) (*sql.Stmt, error) {
	if err := query.Validate(); err != nil {
		return nil, err
	}
	return db.getReplica().PrepareContext(ctx, query.String())
}

// PrepareQueryContextForMaster creates a prepared statement for later queries(SELECT).The caller must call the statement's Close method when the statement is no longer needed.
// This method is executed on the master and can use for SELECT statements only.
func (db *DB) PrepareQueryContextForMaster(ctx context.Context, query SQLQuery) (*sql.Stmt, error) {
	if err := query.Validate(); err != nil {
		return nil, err
	}
	return db.master.PrepareContext(ctx, query.String())
}

// PrepareMutation creates a prepared statement for later executions.The caller must call the statement's Close method when the statement is no longer needed.
// This method is executed on the master and can use for INSERT|UPDATE|DELETE statements only.
func (db *DB) PrepareMutation(query SQLMutation) (*sql.Stmt, error) {
	if err := query.Validate(); err != nil {
		return nil, err
	}
	return db.master.Prepare(query.String())
}

// PrepareMutationContext creates a prepared statement for later executions.The caller must call the statement's Close method when the statement is no longer needed.
// This method is executed on the master and can use for INSERT|UPDATE|DELETE statements only.
func (db *DB) PrepareMutationContext(ctx context.Context, query SQLMutation) (*sql.Stmt, error) {
	if err := query.Validate(); err != nil {
		return nil, err
	}
	return db.master.PrepareContext(ctx, query.String())
}

// Exec executes a query without returning any rows. The args are for any placeholder parameters in the query.
// This method is executed on the master and can use for INSERT|UPDATE|DELETE statements only.
func (db *DB) Exec(query SQLMutation, args ...interface{}) (sql.Result, error) {
	if err := query.Validate(); err != nil {
		return nil, err
	}
	return db.master.Exec(query.String(), args...)
}

// ExecContext executes a query without returning any rows. The args are for any placeholder parameters in the query.
// This method is executed on the master and can use for INSERT|UPDATE|DELETE statements only.
func (db *DB) ExecContext(ctx context.Context, query SQLMutation, args ...interface{}) (sql.Result, error) {
	if err := query.Validate(); err != nil {
		return nil, err
	}
	return db.master.ExecContext(ctx, query.String(), args...)
}

// Transaction executes paramed function in one database transaction. Executes the passed function and commits the transaction if there is no error. If an error occurs when executing the passed function rolls back the transaction.
// see sqlw/TxHandlerFunc
func (db *DB) Transaction(fn TxHandlerFunc) error {
	db.mu.Lock()
	defer db.mu.Unlock()

	origin, err := db.master.Begin()
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %v", err)
	}
	tx := &Tx{origin}

	if err := fn(tx); err != nil {
		if re := tx.parent.Rollback(); re != nil {
			if re.Error() != sql.ErrTxDone.Error() {
				return fmt.Errorf("fialed to rollback: %v", err)
			}
		}
		return fmt.Errorf("failed to execcute transaction: %v", err)
	}
	return tx.parent.Commit()
}

// TransactionTx executes paramed function in one database transaction. Executes the passed function and commits the transaction if there is no error. If an error occurs when executing the passed function rolls back the transaction.
// see sqlw/TxHandlerFunc
func (db *DB) TransactionTx(ctx context.Context, fn TxHandlerFunc, opts *sql.TxOptions) error {
	db.mu.Lock()
	defer db.mu.Unlock()

	origin, err := db.master.BeginTx(ctx, opts)
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %v", err)
	}
	tx := &Tx{origin}

	if err := fn(tx); err != nil {
		if re := tx.parent.Rollback(); re != nil {
			if re.Error() != sql.ErrTxDone.Error() {
				return fmt.Errorf("fialed to rollback: %v", err)
			}
		}
		return fmt.Errorf("failed to execcute transaction: %v", err)
	}
	return tx.parent.Commit()
}
