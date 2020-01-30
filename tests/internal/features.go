package internal

import (
	"encoding/json"
	"fmt"
	"github.com/DATA-DOG/godog"
	"github.com/DATA-DOG/godog/gherkin"
)

func (t *Tests) iRequestToCreateBannerWithData(bannerId int, data *gherkin.DocString) error {
	banner := BannerC{}
	err := json.Unmarshal([]byte(data.Content), &banner)
	if err != nil {
		return err
	}

	sc, err := t.uCreateBanner(banner)
	if err != nil {
		return err
	}

	t.responseStatusCode = sc

	return nil
}

func (t *Tests) theResponseCodeShouldBe(code int) error {
	if t.responseStatusCode != code {
		return fmt.Errorf("bad status code: %d", code)
	}
	return nil
}

func (t *Tests) iRequestBannerToShowForTimes(n int) error {
	for i := 0; i < n; i++ {
		bannerId, err := t.uSelectBanner(slot1Id, usrType1Id)
		if err != nil {
			return err
		}
		t.showCounts[bannerId] += 1
	}
	return nil
}

func (t *Tests) iWillGetShowsForBanner(showCount, bannerId int) error {
	if t.showCounts[int64(bannerId)] != int64(showCount) {
		return fmt.Errorf("show counts not equal, expected %d, got %d", showCount, t.showCounts[int64(bannerId)])
	}
	return nil
}

func (t *Tests) iRequestClickForBanner(bannerId int) error {
	sc, err := t.uBannerAddClick(slot1Id, int64(bannerId), usrType1Id)
	if err != nil {
		return err
	}

	t.responseStatusCode = sc

	return nil
}

func (t *Tests) bannerShowCountMustBeGreaterThanBannerShowCount(bannerId1, bannerId2 int) error {
	showCount1 := t.showCounts[int64(bannerId1)]
	showCount2 := t.showCounts[int64(bannerId2)]
	if showCount1 <= showCount2 {
		return fmt.Errorf("%d is not grater than %d", showCount1, showCount2)
	}
	return nil
}

func (t *Tests) iRequestToDeleteBanner(bannerId int) error {
	sc, err := t.uDeleteBanner(slot1Id, int64(bannerId))
	if err != nil {
		return err
	}

	t.responseStatusCode = sc

	return nil
}

func (t *Tests) FeatureContext(s *godog.Suite) {
	s.Step(`^I request to create banner (\d+) with data:$`, t.iRequestToCreateBannerWithData)
	s.Step(`^The response code should be (\d+)$`, t.theResponseCodeShouldBe)
	s.Step(`^I request banner to show for (\d+) times$`, t.iRequestBannerToShowForTimes)
	s.Step(`^I will get (\d+) shows for banner (\d+)$`, t.iWillGetShowsForBanner)
	s.Step(`^I request click for banner (\d+):$`, t.iRequestClickForBanner)
	s.Step(`^banner (\d+) show count must be greater than banner (\d+) show count$`, t.bannerShowCountMustBeGreaterThanBannerShowCount)
	s.Step(`^I request to delete banner (\d+)$`, t.iRequestToDeleteBanner)
}
