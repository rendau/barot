package entities

// Banner is type for banner
type Banner struct {
	ID     int64
	SlotID int64
	Note   string
}

// BannerFilterPars is type for banner filter params
type BannerFilterPars struct {
	ID     *int64
	SlotID *int64
}
