package cmd

import (
	"context"
	"github.com/rendau/barot/internal/adapters/db/pg"
	"github.com/rendau/barot/internal/adapters/http_api"
	"github.com/rendau/barot/internal/adapters/logger/zap"
	"github.com/rendau/barot/internal/adapters/mq/rmq"
	"github.com/rendau/barot/internal/domain/core"
	"github.com/spf13/viper"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"
)

// Execute - executes root command
func Execute() {
	var err error

	loadConf()

	lg, err := zap.NewSt(
		viper.GetString("log_level"),
		true,
		true,
	)
	if err != nil {
		log.Fatal(err)
	}
	defer lg.Sync()

	db, err := pg.NewSt(
		viper.GetString("pg_dsn"),
		lg,
	)
	if err != nil {
		lg.Fatal(err)
	}

	mq, err := rmq.NewSt(
		viper.GetString("rmq_dsn"),
		lg,
	)
	if err != nil {
		lg.Fatal(err)
	}

	cr := core.NewSt(
		lg,
		db,
		mq,
	)

	httpApi := http_api.CreateApi(
		lg,
		viper.GetString("http_listen"),
		cr,
	)

	lg.Infow("Starting", "http_listen", viper.GetString("http_listen"))

	httpApi.Start()

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt, syscall.SIGTERM, syscall.SIGINT)

	var exitCode int

	select {
	case <-stop:
	case <-httpApi.Wait():
		exitCode = 1
	}

	lg.Infow("Shutting down...")

	ctx, _ := context.WithTimeout(context.Background(), 20*time.Second)
	err = httpApi.Shutdown(ctx)
	if err != nil {
		lg.Errorw("Fail to shutdown http-api", err)
		exitCode = 1
	}

	os.Exit(exitCode)
}

func loadConf() {
	confFilePath := os.Getenv("CONF_PATH")
	if confFilePath != "" {
		viper.SetConfigFile(confFilePath)
		err := viper.ReadInConfig()
		if err != nil {
			log.Fatal(err)
		}
	}

	// env vars are in priority
	viper.AutomaticEnv()
}
