package routes

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/cs3305-team-4/api/pkg/services"
	"github.com/gorilla/mux"
)

func InjectAccountsRoutes(subrouter *mux.Router) {
	subrouter.HandleFunc("", handleAccounts).Methods("POST")
	subrouter.HandleFunc("/{uuid}/verify", handleAccountsVerify).Methods("POST")
}

// Account DTO.
type AccountDTO struct {
	ID           string `json:"id" validate:"len=0"`
	Email        string `json:"email" validate:"required,email"`
	Type         string `json:"type" validate:"required"`
	Password     string `json:"password,omitempty" validate:"required,passwd"`
	ParentsEmail string `json:"parents_email,omitempty" validate:"omitempty,email"`
}

func handleAccounts(w http.ResponseWriter, r *http.Request) {
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
