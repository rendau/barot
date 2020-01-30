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

	rep, err := http.Post(t.apiUrl+"/banners", "application/json", bytes.NewBuffer(dataBytes))
	if err != nil {
		return 0, err
	}
	defer rep.Body.Close()

	return rep.StatusCode, nil
}

func (t *Tests) uSelectBanner(slotId, usrTypeId int64) (int64, error) {
	rep, err := http.Get(t.apiUrl + fmt.Sprintf("/banners/select/%d/%d", slotId, usrTypeId))
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

func (t *Tests) uBannerAddClick(slotId, bannerId, usrTypeId int64) (int, error) {
	rep, err := http.Post(
		t.apiUrl+fmt.Sprintf("/banners/add_click/%d/%d/%d", slotId, bannerId, usrTypeId),
		"application/json",
		nil,
	)
	if err != nil {
		return 0, err
	}
	defer rep.Body.Close()

	return rep.StatusCode, nil
}

func (t *Tests) uDeleteBanner(slotId, bannerId int64) (int, error) {
	req, err := http.NewRequest("DELETE", t.apiUrl+fmt.Sprintf("/banners/%d/%d", slotId, bannerId), nil)
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
