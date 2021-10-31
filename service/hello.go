package service

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/gorilla/mux"
	log "github.com/sirupsen/logrus"
)

var usernameMaxLen = 20

type Message struct {
	Message string `json:"message"`
}

type BDay struct {
	DateOfBirth string `json:"dateOfBirth"`
}

// username process request to /hello/username endpoint
func (b *BDayService) username(w http.ResponseWriter, req *http.Request) {
	vars := mux.Vars(req)
	bDayReqProcessed.Inc()

	switch req.Method {
	case "GET":
		msg, code := b.processGETUsername(vars["username"])
		w.WriteHeader(code)
		if msg != nil {
			fmt.Fprintf(w, "%s\n", msg)
		}
	case "PUT":
		code := b.processPUTUsername(vars["username"], req.Body)
		w.WriteHeader(code)

	default:
		// For an Internet facing app I would prefer to send here http.StatusNotFound,
		// to make scaning the app harder for someone potentially looking for an atack vector.
		// For internal app http.StatusNotImplemented is the way to go.
		w.WriteHeader(http.StatusNotImplemented)
	}
}

// processGETUsername process GET requests for /hello/username
func (b *BDayService) processGETUsername(username string) ([]byte, int) {
	if len(username) > usernameMaxLen {
		log.WithField("context", "processGETUsername").Debugf("username too long")
		return nil, http.StatusBadRequest
	}
	msg := ""
	bBay, err := b.getUsernameBDay(username)
	if err != nil {
		return nil, http.StatusNotFound
	}
	daysToBday, err := b.countDays(bBay)
	if err != nil {
		return nil, http.StatusNotFound
	}
	if daysToBday == 0 {
		msg = fmt.Sprintf("Hello, %s! Happy birthday!", username)
	} else {
		msg = fmt.Sprintf("Hello, %s! Your birthday is in %d day(s)", username, daysToBday)
	}
	j, err := json.Marshal(&Message{
		Message: msg,
	})
	if err != nil {
		return nil, http.StatusInternalServerError
	}
	return j, http.StatusOK
}

// processPUTUsername process PUT requests for /hello/username
func (b *BDayService) processPUTUsername(username string, reqBody io.Reader) int {
	if len(username) > usernameMaxLen {
		log.WithField("context", "processPUTUsername").Debugf("user name too long")
		return http.StatusBadRequest
	}
	body, err := io.ReadAll(reqBody)
	if err != nil {
		log.WithField("context", "processPUTUsername").WithError(err).Debugf("unable to read body")
		return http.StatusInternalServerError
	}

	bday := &BDay{}
	err = json.Unmarshal(body, bday)
	if err != nil {
		log.WithField("context", "processPUTUsername").WithError(err).Debugf("unable to parse json")
		return http.StatusBadRequest
	}
	err = b.checkDate(bday.DateOfBirth)
	if err != nil {
		return http.StatusBadRequest
	}

	err = b.saveUsernameBDay(username, bday.DateOfBirth)
	if err != nil {
		log.WithField("context", "processPUTUsername").WithError(err).Debugf("unable to save data")
		return http.StatusInternalServerError
	}
	return http.StatusNoContent
}
