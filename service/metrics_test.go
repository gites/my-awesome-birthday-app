package service

import (
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestMetricsEndpoint(t *testing.T) {
	b, err := NewBDayService()
	require.NoError(t, err)

	testMux := http.NewServeMux()
	testMux.HandleFunc("/live", b.liveHandle)
	testMux.HandleFunc("/ready", b.readyHandle)
	testMux.HandleFunc("/version", b.versionHandle)
	testMux.Handle("/metrics", promhttp.Handler())
	srv := httptest.NewServer(testMux)
	defer srv.Close()

	client := &http.Client{}

	tests := []struct {
		description  string
		path         string
		ready        bool
		codeExpected int
		empty        bool
	}{
		{"test metrics", "/metrics", false, http.StatusOK, false},
		{"test version", "/version", false, http.StatusOK, false},
		{"test live", "/live", false, http.StatusOK, false},
		{"test not ready", "/ready", false, http.StatusTooEarly, false},
		{"test ready", "/ready", true, http.StatusOK, false},
	}

	for _, test := range tests {

		b.isReady = test.ready
		resp, err := client.Get(srv.URL + test.path)
		require.NoError(t, err, test.description+" failed")
		defer resp.Body.Close()

		assert.Equal(t, test.codeExpected, resp.StatusCode, test.description+" failed")
		data, err := io.ReadAll(resp.Body)
		require.NoError(t, err, test.description+" failed")
		if test.empty {
			assert.Empty(t, data, test.description+" failed")
		} else {
			assert.NotEmpty(t, data, test.description+" failed")
		}
	}
}

func TestRunMetricsServer(t *testing.T) {
	b, err := NewBDayService()
	require.NoError(t, err)
	err = b.ConnectDB()
	require.NoError(t, err)
	go b.runMetricsServer()
	time.Sleep(1 * time.Second)

	client := &http.Client{}
	resp, err := client.Get(fmt.Sprintf("http://localhost:%s/live", b.MetricsPort))
	require.NoError(t, err)
	defer resp.Body.Close()
	assert.Equal(t, http.StatusOK, resp.StatusCode)
}
