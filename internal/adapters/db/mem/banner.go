package mem

import (
	"context"
	"fmt"
	"strconv"
	"strings"

	"github.com/rendau/barot/internal/domain/entities"
)

const defaultIncrementValue int64 = 1

type bannerSt struct {
	id     int64
	slotID int64
	note   string
	stats  map[int64]*statSt
}

type statSt struct {
	showCount  int64
	clickCount int64
}

func (d *MemoryDB) getBannerSlotKey(bannerID, slotID int64) string {
	return fmt.Sprintf("%d_%d", bannerID, slotID)
}

// BannerCreate is for BannerCreate
func (d *MemoryDB) BannerCreate(ctx context.Context, pars entities.BannerCreatePars) error {
	d.banners[d.getBannerSlotKey(pars.ID, pars.SlotID)] = &bannerSt{
		id:     pars.ID,
		slotID: pars.SlotID,
		note:   pars.Note,
		stats:  make(map[int64]*statSt),
	}

	return nil
}

// BannerDelete is for BannerDelete
func (d *MemoryDB) BannerDelete(ctx context.Context, pars entities.BannerDeletePars) error {
	delete(d.banners, d.getBannerSlotKey(pars.ID, pars.SlotID))

	return nil
}

// BannerList is for BannerList
func (d *MemoryDB) BannerList(ctx context.Context, pars entities.BannerListPars) ([]*entities.Banner, error) {
	result := make([]*entities.Banner, 0)

	keySuffix := "_" + strconv.FormatInt(pars.SlotID, 10)

	for k, v := range d.banners {
		if !strings.HasSuffix(k, keySuffix) {
			continue
		}

		banner := &entities.Banner{
			ID:       v.id,
			SlotID:   v.slotID,
			ShowCnt:  0,
			ClickCnt: 0,
		}

		stat := v.stats[pars.UsrTypeID]

		if stat != nil {
			banner.ShowCnt = stat.showCount
			banner.ClickCnt = stat.clickCount
		}

		result = append(result, banner)
	}

	return result, nil
}

// BannerIncShowCount is for BannerIncShowCount
func (d *MemoryDB) BannerIncShowCount(ctx context.Context, pars entities.BannerStatIncPars) error {
	banner := d.banners[d.getBannerSlotKey(pars.ID, pars.SlotID)]
	if banner == nil {
		return fmt.Errorf("banner does not exist")
	}

	v := defaultIncrementValue
	if pars.Value != 0 {
		v = pars.Value
	}

	stat := banner.stats[pars.UsrTypeID]

	if stat == nil {
		banner.stats[pars.UsrTypeID] = &statSt{showCount: v}
	} else {
		stat.showCount += v
	}

	return nil
}

// BannerIncClickCount is for BannerIncClickCount
func (d *MemoryDB) BannerIncClickCount(ctx context.Context, pars entities.BannerStatIncPars) error {
	banner := d.banners[d.getBannerSlotKey(pars.ID, pars.SlotID)]
	if banner == nil {
		return fmt.Errorf("banner does not exist")
	}

	v := defaultIncrementValue
	if pars.Value != 0 {
		v = pars.Value
	}

	stat := banner.stats[pars.UsrTypeID]

	if stat == nil {
		banner.stats[pars.UsrTypeID] = &statSt{showCount: v}
	} else {
		stat.clickCount += v
	}

	return nil
}
