package routes

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"

	"github.com/cs3305-team-4/api/pkg/services"
	"github.com/gorilla/mux"
)

func authAccount() func(next http.Handler) http.Handler {
	return authMiddleware(func(w http.ResponseWriter, r *http.Request, ac *AuthContext) error {
		id, err := getUUID(r, "uuid")
		if err != nil {
			return err
		}

		if ac.Account.ID != id {
			return errors.New("cannot operate on a resource you do not own")
		}

		return nil
	})
}

func InjectAccountsRoutes(subrouter *mux.Router) {
	subrouter.HandleFunc("", handleAccountsPost).Methods("POST")

	// Only allow users to access routes relevant to their own account
	accountResource := subrouter.PathPrefix("/{uuid}").Subrouter()
	accountResource.Use(authAccount())
	accountResource.HandleFunc("", handleAccountsGet).Methods("GET")
	accountResource.HandleFunc("/verify", handleAccountsVerify).Methods("POST")
	accountResource.HandleFunc("/email", handleAccountsUpdateEmail).Methods("POST")
	accountResource.HandleFunc("/password", handleAccountsUpdatePassword).Methods("POST")
	accountResource.HandleFunc("/lessons", handleAccountsLessonsGet).Methods("GET")
}

// Account DTO.
type AccountDTO struct {
	ID           string `json:"id" validate:"len=0"`
	Email        string `json:"email" validate:"required,email"`
	Type         string `json:"type" validate:"required"`
	Password     string `json:"password,omitempty" validate:"required,passwd"`
	ParentsEmail string `json:"parents_email,omitempty" validate:"omitempty,email"`
}

func dtoFromAccount(a *services.Account) *AccountDTO {
	return &AccountDTO{
		ID:    a.ID.String(),
		Email: a.Email,
		Type:  string(a.Type),
	}
}

func handleAccountsGet(w http.ResponseWriter, r *http.Request) {
	id, err := getUUID(r, "uuid")
	if err != nil {
		restError(w, r, err, http.StatusBadRequest)
		return
	}
	serviceAccount, err := services.ReadAccountByID(id, nil)
	if err != nil {
		restError(w, r, err, http.StatusBadRequest)
		return
	}
	outAccount := dtoFromAccount(serviceAccount)
	if err = json.NewEncoder(w).Encode(outAccount); err != nil {
		restError(w, r, err, http.StatusInternalServerError)
		return
	}
}

func handleAccountsUpdateEmail(w http.ResponseWriter, r *http.Request) {
	id, err := getUUID(r, "uuid")
	if err != nil {
		restError(w, r, err, http.StatusBadRequest)
		return
	}
	// TODO(ericm): Add email verification.
	field := ParseUpdateString(w, r)
	if err = validateUpdate("Email", field, &AccountDTO{}); err != nil {
		restError(w, r, err, http.StatusBadRequest)
		return
	}
	var account *services.Account
	if account, err = services.UpdateAccountField(id, "email", field); err != nil {
		restError(w, r, err, http.StatusBadRequest)
		return
	}
	outAccount := dtoFromAccount(account)
	if err = json.NewEncoder(w).Encode(outAccount); err != nil {
		restError(w, r, err, http.StatusInternalServerError)
		return
	}
}

// UpdateDTO used for single field update route posts.
type UpdatePasswordDTO struct {
	Value struct {
		OldPassword string `json:"old_password"`
		NewPassword string `json:"new_password"`
	} `json:"value"`
}

func handleAccountsUpdatePassword(w http.ResponseWriter, r *http.Request) {
	id, err := getUUID(r, "uuid")
	if err != nil {
		restError(w, r, err, http.StatusBadRequest)
		return
	}
	update := &UpdatePasswordDTO{}
	if !ParseBody(w, r, update) {
		return
	}
	newPassword := update.Value.NewPassword
	// TODO(ericm): Validate OldPassword hash.
	if err = validateUpdate("Password", newPassword, &AccountDTO{}); err != nil {
		restError(w, r, err, http.StatusBadRequest)
		return
	}
	passwordHash, err := services.NewPasswordHash(newPassword)
	if err != nil {
		restError(w, r, err, http.StatusBadRequest)
		return
	}
	var account *services.Account
	if account, err = passwordHash.SetOnAccountByID(id); err != nil {
		restError(w, r, err, http.StatusBadRequest)
		return
	}
	outAccount := dtoFromAccount(account)
	if err = json.NewEncoder(w).Encode(outAccount); err != nil {
		restError(w, r, err, http.StatusInternalServerError)
		return
	}
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
	// TODO(ericm): Add email verification.
	serviceAccount := &services.Account{
		Email:        account.Email,
		Type:         accountType,
		PasswordHash: *hash,
	}
	if err = services.CreateAccount(serviceAccount); err != nil {
		restError(w, r, err, http.StatusInternalServerError)
		return
	}
	outAccount := dtoFromAccount(serviceAccount)
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
