package service

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gorilla/mux"
	_ "github.com/lib/pq"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestUsername(t *testing.T) {
	b, err := NewBDayService()
	require.NoError(t, err)

	err = b.ConnectDB()
	require.NoError(t, err)

	testMux := mux.NewRouter()
	testMux.HandleFunc("/hello/{username:[a-zA-Z]+}", b.username)

	srv := httptest.NewServer(testMux)
	defer srv.Close()

	testsPUT := []struct {
		description  string
		username     string
		body         *BDay
		codeExpected int
	}{
		{"test PUT user", "someuser", &BDay{DateOfBirth: "2012-01-12"}, http.StatusNoContent},
		{"test PUT user again", "someuser", &BDay{DateOfBirth: "2013-01-12"}, http.StatusNoContent},
		{"test PUT user bad date", "otheruser", &BDay{DateOfBirth: "2213-01-12"}, http.StatusBadRequest},
		{"test PUT user too long", "zxcvbnmasdfghjklqwerty", &BDay{DateOfBirth: "2013-01-12"}, http.StatusBadRequest},
	}
	for _, test := range testsPUT {
		reqBody, err := json.Marshal(test.body)
		require.NoError(t, err)
		req, err := http.NewRequest(http.MethodPut, srv.URL+"/hello/"+test.username, bytes.NewBuffer(reqBody))
		client := &http.Client{}
		resp, err := client.Do(req)
		require.NoError(t, err)
		assert.Equal(t, test.codeExpected, resp.StatusCode, test.description+" failed")
		defer resp.Body.Close()
	}

	testsGET := []struct {
		description  string
		username     string
		codeExpected int
	}{
		{"test GET user", "someuser", http.StatusOK},
		{"test GET user again", "someuser", http.StatusOK},
		{"test GET user nonexisting", "xxxxxotheruser", http.StatusNotFound},
		{"test GET user too long", "zxcvbnmasdfghjklqwerty", http.StatusBadRequest},
	}
	for _, test := range testsGET {

		client := &http.Client{}
		resp, err := client.Get(srv.URL + "/hello/" + test.username)
		require.NoError(t, err)
		assert.Equal(t, test.codeExpected, resp.StatusCode, test.description+" failed")
		defer resp.Body.Close()
	}
}

func TestProcessGETUsername(t *testing.T) {
	b, err := NewBDayService()
	require.NoError(t, err)

	err = b.ConnectDB()
	require.NoError(t, err)

	err = b.saveUsernameBDay("kot", "2012-12-12")
	require.NoError(t, err)
	testTime := time.Now()
	err = b.saveUsernameBDay("zenek", fmt.Sprintf("%d-%d-%d", testTime.Year(), testTime.Month(), testTime.Day()))
	require.NoError(t, err)

	tests := []struct {
		description      string
		username         string
		expectedHttpCode int
	}{
		{"no username in db", "michal", http.StatusNotFound},
		{"username in db", "kot", http.StatusOK},
		{"bday username", "zenek", http.StatusOK},
		{"too long username", "qwertyuiopasdfghjklzzx", http.StatusBadRequest},
	}
	for _, test := range tests {
		msg, code := b.processGETUsername(test.username)
		assert.Equal(t, test.expectedHttpCode, code, test.description+" failed")
		if msg != nil {
			bday := &BDay{}
			err = json.Unmarshal(msg, bday)
			require.NoError(t, err, test.description+" failed")
		}
	}

	err = b.db.Close()
	require.NoError(t, err)
	_, code := b.processGETUsername("kot")
	assert.Equal(t, http.StatusNotFound, code, "broken db connection failed")
}

func TestProcessPUTUsername(t *testing.T) {
	b, err := NewBDayService()
	require.NoError(t, err)

	err = b.ConnectDB()
	require.NoError(t, err)

	err = b.saveUsernameBDay("kot", "2012-12-12")
	require.NoError(t, err)

	tests := []struct {
		description      string
		username         string
		body             string
		expectedHttpCode int
	}{
		{"no username in db", "michal", "{ \"dateOfBirth\": \"2012-12-12\" }", http.StatusNoContent},
		{"username in db", "michal", "{ \"dateOfBirth\": \"2013-12-12\" }", http.StatusNoContent},
		{"username bad date", "zzzz", "{ \"dateOfBirth\": \"2013-13-12\" }", http.StatusBadRequest},
		{"username bad date", "xxxx", "{ \"dateOfBirth\": \"2213-01-12\" }", http.StatusBadRequest},
		{"username no data", "yyyy", "", http.StatusBadRequest},
		{"username bad json", "yyyy", "{ \"dateOfBirth\": \"2213-13-12 }", http.StatusBadRequest},
	}
	for _, test := range tests {
		buffer := new(bytes.Buffer)
		_, err := buffer.Write([]byte(test.body))
		require.NoError(t, err)
		code := b.processPUTUsername(test.username, buffer)
		assert.Equal(t, test.expectedHttpCode, code, test.description+" failed")
	}

	err = b.db.Close()
	require.NoError(t, err)
	buffer := new(bytes.Buffer)
	_, err = buffer.Write([]byte(tests[0].body))
	require.NoError(t, err)
	code := b.processPUTUsername(tests[0].username, buffer)
	assert.Equal(t, http.StatusInternalServerError, code, "broken db connection failed")
}
