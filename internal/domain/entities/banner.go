package entities

import (
	"time"

	"github.com/rendau/barot/internal/constant"
)

// Banner is type for banner
type Banner struct {
	ID       int64
	SlotID   int64
	ShowCnt  int64
	ClickCnt int64
}

type BannerCreatePars struct {
	ID     int64
	SlotID int64
	Note   string
}

type BannerDeletePars struct {
	ID     int64
	SlotID int64
}

type BannerSelectPars struct {
	SlotID    int64
	UsrTypeID int64
}

type BannerListPars struct {
	SlotID    int64
	UsrTypeID int64
}

type BannerStatIncPars struct {
	ID        int64
	SlotID    int64
	UsrTypeID int64
	Value     int64
}

type BannerEvent struct {
	Type      constant.BannerEventType
	BannerID  int64
	SlotID    int64
	UsrTypeID int64
	DateTime  time.Time
}
