package main

import (
	"context"
	"github.com/rendau/barot/internal/adapters/db/pg"
	"github.com/rendau/barot/internal/adapters/logger/zap"
	"github.com/rendau/barot/internal/constant"
	"github.com/rendau/barot/internal/domain/core"
	"github.com/rendau/barot/internal/domain/entities"
	"github.com/spf13/viper"
	"github.com/stretchr/testify/require"
	"log"
	"os"
	"testing"
)

const confFilePath = "conf_test.yml"

var (
	app = struct {
		lg *zap.St
		db *pg.St
		cr *core.St
	}{}
)

func TestMain(m *testing.M) {
	var err error

	viper.SetConfigFile(confFilePath)
	_ = viper.ReadInConfig()
	viper.AutomaticEnv()

	app.lg, err = zap.NewSt(
		viper.GetString("log_level"),
		true,
		true,
	)
	if err != nil {
		log.Fatal(err)
	}
	defer app.lg.Sync()

	app.db, err = pg.NewSt(
		viper.GetString("pg_dsn"),
		app.lg,
	)
	if err != nil {
		app.lg.Fatal(err)
	}

	app.cr = core.NewSt(
		app.lg,
		app.db,
	)

	if err != nil {
		app.lg.Fatal(err)
	}

	ec := m.Run()
	os.Exit(ec)
}

func cleanDb() {
	_, err := app.db.Db.Exec(`
		truncate stat restart identity cascade
	`)
	if err != nil {
		app.lg.Fatal(err)
	}

	_, err = app.db.Db.Exec(`
		truncate banner restart identity cascade
	`)
	if err != nil {
		app.lg.Fatal(err)
	}
}

func TestMabCalc(t *testing.T) {
	v := app.cr.MabCalc(0, 0, 100)
	require.Equal(t, constant.MabCalcInitValue, v)
}

func TestBannerCreation(t *testing.T) {
	var err error

	cleanDb()

	ctx := context.Background()

	var banner1Id int64 = 1
	var banner2Id int64 = 2

	var slot1Id int64 = 1
	var slot2Id int64 = 2

	note := "some banner note"

	err = app.cr.BannerCreate(ctx, entities.BannerCreatePars{
		ID:     banner1Id,
		SlotID: slot1Id,
		Note:   note,
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
		Note:   note,
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

	cleanDb()

	ctx := context.Background()

	bIds := []int64{1, 2, 3}

	var slotId int64 = 1
	var usrTypeId int64 = 1

	note := "some banner note"

	for _, bId := range bIds {
		err = app.cr.BannerCreate(ctx, entities.BannerCreatePars{
			ID:     bId,
			SlotID: slotId,
			Note:   note,
		})
		require.Nil(t, err)
	}

	banners, err := selectBannerInLoop(ctx, slotId, usrTypeId, 90)
	require.Nil(t, err)
	require.Equal(t, 3, len(banners))
	require.Equal(t, banners[bIds[0]].ShowCnt, banners[bIds[1]].ShowCnt)
	require.Equal(t, banners[bIds[1]].ShowCnt, banners[bIds[2]].ShowCnt)

	err = app.cr.BannerAddClick(ctx, entities.BannerStatIncPars{
		ID:        bIds[0],
		SlotID:    slotId,
		UsrTypeID: usrTypeId,
		Value:     2,
	})
	require.Nil(t, err)

	err = app.cr.BannerAddClick(ctx, entities.BannerStatIncPars{
		ID:        bIds[1],
		SlotID:    slotId,
		UsrTypeID: usrTypeId,
	})
	require.Nil(t, err)

	banners, err = selectBannerInLoop(ctx, slotId, usrTypeId, 90)
	require.Nil(t, err)
	require.Equal(t, 3, len(banners))
	require.True(t, banners[bIds[0]].ShowCnt > banners[bIds[1]].ShowCnt)
	require.True(t, banners[bIds[1]].ShowCnt > banners[bIds[2]].ShowCnt)
}

func selectBannerInLoop(ctx context.Context, slotId, usrTypeId int64, n int) (map[int64]*entities.Banner, error) {
	var err error

	for i := 0; i < n; i++ {
		_, err = app.cr.BannerSelectId(ctx, entities.BannerListPars{
			SlotID:    slotId,
			UsrTypeID: usrTypeId,
		})
		if err != nil {
			return nil, err
		}
	}
	var banners []*entities.Banner
	banners, err = app.db.BannerList(ctx, entities.BannerListPars{
		SlotID:    slotId,
		UsrTypeID: usrTypeId,
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
