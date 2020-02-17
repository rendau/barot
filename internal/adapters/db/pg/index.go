package pg

import (
	"context"
	"fmt"
	"time"

	"github.com/rendau/barot/internal/adapters/logger/zap"
	// driver for migration
	_ "github.com/jackc/pgx/v4/stdlib"
	"github.com/jmoiron/sqlx"
)

const (
	// ErrMsg is for ErrMsg
	ErrMsg        = "PG-error"
	dbWaitTimeout = 30 * time.Second
)

// PostgresDB - is type for pg-Db
type PostgresDB struct {
	lg *zap.Logger
	Db *sqlx.DB
}

// NewPostgresDB - creates new PostgresDB instance
func NewPostgresDB(dsn string, lg *zap.Logger) (*PostgresDB, error) {
	var err error

	if dsn == "" {
		return nil, fmt.Errorf("bad dsn for postgresql")
	}

	res := &PostgresDB{
		lg: lg,
	}

	connectionContext, connectionContextCancel := context.WithTimeout(context.Background(), dbWaitTimeout)
	defer connectionContextCancel()

	res.Db, err = res.dbWait(connectionContext, dsn)
	if err != nil {
		return nil, err
	}

	res.Db.SetMaxOpenConns(10)
	res.Db.SetMaxIdleConns(5)
	res.Db.SetConnMaxLifetime(10 * time.Minute)

	return res, nil
}

func (d *PostgresDB) dbWait(ctx context.Context, dsn string) (*sqlx.DB, error) {
	var err error

	var cnt uint32

	var db *sqlx.DB

	db, err = sqlx.Open("pgx", dsn)
	if err != nil {
		return nil, err
	}

	for {
		err = db.PingContext(ctx)
		if err == nil || ctx.Err() != nil {
			break
		}

		time.Sleep(time.Second)
	}

	if err != nil {
		return nil, err
	}

	for {
		err = db.GetContext(ctx, &cnt, `select count(*) from banner`)
		if err == nil || ctx.Err() != nil {
			break
		}

		time.Sleep(time.Second)
	}

	if err != nil {
		return nil, err
	}

	return db, nil
}
