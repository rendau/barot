package main

import (
	"context"
	"log"
	"os"
	"testing"

	dbMem "github.com/rendau/barot/internal/adapters/db/mem"
	"github.com/rendau/barot/internal/adapters/logger/zap"
	mqMock "github.com/rendau/barot/internal/adapters/mq/mock"
	"github.com/rendau/barot/internal/constant"
	"github.com/rendau/barot/internal/domain/core"
	"github.com/rendau/barot/internal/domain/entities"
	"github.com/stretchr/testify/require"
)

const bannerNote = "some banner note"

var (
	app = struct {
		lg *zap.Logger
		db *dbMem.MemoryDB
		mq *mqMock.MessageQueueMock
		cr *core.St
	}{}
)

func TestMain(m *testing.M) {
	var err error

	app.lg, err = zap.NewLogger("debug", true, false)
	if err != nil {
		log.Fatal(err)
	}

	defer app.lg.Sync()

	app.db = dbMem.NewMemoryDB(app.lg)

	app.mq = mqMock.NewMessageQueueMock()

	app.cr = core.NewSt(app.lg, app.db, app.mq)

	ec := m.Run()

	os.Exit(ec)
}

func cleanCtx() {
	app.db.Clean()
	app.mq.Clean()
}

func TestMabCalc(t *testing.T) {
	v := app.cr.MultiArmedBanditCalculate(0, 0, 100)
	require.Equal(t, constant.MabCalcInitValue, v)
}

func TestBannerCreation(t *testing.T) {
	var err error

	cleanCtx()

	ctx := context.Background()

	var banner1Id int64 = 1

	var banner2Id int64 = 2

	var slot1Id int64 = 1

	var slot2Id int64 = 2

	err = app.cr.BannerCreate(ctx, entities.BannerCreatePars{
		ID:     banner1Id,
		SlotID: slot1Id,
		Note:   bannerNote,
	})
	require.Nil(t, err)

	banners, err := app.db.BannerList(ctx, entities.BannerListPars{
		SlotID:    slot1Id,
		UsrTypeID: 0,
	})
	require.Nil(t, err)
	require.Equal(t, 1, len(banners))
	require.Equal(t, banner1Id, banners[0].ID)
	require.Equal(t, slot1Id, banners[0].SlotID)
	require.Equal(t, int64(0), banners[0].ShowCnt)
	require.Equal(t, int64(0), banners[0].ClickCnt)

	err = app.cr.BannerCreate(ctx, entities.BannerCreatePars{
		ID:     banner2Id,
		SlotID: slot2Id,
		Note:   bannerNote,
	})
	require.Nil(t, err)

	banners, err = app.db.BannerList(ctx, entities.BannerListPars{
		SlotID:    slot1Id,
		UsrTypeID: 0,
	})
	require.Nil(t, err)
	require.Equal(t, 1, len(banners))
	require.Equal(t, banner1Id, banners[0].ID)
	require.Equal(t, slot1Id, banners[0].SlotID)
	require.Equal(t, int64(0), banners[0].ShowCnt)
	require.Equal(t, int64(0), banners[0].ClickCnt)

	banners, err = app.db.BannerList(ctx, entities.BannerListPars{
		SlotID:    slot2Id,
		UsrTypeID: 0,
	})
	require.Nil(t, err)
	require.Equal(t, 1, len(banners))
	require.Equal(t, banner2Id, banners[0].ID)
	require.Equal(t, slot2Id, banners[0].SlotID)
	require.Equal(t, int64(0), banners[0].ShowCnt)
	require.Equal(t, int64(0), banners[0].ClickCnt)
}

