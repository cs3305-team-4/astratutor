package routes

import (
	"encoding/json"
	"fmt"
	"net/http"
	"regexp"
	"strings"

	log "github.com/sirupsen/logrus"
)

type returnDetail struct {
	Msg string `json:"msg"`
}
type returnError struct {
	Detail returnDetail `json:"detail"`
}

func httpError(w http.ResponseWriter, r *http.Request, err error, code int) {
	log.WithContext(r.Context()).WithError(err).Error("Error parsing body")
	w.WriteHeader(code)

	err = customErrors(err)
	if err == nil {
		return
	}
	outErr := returnError{
		Detail: returnDetail{
			Msg: err.Error(),
		},
	}

	if err := json.NewEncoder(w).Encode(&outErr); err != nil {
		log.WithContext(r.Context()).WithError(err).Error("Could not return error to user")
		return
	}
}

const (
	sqlDuplicate = "(SQLSTATE 23505)"
)

func customErrors(err error) error {
	if err == nil {
		return err
	}
	switch {
	case strings.Contains(err.Error(), sqlDuplicate):
		re, e := regexp.Compile("_([a-z]+)_key")
		if e != nil {
			log.WithError(err).Error("Error parsing body")
			return err
		}
		match := re.FindStringSubmatch(err.Error())
		if len(match) > 1 {
			return fmt.Errorf("%s already exists", match[1])
		}
	}
	return err
}
