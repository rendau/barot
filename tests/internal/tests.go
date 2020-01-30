package internal

import (
	"context"
	"github.com/DATA-DOG/godog"
	"log"
	"net/http"
	"time"
)

const (
	servicesWaitTimeout = 30 * time.Second
	slot1Id             = 1
	usrType1Id          = 1
)

// Tests - is type for tests
type Tests struct {
	apiUrl string

	responseStatusCode int
	showCounts         map[int64]int64
}

func NewTests(apiUrl string) *Tests {
	return &Tests{
		apiUrl:     apiUrl,
		showCounts: map[int64]int64{},
	}
}

func (t *Tests) Run(outFormat, featuresPath string) int {
	ctx, _ := context.WithTimeout(context.Background(), servicesWaitTimeout)
	err := t.waitServices(ctx)
	if err != nil {
		log.Println("Fail to wait services")
		return 1
	}

	return godog.RunWithOptions(
		"integration",
		t.FeatureContext,
		godog.Options{
			Format: outFormat,
			Paths:  []string{featuresPath},
		},
	)
}

func (t *Tests) waitServices(ctx context.Context) error {
	var err error
	var req *http.Request
	var rep *http.Response

	client := http.Client{Timeout: 1 * time.Second}

	req, err = http.NewRequestWithContext(ctx, "GET", t.apiUrl+"/ready", nil)
	if err != nil {
		return err
	}

	for {
		rep, err = client.Do(req)
		if err == nil {
			_ = rep.Body.Close()
			break
		}
		if ctx.Err() != nil {
			err = ctx.Err()
			break
		}
		time.Sleep(time.Second)
	}

	return err
}
