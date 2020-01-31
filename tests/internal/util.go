package internal

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"
	"time"
)

func (t *Tests) uCreateBanner(banner BannerC) (int, error) {
	dataBytes, err := json.Marshal(banner)
	if err != nil {
		return 0, err
	}

	rep, err := http.Post(t.apiURL+"/banners", "application/json", bytes.NewBuffer(dataBytes))
	if err != nil {
		return 0, err
	}
	defer rep.Body.Close()

	return rep.StatusCode, nil
}

func (t *Tests) uSelectBanner(slotID, usrTypeID int64) (int64, error) {
	rep, err := http.Get(t.apiURL + fmt.Sprintf("/banners/select/%d/%d", slotID, usrTypeID))
	if err != nil {
		return 0, err
	}
	defer rep.Body.Close()

	repBytes, err := ioutil.ReadAll(rep.Body)
	if err != nil {
		return 0, err
	}

	return strconv.ParseInt(string(repBytes), 10, 64)
}

func (t *Tests) uBannerAddClick(slotID, bannerID, usrTypeID int64) (int, error) {
	rep, err := http.Post(
		t.apiURL+fmt.Sprintf("/banners/add_click/%d/%d/%d", slotID, bannerID, usrTypeID),
		"application/json",
		nil,
	)
	if err != nil {
		return 0, err
	}
	defer rep.Body.Close()

	return rep.StatusCode, nil
}

func (t *Tests) uDeleteBanner(slotID, bannerID int64) (int, error) {
	req, err := http.NewRequest("DELETE", t.apiURL+fmt.Sprintf("/banners/%d/%d", slotID, bannerID), nil)
	if err != nil {
		return 0, err
	}

	client := &http.Client{Timeout: 5 * time.Second}

	rep, err := client.Do(req)
	if err != nil {
		return 0, err
	}
	defer rep.Body.Close()

	return rep.StatusCode, nil
}
