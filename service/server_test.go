package service

import (
	"fmt"
	"net/http"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestRunServer(t *testing.T) {
	b, err := NewBDayService()
	require.NoError(t, err)
	err = b.ConnectDB()
	require.NoError(t, err)
	go b.runServer()
	time.Sleep(1 * time.Second)

	client := &http.Client{}
	resp, err := client.Get(fmt.Sprintf("http://localhost:%s/username/xxxx", b.HttpPort))
	require.NoError(t, err)
	defer resp.Body.Close()
	assert.Equal(t, http.StatusNotFound, resp.StatusCode)
}
