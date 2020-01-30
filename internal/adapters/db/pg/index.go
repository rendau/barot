package pg

import (
	"context"
	"fmt"
	"github.com/rendau/barot/internal/adapters/logger/zap"

	// driver for migration
	_ "github.com/jackc/pgx/v4/stdlib"
	"github.com/jmoiron/sqlx"
	"time"
)

const (
	ErrMsg        = "PG-error"
	dbWaitTimeout = 30 * time.Second
)

// St - is type for pg-Db
type St struct {
	lg *zap.St
	Db *sqlx.DB
}

// NewSt - creates new St instance
func NewSt(dsn string, lg *zap.St) (*St, error) {
	var err error

	if dsn == "" {
		return nil, fmt.Errorf("bad dsn for postgresql")
	}

	res := &St{
		lg: lg,
	}

	connectionContext, _ := context.WithTimeout(context.Background(), dbWaitTimeout)
	res.Db, err = res.dbWait(dsn, connectionContext)
	if err != nil {
		return nil, err
	}

	res.Db.SetMaxOpenConns(10)
	res.Db.SetMaxIdleConns(5)
	res.Db.SetConnMaxLifetime(10 * time.Minute)

	return res, nil
}

func (d *St) dbWait(dsn string, ctx context.Context) (*sqlx.DB, error) {
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
