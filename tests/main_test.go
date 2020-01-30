package main

import (
	"github.com/rendau/my-otus/task8/integration_tests/internal"
	"github.com/spf13/viper"
	"os"
	"testing"
)

const (
	outFormat       = "progress"
	featuresDirPath = "gherkin/features"
)

func TestMain(m *testing.M) {
	viper.AutomaticEnv()

	sc := m.Run()

	os.Exit(sc)
}

func TestIntegration(t *testing.T) {
	tests := internal.NewTests(
		viper.GetString("api_url"),
	)

	status := tests.Run(
		outFormat,
		featuresDirPath,
	)
	if status != 0 {
		t.Fail()
	}
}
