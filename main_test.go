package main

import (
	"github.com/rendau/barot/internal/adapters/db/pg"
	"github.com/rendau/barot/internal/adapters/logger/zap"
	"github.com/rendau/barot/internal/constant"
	"github.com/rendau/barot/internal/domain/core"
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
	)
	if err != nil {
		app.lg.Fatal(err)
	}

	app.cr = core.NewSt()

	ec := m.Run()
	os.Exit(ec)
}

func TestMabCalc(t *testing.T) {
	v := app.cr.MabCalc(0, 0, 100)
	require.Equal(t, constant.MabCalcInitValue, v)
}
