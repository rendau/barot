package constant

type BannerEventType string

const (
	MabCalcInitValue = float64(999999) //nolint

	BannerEventTypeShow  = BannerEventType("show")
	BannerEventTypeClick = BannerEventType("click")
)

func (s BannerEventType) String() string {
	return string(s)
}
