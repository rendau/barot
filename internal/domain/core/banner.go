package core

import (
	"context"
	"github.com/rendau/barot/internal/constant"
	"github.com/rendau/barot/internal/domain/entities"
	"math"
)

func (c *St) BannerCreate(ctx context.Context, obj entities.BannerCreatePars) error {
	return c.db.BannerCreate(ctx, obj)
}

func (c *St) BannerDelete(ctx context.Context, pars entities.BannerDeletePars) error {
	return c.db.BannerDelete(ctx, pars)
}

func (c *St) BannerSelectId(ctx context.Context, pars entities.BannerListPars) (int64, error) {
	banners, err := c.db.BannerList(ctx, pars)
	if err != nil {
		return 0, err
	}

	var allBannersShowCount int64

	// show count of all banners
	for _, banner := range banners {
		allBannersShowCount += banner.ShowCnt
	}

	var selectedBannerId int64

	// find the best of the best of the best :)
	var maxPoint float64
	var point float64
	for _, banner := range banners {
		point = c.MabCalc(banner.ShowCnt, banner.ClickCnt, allBannersShowCount)
		if point > maxPoint || selectedBannerId == 0 {
			selectedBannerId = banner.ID
			maxPoint = point
		}
	}

	// increment show counter
	if selectedBannerId > 0 {
		err = c.db.BannerIncShowCount(ctx, entities.BannerStatIncPars{
			ID:        selectedBannerId,
			SlotID:    pars.SlotID,
			UsrTypeID: pars.UsrTypeID,
		})
	}

	return selectedBannerId, nil
}

func (c *St) BannerAddClick(ctx context.Context, pars entities.BannerStatIncPars) error {
	return c.db.BannerIncClickCount(ctx, pars)
}

// MabCalc is calculates "multiarmed bandit" algorithm by input args
func (c *St) MabCalc(bannerShowCount, bannerClickCount, allBannersShowCount int64) float64 {
	if bannerShowCount == 0 {
		return constant.MabCalcInitValue
	}

	return (float64(bannerClickCount) / float64(bannerShowCount)) +
		math.Sqrt(2*math.Log(float64(allBannersShowCount))/float64(bannerShowCount))
}
