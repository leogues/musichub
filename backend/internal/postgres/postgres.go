package postgres

import (
	"context"
	"database/sql"
	"database/sql/driver"
	"embed"
	"fmt"
	"io/fs"
	"sort"
	"time"

	musichub "github.com/leogues/MusicSyncHub"
	_ "github.com/lib/pq"
)

type DB struct {
	db     *sql.DB
	ctx    context.Context
	cancel func()

	DatasourceName string

	Now func() time.Time
}

func now() time.Time {
	return time.Now().Truncate(time.Second)
}

func NewDB(datasourceName string) *DB {
	db := &DB{
		DatasourceName: datasourceName,
		Now:            now,
	}

	db.ctx, db.cancel = context.WithCancel(context.Background())
	return db
}

func (db *DB) SqlDB() *sql.DB {
	return db.db
}

func (db *DB) Open() (err error) {
	if db.DatasourceName == "" {
		return fmt.Errorf("missing datasource name")
	}

	if db.db, err = sql.Open("postgres", db.DatasourceName); err != nil {
		return err
	}

	if err = db.db.Ping(); err != nil {
		return fmt.Errorf("ping database: %w", err)
	}

	if err = db.migrate(); err != nil {
		return fmt.Errorf("migrate: %w", err)
	}

	return nil
}

func (db *DB) Close() error {
	db.cancel()

	if db.db != nil {
		db.db.Close()
	}

	return nil
}

//go:embed migration/*.sql
var migrationFS embed.FS

func (db *DB) migrate() error {
	if _, err := db.db.Exec(`CREATE TABLE IF NOT EXISTS migrations (name TEXT PRIMARY KEY);`); err != nil {
		return fmt.Errorf("cannot create migrations table: %w", err)
	}

	// Read migration files from our embedded file system.
	// This uses Go 1.16's 'embed' package.
	names, err := fs.Glob(migrationFS, "migration/*.sql")
	if err != nil {
		return err
	}
	sort.Strings(names)

	for _, name := range names {
		if err := db.migrateFile(name); err != nil {
			return fmt.Errorf("migration error: name=%q err=%w", name, err)
		}
	}

	return nil
}

func (db *DB) migrateFile(name string) error {
	tx, err := db.db.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	var count int
	if err := tx.QueryRow(`SELECT COUNT(*) FROM migrations WHERE name = $1;`, name).Scan(&count); err != nil {
		return err
	} else if count != 0 {
		return nil
	}

	if buf, err := migrationFS.ReadFile(name); err != nil {
		return err
	} else if _, err := tx.Exec(string(buf)); err != nil {
		return err
	}

	if _, err := tx.Exec(`INSERT INTO migrations (name) VALUES ($1);`, name); err != nil {
		return err
	}

	return tx.Commit()
}

func FormatLimitOffset(limit, offset int) string {
	if limit > 0 && offset > 0 {
		return fmt.Sprintf(`LIMIT %d OFFSET %d`, limit, offset)
	} else if limit > 0 {
		return fmt.Sprintf(`LIMIT %d`, limit)
	} else if offset > 0 {
		return fmt.Sprintf(`OFFSET %d`, offset)
	}
	return ""
}

type TxOrDb interface {
	Exec(query string, args ...interface{}) (sql.Result, error)
	Query(query string, args ...interface{}) (*sql.Rows, error)
	ExecContext(ctx context.Context, query string, args ...interface{}) (sql.Result, error)
	QueryContext(ctx context.Context, query string, args ...interface{}) (*sql.Rows, error)
	QueryRowContext(ctx context.Context, query string, args ...interface{}) *sql.Row
}
type Tx struct {
	TxOrDb
	db  *DB
	now time.Time
}

func (db *DB) BeginTx(ctx context.Context) (*Tx, error) {
	tx, ok := musichub.TransactionFromContext(ctx)

	if ok {
		return &Tx{
			TxOrDb: tx,
			db:     db,
			now:    db.Now(),
		}, nil
	}

	sqlDB := db.db

	return &Tx{
		TxOrDb: sqlDB,
		db:     db,
		now:    db.Now(),
	}, nil
}

type NullTime time.Time

func (n *NullTime) Scan(value interface{}) error {
	if value == nil {
		*(*time.Time)(n) = time.Time{}
		return nil
	} else if value, ok := value.(time.Time); ok {
		*(*time.Time)(n) = value.Local()
		return nil
	}
	return fmt.Errorf("NullTime: cannot scan to time.Time: %T", value)
}

func (n *NullTime) Value() (driver.Value, error) {
	if n == nil || (*time.Time)(n).IsZero() {
		return nil, nil
	}
	return (*time.Time)(n).Local().Format(time.RFC3339), nil
}
