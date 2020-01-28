package interfaces

import (
	"context"
	"github.com/rendau/barot/internal/domain/entities"
)

type Db interface {
	// banner
	BannerCreate(ctx context.Context, obj *entities.Banner) error
	BannerList(ctx context.Context, pars entities.BannerFilterPars) ([]*entities.Banner, error)
	BannerDelete(ctx context.Context, pars entities.BannerFilterPars) error

	// stat
	StatIncShowCount(ctx context.Context, pars *entities.StatIncPars) error
}
