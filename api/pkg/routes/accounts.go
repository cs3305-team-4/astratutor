package routes

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"

	"github.com/cs3305-team-4/api/pkg/services"
	"github.com/gorilla/mux"
)

func InjectAccountsRoutes(subrouter *mux.Router) {
	// /
	subrouter.HandleFunc("", handleAccountsPost).Methods("POST")

	// /{uuid}/
	accountResource := subrouter.PathPrefix("/{uuid}").Subrouter()

	// Only allow users to access routes relevant to their own account
	accountResource.Use(authMiddleware(func(w http.ResponseWriter, r *http.Request, ac *AuthContext) error {
		id, err := getUUID(r, "uuid")
		if err != nil {
			return err
		}

		if ac.Account.ID != id {
			return errors.New("cannot operate on a resource you do not own")
		}

		return nil
	}))

	// /{uuid}/lessons
	accountResource.HandleFunc("/lessons", handleAccountsLessonsGet).Methods("GET")
}

// Account DTO.
type AccountDTO struct {
	ID           string `json:"id" validate:"len=0"`
	Email        string `json:"email" validate:"nonzero"`
	Type         string `json:"type"`
	Password     string `json:"password,omitempty" validate:"passwd"`
	ParentsEmail string `json:"parents_email,omitempty"`
}

func handleAccountsPost(w http.ResponseWriter, r *http.Request) {
	account := &AccountDTO{}
	if !ParseBody(w, r, account) {
		return
	}
	accountType, err := services.ToAccountType(account.Type)
	if err != nil {
		restError(w, r, err, http.StatusBadRequest)
		return
	}
	hash, err := services.NewPasswordHash(account.Password)
	if err != nil {
		restError(w, r, err, http.StatusInternalServerError)
		return
	}
	serviceAccount := &services.Account{
		Email:        account.Email,
		Type:         accountType,
		PasswordHash: *hash,
	}
	if err = services.CreateAccount(serviceAccount); err != nil {
		restError(w, r, err, http.StatusInternalServerError)
		return
	}
	outAccount := &AccountDTO{
		ID:    serviceAccount.ID.String(),
		Email: serviceAccount.Email,
		Type:  string(serviceAccount.Type),
	}
	if err = json.NewEncoder(w).Encode(outAccount); err != nil {
		restError(w, r, err, http.StatusInternalServerError)
		return
	}
}

func handleAccountsVerify(w http.ResponseWriter, r *http.Request) {
	id, err := getUUID(r, "uuid")
	if err != nil {
		restError(w, r, err, http.StatusBadRequest)
	}
	fmt.Println(id)
}

func handleAccountsLessonsGet(w http.ResponseWriter, r *http.Request) {
	id, err := getUUID(r, "uuid")
	if err != nil {
		restError(w, r, err, http.StatusBadRequest)
		return
	}

	lessons, err := services.ReadLessonsByAccountID(id)
	if err != nil {
		restError(w, r, err, http.StatusBadRequest)
		return
	}

	dtoLessons := dtoFromLessons(lessons)
	if err = json.NewEncoder(w).Encode(dtoLessons); err != nil {
		restError(w, r, err, http.StatusInternalServerError)
		return
	}
}
