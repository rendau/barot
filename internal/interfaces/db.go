package interfaces

import (
	"context"
	"github.com/rendau/barot/internal/domain/entities"
)

type Db interface {
	// banner
	BannerCreate(ctx context.Context, pars entities.BannerCreatePars) error
	BannerDelete(ctx context.Context, pars entities.BannerDeletePars) error
	BannerList(ctx context.Context, pars entities.BannerListPars) ([]*entities.Banner, error)
	BannerIncShowCount(ctx context.Context, pars entities.BannerStatIncPars) error
	BannerIncClickCount(ctx context.Context, pars entities.BannerStatIncPars) error
}
