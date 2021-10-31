package service

import (
	"errors"
	"fmt"
	"time"

	log "github.com/sirupsen/logrus"
)

// countDays counts the days to birthday
func (b *BDayService) countDays(bday string) (int, error) {
	t, err := time.Parse(time.RFC3339, bday)
	if err != nil {
		log.WithField("context", "checkTime").WithError(err).Debugf("time parse failed")
		return 0, err
	}
	if t.After(time.Now()) {
		return -1, errors.New(("bday is from future"))
	}
	var days int
	nowYearDay := time.Now().YearDay()
	bdayYearDay := time.Date(time.Now().Year(), t.Month(), t.Day(), 23, 59, 59, 0, time.UTC).YearDay()
	lastYearDay := time.Date(time.Now().Year(), 12, 31, 23, 59, 59, 0, time.UTC).YearDay()
	days = bdayYearDay - nowYearDay
	if days < 0 {
		days = days + lastYearDay
	}
	return days, nil
}

// checkDate checks if a date is a date before today date
func (b *BDayService) checkDate(bday string) error {
	t, err := time.Parse(time.RFC3339, fmt.Sprintf("%sT00:00:00.000Z", bday))
	if err != nil {
		log.WithField("context", "checkDate").WithError(err).Debugf("time parse failed")
		return err
	}
	// Get time one day ago
	yesterday := time.Now().AddDate(0, 0, -1)
	if !t.Before(yesterday) {
		log.WithField("context", "checkDate").WithError(err).Debugf("time not before today")
		return errors.New("time not before today")
	}
	return nil
}
