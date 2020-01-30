package internal

// BannerC - is type for event create
type BannerC struct {
	ID     int64  `json:"id"`
	SlotId int64  `json:"slot_id"`
	Note   string `json:"note"`
}
