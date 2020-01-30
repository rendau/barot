package core

import "github.com/rendau/barot/internal/interfaces"

type St struct {
	lg interfaces.Logger
	db interfaces.Db
	mq interfaces.Mq
}

func NewSt(lg interfaces.Logger, db interfaces.Db, mq interfaces.Mq) *St {
	return &St{
		lg: lg,
		db: db,
		mq: mq,
	}
}
