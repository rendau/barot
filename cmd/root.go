package cmd

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/rendau/barot/internal/adapters/db/pg"
	"github.com/rendau/barot/internal/adapters/httpapi"
	"github.com/rendau/barot/internal/adapters/logger/zap"
	"github.com/rendau/barot/internal/adapters/mq/rmq"
	"github.com/rendau/barot/internal/domain/core"
	"github.com/spf13/viper"
)

// Execute - executes root command
// nolint:funlen
func Execute() {
	var err error

	loadConf()

	lg, err := zap.NewLogger(viper.GetString("log_level"), true, true)
	if err != nil {
		log.Fatal(err)
	}

	defer lg.Sync()

	db, err := pg.NewPostgresDB(viper.GetString("pg_dsn"), lg)
	if err != nil {
		lg.Fatal(err)
	}

	mq, err := rmq.NewRabbitMQ(viper.GetString("rmq_dsn"), lg)
	if err != nil {
		lg.Fatal(err)
	}

	cr := core.NewSt(lg, db, mq)

	httpAPIInst := httpapi.CreateAPI(lg, viper.GetString("http_listen"), cr)

	lg.Infow("Starting", "http_listen", viper.GetString("http_listen"))

	httpAPIInst.Start()

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt, syscall.SIGTERM, syscall.SIGINT)

	var exitCode int

	select {
	case <-stop:
	case <-httpAPIInst.Wait():
		exitCode = 1
	}

	lg.Infow("Shutting down...")

	ctx, ctxCancel := context.WithTimeout(context.Background(), 20*time.Second)
	defer ctxCancel()

	err = httpAPIInst.Shutdown(ctx)
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
