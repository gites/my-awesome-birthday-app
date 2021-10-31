package service

import (
	"log"
	"os"
	"testing"

	"github.com/golang-migrate/migrate"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestMain(m *testing.M) {

	b, err := NewBDayService()
	if err != nil {
		log.Fatal(err)
	}
	err = RunMigrations(b)
	if err != nil {
		log.Fatal(err)
	}
	os.Exit(m.Run())
}

func TestMigrations(t *testing.T) {

	b, err := NewBDayService()
	require.NoError(t, err)
	err = RunMigrations(b)
	assert.Equal(t, migrate.ErrNoChange, err)
	err = RunMigrations(b)
	assert.Equal(t, migrate.ErrNoChange, err)

}

func TestNewBDayService(t *testing.T) {
	b, err := NewBDayService()
	require.NoError(t, err)
	require.NotEmpty(t, b)
}

func TestConnectDB(t *testing.T) {
	b, err := NewBDayService()
	require.NoError(t, err)
	err = b.ConnectDB()
	require.NoError(t, err)
	require.NotEmpty(t, b.db)
	defer b.db.Close()
}
