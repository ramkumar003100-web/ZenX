package database

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"

	"zenx/pkg/metrics"
)

type DB struct {
	SQL *sql.DB
}

func New(driver, dsn string) (*DB, error) {
	db, err := sql.Open(driver, dsn)
	if err != nil {
		return nil, err
	}
	if err := db.Ping(); err != nil {
		return nil, err
	}
	return &DB{SQL: db}, nil
}

func (d *DB) Close() error { return d.SQL.Close() }

func (d *DB) Migrate(ctx context.Context, dir string) error {
	entries, err := os.ReadDir(dir)
	if err != nil {
		return err
	}
	files := make([]string, 0, len(entries))
	for _, e := range entries {
		if !e.IsDir() && strings.HasSuffix(e.Name(), ".sql") {
			files = append(files, e.Name())
		}
	}
	sort.Strings(files)

	tx, err := d.SQL.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()

	for _, file := range files {
		stmt, readErr := os.ReadFile(filepath.Join(dir, file))
		if readErr != nil {
			return readErr
		}
		if _, execErr := tx.ExecContext(ctx, string(stmt)); execErr != nil {
			return fmt.Errorf("migration %s failed: %w", file, execErr)
		}
		metrics.DBQueries.WithLabelValues("migration").Inc()
	}
	return tx.Commit()
}

func (d *DB) ExecContext(ctx context.Context, query string, args ...any) (sql.Result, error) {
	metrics.DBQueries.WithLabelValues("exec").Inc()
	return d.SQL.ExecContext(ctx, query, args...)
}

func (d *DB) QueryContext(ctx context.Context, query string, args ...any) (*sql.Rows, error) {
	metrics.DBQueries.WithLabelValues("query").Inc()
	return d.SQL.QueryContext(ctx, query, args...)
}

func (d *DB) QueryRowContext(ctx context.Context, query string, args ...any) *sql.Row {
	metrics.DBQueries.WithLabelValues("query_row").Inc()
	return d.SQL.QueryRowContext(ctx, query, args...)
}

// ORM-like helpers
func (d *DB) Insert(ctx context.Context, table string, fields map[string]any) (sql.Result, error) {
	if len(fields) == 0 {
		return nil, errors.New("insert fields cannot be empty")
	}
	cols := make([]string, 0, len(fields))
	placeholders := make([]string, 0, len(fields))
	args := make([]any, 0, len(fields))
	for k, v := range fields {
		cols = append(cols, k)
		placeholders = append(placeholders, "?")
		args = append(args, v)
	}
	q := fmt.Sprintf("INSERT INTO %s (%s) VALUES (%s)", table, strings.Join(cols, ","), strings.Join(placeholders, ","))
	return d.ExecContext(ctx, q, args...)
}

func (d *DB) Delete(ctx context.Context, table, where string, args ...any) (sql.Result, error) {
	q := fmt.Sprintf("DELETE FROM %s WHERE %s", table, where)
	return d.ExecContext(ctx, q, args...)
}

type Query struct {
	table   string
	columns []string
	where   []string
	args    []any
	limit   int
}

func Table(name string) *Query { return &Query{table: name, columns: []string{"*"}} }

func (q *Query) Select(columns ...string) *Query {
	if len(columns) > 0 {
		q.columns = columns
	}
	return q
}

func (q *Query) Where(clause string, args ...any) *Query {
	q.where = append(q.where, clause)
	q.args = append(q.args, args...)
	return q
}

func (q *Query) Limit(v int) *Query { q.limit = v; return q }

func (q *Query) Build() (string, []any, error) {
	if q.table == "" {
		return "", nil, errors.New("missing table")
	}
	sql := fmt.Sprintf("SELECT %s FROM %s", strings.Join(q.columns, ","), q.table)
	if len(q.where) > 0 {
		sql += " WHERE " + strings.Join(q.where, " AND ")
	}
	if q.limit > 0 {
		sql += fmt.Sprintf(" LIMIT %d", q.limit)
	}
	return sql, q.args, nil
}
