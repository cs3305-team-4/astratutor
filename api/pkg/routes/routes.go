package routes

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"reflect"
	"regexp"
	"strings"

	"github.com/cs3305-team-4/api/pkg/services"
	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
	"github.com/gorilla/mux"
	"github.com/rs/cors"
)

var validate *validator.Validate

func GetHandler() http.Handler {
	validate = validator.New()
	setCustomValidators()
	services.SetCustomValidators()
	r := mux.NewRouter()

	r.Use(
		loggingMiddleware,
		jsonMiddleware,
	)

	InjectAccountsRoutes(r.PathPrefix("/accounts").Subrouter())
	InjectAuthRoutes(r.PathPrefix("/auth").Subrouter())
	InjectLessonsRoutes(r.PathPrefix("/lessons").Subrouter())
	InjectStudentsRoutes(r.PathPrefix("/students").Subrouter())
	InjectSubjectsRoutes(r.PathPrefix("/subjects").Subrouter())
	InjectTutorsRoutes(r.PathPrefix("/tutors").Subrouter())

	return cors.New(cors.Options{
		AllowedOrigins: []string{"*"},
		AllowedMethods: []string{
			http.MethodGet,
			http.MethodPost,
			http.MethodPut,
			http.MethodPatch,
			http.MethodDelete,
			http.MethodOptions,
			http.MethodHead,
		},

		AllowedHeaders: []string{
			"*",
		},
	}).Handler(r)
}

func setCustomValidators() {
	validate.RegisterValidationCtx("passwd", func(ctx context.Context, v validator.FieldLevel) bool {
		errs, ok := ctx.Value("error").(*map[string]error)
		if !ok {
			return false
		}
		st := v.Field()
		if st.Kind() != reflect.String {
			(*errs)["error"] = errors.New("passwd only validates strings")
			return false
		}
		val := st.String()
		if len(val) < 8 {
			(*errs)["error"] = errors.New("Password must have at least 8 characters.")
			return false
		}
		if strings.ToLower(val) == val {
			(*errs)["error"] = errors.New("Password must have at least one upper case letter.")
			return false
		}
		if strings.ToUpper(val) == val {
			(*errs)["error"] = errors.New("Password must have at least one lower case letter.")
			return false
		}
		numRe := regexp.MustCompile(`[0-9]+`)
		if !numRe.MatchString(val) {
			(*errs)["error"] = errors.New("Password must have at least one number.")
			return false
		}
		return true
	})
}

// ParseBody inplace. Returns false if error has been written.
func ParseBody(w http.ResponseWriter, r *http.Request, i interface{}) bool {
	if err := json.NewDecoder(r.Body).Decode(i); err != nil {
		restError(w, r, err, http.StatusBadRequest)
		return false
	}
	errs := &map[string]error{}
	ctx := context.WithValue(context.Background(), "error", errs)
	if err := validate.StructCtx(ctx, i); err != nil {
		if _, ok := (*errs)["error"]; ok {
			err = (*errs)["error"]
		}
		restError(w, r, err, http.StatusBadRequest)
		return false
	}
	return true
}

// getUUID can parse a UUID from the router variables
// if param is nil, the default variable used "uuid"
func getUUID(r *http.Request, param string) (uuid.UUID, error) {
	para := "uuid"
	if param != "" {
		para = param
	}

	vars := mux.Vars(r)
	val, ok := vars[para]
	if !ok {
		return uuid.UUID{}, errors.New("no uuid found")
	}
	return uuid.Parse(val)
}

// UpdateDTO used for single field update route posts.
type UpdateDTO struct {
	Value string `json:"value"`
}

// ParseUpdateString will parse
func ParseUpdateString(w http.ResponseWriter, r *http.Request) string {
	update := &UpdateDTO{}
	if !ParseBody(w, r, update) {
		return ""
	}
	return update.Value
}
