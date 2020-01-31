package constant

// BannerEventType is type for BannerEventType
type BannerEventType string

const (
	// MabCalcInitValue is constant
	MabCalcInitValue = float64(999999) //nolint

	// BannerEventTypeShow is constant
	BannerEventTypeShow = BannerEventType("show")
	// BannerEventTypeClick is constant
	BannerEventTypeClick = BannerEventType("click")
)

func (s BannerEventType) String() string {
	return string(s)
}
