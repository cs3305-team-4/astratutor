package routes

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"regexp"
	"strings"

	"github.com/cs3305-team-4/api/pkg/services"
	"github.com/go-playground/validator/v10"
	log "github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

type returnDetail struct {
	Msg string `json:"msg"`
}
type returnError struct {
	Detail returnDetail `json:"detail"`
}

func restError(w http.ResponseWriter, r *http.Request, err error, code int) {
	log.WithContext(r.Context()).WithError(err).Error("Error parsing body")
	err, code = customErrors(err, code)
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
	switch {
	case errors.Is(in, services.AccountErrorProfileExists):
		codeOut = http.StatusBadRequest
	case errors.Is(in, services.AccountErrorProfileDoesNotExists):
		codeOut = http.StatusNotFound
	case errors.Is(in, services.AccountErrorAccountDoesNotExist):
		codeOut = http.StatusNotFound
	case errors.Is(in, services.AccountErrorEntryDoesNotExists):
		codeOut = http.StatusNotFound
	case errors.Is(in, gorm.ErrRecordNotFound):
		codeOut = http.StatusNotFound
		out = errors.New("No record matching provided ID found.")
	case errors.Is(in, io.EOF):
		fallthrough
	case errors.Is(in, &json.SyntaxError{}):
		codeOut = http.StatusBadRequest
		out = errors.New("Validation failed for body: invalid format.")

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

	// Validation error
	if errs, ok := in.(validator.ValidationErrors); ok {
		outs := "Validation failed for "
		for _, err := range errs {
			field := err.Field()
			kind := err.ActualTag()
			outs += fmt.Sprintf("('%s' on field '%s') ", kind, field)
		}
		codeOut = http.StatusBadRequest
		out = errors.New(outs)
	}
	return
}
