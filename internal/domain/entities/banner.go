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

// BannerCreatePars is type for BannerCreatePars
type BannerCreatePars struct {
	ID     int64
	SlotID int64
	Note   string
}

// BannerDeletePars is type for BannerDeletePars
type BannerDeletePars struct {
	ID     int64
	SlotID int64
}

// BannerSelectPars is type for BannerSelectPars
type BannerSelectPars struct {
	SlotID    int64
	UsrTypeID int64
}

// BannerListPars is type for BannerListPars
type BannerListPars struct {
	SlotID    int64
	UsrTypeID int64
}

// BannerStatIncPars is type for BannerStatIncPars
type BannerStatIncPars struct {
	ID        int64
	SlotID    int64
	UsrTypeID int64
	Value     int64
}

// BannerEvent is type for banner-event
type BannerEvent struct {
	Type      constant.BannerEventType
	BannerID  int64
	SlotID    int64
	UsrTypeID int64
	DateTime  time.Time
}
