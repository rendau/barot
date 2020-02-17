package mem

import (
	"github.com/rendau/barot/internal/adapters/logger/zap"
	// driver for migration
	_ "github.com/jackc/pgx/v4/stdlib"
)

// MemoryDB - is type for mem-Db
type MemoryDB struct {
	lg *zap.Logger

	banners map[string]*bannerSt
}

// NewMemoryDB - creates new MemoryDB instance
func NewMemoryDB(lg *zap.Logger) *MemoryDB {
	return &MemoryDB{
		lg: lg,

		banners: make(map[string]*bannerSt),
	}
}

// Clean - cleans mem-db
func (d *MemoryDB) Clean() {
	d.banners = make(map[string]*bannerSt)
}
