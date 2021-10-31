package service

import (
	"errors"
	"fmt"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestCheckDate(t *testing.T) {
	b, err := NewBDayService()
	require.NoError(t, err)

	tests := []struct {
		description string
		dateShift   int
		isBefore    bool
	}{
		{"test date before", -1, true},
		{"test date today", 0, false},
		{"test date after", 1, false},
	}

	for _, test := range tests {
		testTime := time.Now().AddDate(0, 0, test.dateShift)
		dateString := fmt.Sprintf("%d-%d-%d", testTime.Year(), testTime.Month(), testTime.Day())
		err := b.checkDate(dateString)
		if test.isBefore {
			require.NoError(t, err)
		} else {
			require.Error(t, err)
		}
	}
	err = b.checkDate("dasdadada")
	require.Error(t, err)
}

func TestCountDays(t *testing.T) {
	b, err := NewBDayService()
	require.NoError(t, err)

	tests := []struct {
		description string
		yearShift   int
		monthShift  int
		dayShift    int
		isBDay      bool
		expectedErr error
	}{
		{"test birthday", -10, 0, 0, true, nil},
		{"test 1 month, 1 day", 0, -1, -1, false, nil},
		{"test 1 year, 1 month, 1 day", -1, -1, -1, false, nil},
		{"test date future", 0, 0, 1, false, errors.New("bday is from future")},
	}

	for _, test := range tests {
		testTime := time.Now().AddDate(test.yearShift, test.monthShift, test.dayShift)
		dateString := fmt.Sprintf("%d-%02d-%02dT00:00:00Z", testTime.Year(), testTime.Month(), testTime.Day())
		i, err := b.countDays(dateString)
		assert.Equal(t, test.expectedErr, err)
		if test.isBDay {
			assert.Equal(t, 0, i)
		} else {
			assert.NotEqual(t, 0, i)
		}
	}

	tests2 := []struct {
		description  string
		dayShift     int
		expectedDays int
		expectedErr  error
	}{
		{"test birthday", 0, 0, nil},
		{"test 10 days after", -10, -10 + time.Date(time.Now().Year(), 12, 31, 23, 59, 59, 0, time.UTC).YearDay(), nil},
		{"test 11 days before ", 11, 11, nil},
	}
	for _, test := range tests2 {
		testTime := time.Now().AddDate(-1, 0, test.dayShift)
		dateString := fmt.Sprintf("%d-%02d-%02dT00:00:00Z", testTime.Year(), testTime.Month(), testTime.Day())
		i, err := b.countDays(dateString)
		assert.Equal(t, test.expectedErr, err)
		assert.Equal(t, test.expectedDays, i)
	}
}
