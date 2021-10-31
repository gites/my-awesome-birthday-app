package service

import (
	"errors"
	"testing"

	_ "github.com/lib/pq"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestSaveUsernameBDay(t *testing.T) {
	b, err := NewBDayService()
	require.NoError(t, err)

	err = b.ConnectDB()
	require.NoError(t, err)

	tests := []struct {
		description  string
		username     string
		bday         string
		expectedBDay string
	}{
		{"test username", "ala", "1912-12-12", "1912-12-12T00:00:00Z"},
		{"test username", "ala", "1912-12-12", "1912-12-12T00:00:00Z"},
		{"test username", "ala", "2012-12-12", "2012-12-12T00:00:00Z"},
	}

	for _, test := range tests {
		err := b.saveUsernameBDay(test.username, test.bday)
		require.NoError(t, err)
		actualBDay, err := b.getUsernameBDay(test.username)
		require.NoError(t, err)
		assert.Equal(t, test.expectedBDay, actualBDay)
	}
	b.db.Close()
	err = b.saveUsernameBDay(tests[0].username, tests[0].bday)
	require.Error(t, err)
}

func TestGetUsernameBDay(t *testing.T) {
	b, err := NewBDayService()
	require.NoError(t, err)

	err = b.ConnectDB()
	require.NoError(t, err)
	tests := []struct {
		description  string
		username     string
		expectedBDay string
		expectErr    error
	}{
		{"test get username", "zxxc", "2012-12-12T00:00:00Z", nil},
		{"test get nonexisting username", "dsadadsa", "", errors.New("no user")},
	}

	err = b.saveUsernameBDay("zxxc", "2012-12-12")
	require.NoError(t, err)

	for _, test := range tests {
		actualBDay, err := b.getUsernameBDay(test.username)
		assert.Equal(t, test.expectErr, err)
		assert.Equal(t, test.expectedBDay, actualBDay)
	}

	b.db.Close()
	_, err = b.getUsernameBDay(tests[0].username)
	require.Error(t, err)
}
