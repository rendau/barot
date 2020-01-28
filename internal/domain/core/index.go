package core

import (
	"github.com/rendau/barot/internal/constant"
	"math"
)

type St struct {
}

func NewSt() *St {
	return &St{}
}

// MabCalc is calculates "multiarmed bandit" algorithm by input args
func (c *St) MabCalc(bannerShowCount, bannerClickCount, allBannersShowCount int64) float64 {
	if bannerShowCount == 0 {
		return constant.MabCalcInitValue
	}

	return (float64(bannerClickCount) / float64(bannerShowCount)) +
		math.Sqrt(2*math.Log(float64(allBannersShowCount))/float64(bannerShowCount))
}
