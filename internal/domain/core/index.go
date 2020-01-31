package core

import "github.com/rendau/barot/internal/interfaces"

// St is type for st
type St struct {
	lg interfaces.Logger
	db interfaces.Db
	mq interfaces.Mq
}

// NewSt is create new instance of St
func NewSt(lg interfaces.Logger, db interfaces.Db, mq interfaces.Mq) *St {
	return &St{
		lg: lg,
		db: db,
		mq: mq,
	}
}
