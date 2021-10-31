package service

import (
	"database/sql"
	"errors"
	"fmt"
	"time"

	log "github.com/sirupsen/logrus"
)

// saveUsernameBDay saves or update username and birth date in database
func (b *BDayService) saveUsernameBDay(username, bday string) error {
	log.SetLevel(log.DebugLevel)
	var id int64
	var oldBday time.Time
	var stmt *sql.Stmt
	t, err := time.Parse(time.RFC3339, fmt.Sprintf("%sT00:00:00.000Z", bday))
	if err != nil {
		log.WithField("context", "saveUsernameBDay").WithError(err).Debugf("time parse failed")
		return err
	}
	result, err := b.db.Query("SELECT id, bday FROM bday WHERE username = $1;", username)
	if err != nil {
		log.WithField("context", "saveUsernameBDay").WithError(err).Debugf("check existing user query failed")
		return err
	}

	if result.Next() {
		// existing username
		err = result.Scan(&id, &oldBday)
		if err != nil {
			log.WithField("context", "saveUsernameBDay").WithError(err).Debugf("check existing user scan failed")
			return err
		}
		// don't update if the bday is the same day
		if t.Equal(oldBday) {
			return nil
		}
		stmt, err = b.db.Prepare("UPDATE bday SET bday = $1 WHERE id = $2;")
		if err != nil {
			log.WithField("context", "saveUsernameBDay").WithError(err).Debugf("update existing user perare failed")
			return err
		}
		_, err = stmt.Exec(bday, id)
		if err != nil {
			log.WithField("context", "saveUsernameBDay").WithError(err).Debugf("update existing user exec failed")
			return err
		}
	} else {
		// new username
		stmt, err = b.db.Prepare("INSERT INTO bday(username, bday) VALUES($1,$2);")
		if err != nil {
			log.WithField("context", "saveUsernameBDay").WithError(err).Debugf("insert new user perare failed")
			return err
		}
		_, err = stmt.Exec(username, t)
		if err != nil {
			log.WithField("context", "saveUsernameBDay").WithError(err).Debugf("insert new user exec failed")
			return err
		}
	}
	return nil
}

// getUsernameBDay reads birth day from database
func (b *BDayService) getUsernameBDay(username string) (string, error) {

	var bday string
	result, err := b.db.Query("SELECT bday FROM bday WHERE username = $1;", username)
	if err != nil {
		log.WithField("context", "getUsernameBDay").WithError(err).Debugf("get user query failed")
		return "", err
	}

	if result.Next() {
		err = result.Scan(&bday)
		if err != nil {
			log.WithField("context", "getUsernameBDay").WithError(err).Debugf("get user scan failed")
			return "", err
		}
		return bday, nil
	}
	return "", errors.New("no user")
}
