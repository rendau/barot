package core

import (
	"context"
	"math"
	"time"

	"github.com/rendau/barot/internal/constant"
	"github.com/rendau/barot/internal/domain/entities"
)

func (c *St) BannerCreate(ctx context.Context, obj entities.BannerCreatePars) error {
	return c.db.BannerCreate(ctx, obj)
}

func (c *St) BannerDelete(ctx context.Context, pars entities.BannerDeletePars) error {
	return c.db.BannerDelete(ctx, pars)
}

func (c *St) BannerSelectID(ctx context.Context, pars entities.BannerSelectPars) (int64, error) {
	banners, err := c.db.BannerList(ctx, entities.BannerListPars(pars))
	if err != nil {
		return 0, err
	}

	var allBannersShowCount int64

	// show count of all banners
	for _, banner := range banners {
		allBannersShowCount += banner.ShowCnt
	}

	var selectedBannerID int64

	// find the best of the best of the best :)
	var maxPoint, point float64
	for _, banner := range banners {
		point = c.MabCalc(banner.ShowCnt, banner.ClickCnt, allBannersShowCount)
		if point > maxPoint || selectedBannerID == 0 {
			selectedBannerID = banner.ID
			maxPoint = point
		}
	}

	// increment show counter
	if selectedBannerID > 0     {
		err = c.db.BannerIncShowCount(ctx, entities.BannerStatIncPars{
			ID:        selectedBannerID,
			SlotID:    pars.SlotID,
			UsrTypeID: pars.UsrTypeID,
		})
		if err != nil {
			return 0, err
		}

		err = c.mq.PublishBannerEvent(&entities.BannerEvent{
			Type:      constant.BannerEventTypeShow,
			BannerID:  selectedBannerID,
			SlotID:    pars.SlotID,
			UsrTypeID: pars.UsrTypeID,
			DateTime:  time.Now(),
		})
		if err != nil {
			return 0, err
		}
	}

	return selectedBannerID, nil
}

func (c *St) BannerAddClick(ctx context.Context, pars entities.BannerStatIncPars) error {
	err := c.db.BannerIncClickCount(ctx, pars)
	if err != nil {
		return err
	}

	err = c.mq.PublishBannerEvent(&entities.BannerEvent{
		Type:      constant.BannerEventTypeClick,
		BannerID:  pars.ID,
		SlotID:    pars.SlotID,
		UsrTypeID: pars.UsrTypeID,
		DateTime:  time.Now(),
	})
	if err != nil {
		return err
	}

	return nil
}

// MabCalc is calculates "multiarmed bandit" algorithm by input args
func (c *St) MabCalc(bannerShowCount, bannerClickCount, allBannersShowCount int64) float64 {
	if bannerShowCount == 0 {
		return constant.MabCalcInitValue
	}

	return (float64(bannerClickCount) / float64(bannerShowCount)) +
		math.Sqrt(2*math.Log(float64(allBannersShowCount))/float64(bannerShowCount))
}
