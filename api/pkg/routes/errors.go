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
	err, code = customErrors(err, code)
	if err == nil {
		return
	}
	w.WriteHeader(code)
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

func customErrors(in error, code int) (out error, codeOut int) {
	codeOut = code
	out = in
	if in == nil {
		return
	}
	switch {
	case strings.Contains(in.Error(), sqlDuplicate):
		re, e := regexp.Compile("_([a-z]+)_key")
		if e != nil {
			log.WithError(e).Error("Error parsing body")
			return
		}
		match := re.FindStringSubmatch(in.Error())
		if len(match) > 1 {
			out = fmt.Errorf("%s already exists", match[1])
		}
		codeOut = http.StatusBadRequest
	}
	return
}