func TestBannerSelect(t *testing.T) {
	var err error

	cleanCtx()

	ctx := context.Background()

	bIds := []int64{1, 2, 3}

	var slotID int64 = 1

	var usrTypeID int64 = 1

	for _, bID := range bIds {
		err = app.cr.BannerCreate(ctx, entities.BannerCreatePars{
			ID:     bID,
			SlotID: slotID,
			Note:   bannerNote,
		})
		require.Nil(t, err)
	}

	banners, err := selectBannerInLoop(ctx, slotID, usrTypeID, 90)
	require.Nil(t, err)
	require.Equal(t, 3, len(banners))
	require.Equal(t, banners[bIds[0]].ShowCnt, banners[bIds[1]].ShowCnt)
	require.Equal(t, banners[bIds[1]].ShowCnt, banners[bIds[2]].ShowCnt)

	err = app.cr.BannerAddClick(ctx, entities.BannerStatIncPars{
		ID:        bIds[0],
		SlotID:    slotID,
		UsrTypeID: usrTypeID,
		Value:     2,
	})
	require.Nil(t, err)

	err = app.cr.BannerAddClick(ctx, entities.BannerStatIncPars{
		ID:        bIds[1],
		SlotID:    slotID,
		UsrTypeID: usrTypeID,
	})
	require.Nil(t, err)

	banners, err = selectBannerInLoop(ctx, slotID, usrTypeID, 90)
	require.Nil(t, err)
	require.Equal(t, 3, len(banners))
	require.True(t, banners[bIds[0]].ShowCnt > banners[bIds[1]].ShowCnt)
	require.True(t, banners[bIds[1]].ShowCnt > banners[bIds[2]].ShowCnt)
}

func selectBannerInLoop(ctx context.Context, slotID, usrTypeID int64, n int) (map[int64]*entities.Banner, error) {
	var err error

	for i := 0; i < n; i++ {
		_, err = app.cr.BannerSelectID(ctx, entities.BannerSelectPars{
			SlotID:    slotID,
			UsrTypeID: usrTypeID,
		})
		if err != nil {
			return nil, err
		}
	}

	var banners []*entities.Banner

	banners, err = app.db.BannerList(ctx, entities.BannerListPars{
		SlotID:    slotID,
		UsrTypeID: usrTypeID,
	})

	if err != nil {
		return nil, err
	}

	res := map[int64]*entities.Banner{}

	for _, banner := range banners {
		res[banner.ID] = banner
	}

	return res, nil
}

func TestBannerEvent(t *testing.T) {
	var err error

	cleanCtx()

	ctx := context.Background()

	var bannerID int64 = 1

	var selectedBannerID int64

	var slotID int64 = 1

	var usrTypeID int64 = 1

	err = app.cr.BannerCreate(ctx, entities.BannerCreatePars{
		ID:     bannerID,
		SlotID: slotID,
		Note:   bannerNote,
	})
	require.Nil(t, err)

	selectedBannerID, err = app.cr.BannerSelectID(ctx, entities.BannerSelectPars{
		SlotID:    slotID,
		UsrTypeID: usrTypeID,
	})
	require.Nil(t, err)
	require.Equal(t, bannerID, selectedBannerID)

	events := app.mq.PullAll()
	require.Equal(t, 1, len(events))
	require.Equal(t, constant.BannerEventTypeShow, events[0].Type)
	require.Equal(t, selectedBannerID, events[0].BannerID)
	require.Equal(t, slotID, events[0].SlotID)
	require.Equal(t, usrTypeID, events[0].UsrTypeID)

	err = app.cr.BannerAddClick(ctx, entities.BannerStatIncPars{
		ID:        bannerID,
		SlotID:    slotID,
		UsrTypeID: usrTypeID,
		Value:     2,
	})
	require.Nil(t, err)

	events = app.mq.PullAll()
	require.Equal(t, 1, len(events))
	require.Equal(t, constant.BannerEventTypeClick, events[0].Type)
	require.Equal(t, selectedBannerID, events[0].BannerID)
	require.Equal(t, slotID, events[0].SlotID)
	require.Equal(t, usrTypeID, events[0].UsrTypeID)
}
