package entities

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
