package httpapi

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"github.com/rendau/barot/internal/domain/entities"
)

func (a *API) hBannerAdd(w http.ResponseWriter, r *http.Request) {
	var err error

	var reqObj struct {
		ID     int64  `json:"id"`
		SlotID int64  `json:"slot_id"`
		Note   string `json:"note"`
	}

	if err = json.NewDecoder(r.Body).Decode(&reqObj); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	err = a.cr.BannerCreate(r.Context(), entities.BannerCreatePars{
		ID:     reqObj.ID,
		SlotID: reqObj.SlotID,
		Note:   reqObj.Note,
	})
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func (a *API) hBannerRemove(w http.ResponseWriter, r *http.Request) {
	var err error

	args := mux.Vars(r)
	slotID, _ := strconv.ParseInt(args["slot_id"], 10, 64)
	bannerID, _ := strconv.ParseInt(args["banner_id"], 10, 64)

	err = a.cr.BannerDelete(r.Context(), entities.BannerDeletePars{
		ID:     bannerID,
		SlotID: slotID,
	})
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func (a *API) hBannerSelect(w http.ResponseWriter, r *http.Request) {
	var err error

	args := mux.Vars(r)
	slotID, _ := strconv.ParseInt(args["slot_id"], 10, 64)
	usrTypeID, _ := strconv.ParseInt(args["usr_type_id"], 10, 64)

	var id int64

	id, err = a.cr.BannerSelectID(r.Context(), entities.BannerSelectPars{
		SlotID:    slotID,
		UsrTypeID: usrTypeID,
	})

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json; charset=utf-8")

	_, err = w.Write([]byte(strconv.FormatInt(id, 10)))
	if err != nil {
		a.lg.Errorw("Fail to respond data", err)
	}
}

func (a *API) hBannerAddClick(w http.ResponseWriter, r *http.Request) {
	var err error

	args := mux.Vars(r)
	slotID, _ := strconv.ParseInt(args["slot_id"], 10, 64)
	bannerID, _ := strconv.ParseInt(args["banner_id"], 10, 64)
	usrTypeID, _ := strconv.ParseInt(args["usr_type_id"], 10, 64)

	err = a.cr.BannerAddClick(r.Context(), entities.BannerStatIncPars{
		ID:        bannerID,
		SlotID:    slotID,
		UsrTypeID: usrTypeID,
	})
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}
