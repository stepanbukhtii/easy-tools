package elog

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
	"resty.dev/v3"
)

func TestResty(t *testing.T) {
	slog.SetDefault(slog.New(slog.NewJSONHandler(os.Stdout, nil)))

	restyLogger := &RestyLogger{ServiceName: "target_name"}

	client := resty.New().EnableTrace()
	defer client.Close()

	resp, err := client.R().SetContext(context.Background()).Get("https://httpbin.org/get")
	require.NoError(t, err)

	restyLogger.Info(resp, "request 1")

	testServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		time.Sleep(2 * time.Second)
		fmt.Fprint(w, "Hello")
	}))
	defer testServer.Close()

	resp, err = client.R().SetTimeout(time.Second).SetContext(context.Background()).Get(testServer.URL)
	require.Error(t, err)
	restyLogger.Error(resp, err, "request 2")

}
