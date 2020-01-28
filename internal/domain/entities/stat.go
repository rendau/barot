package entities

// Stat is type for stat
type Stat struct {
	BannerID  int64
	SlotID    int64
	UsrTypeID int64
	ShowCnt   int64
	ClickCnt  int64
}

// StatIncPars is type for stat increment params
type StatIncPars struct {
	BannerID  int64
	SlotID    int64
	UsrTypeID int64
}

// StatFilterPars is type for stat filter params
type StatFilterPars struct {
	BannerIDs *[]int64
	SlotID    *int64
	UsrTypeID *int64
}
