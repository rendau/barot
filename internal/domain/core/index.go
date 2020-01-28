package core

import "github.com/rendau/barot/internal/interfaces"

type St struct {
	lg interfaces.Logger
	db interfaces.Db
}

func NewSt(lg interfaces.Logger, db interfaces.Db) *St {
	return &St{
		lg: lg,
		db: db,
	}
}
