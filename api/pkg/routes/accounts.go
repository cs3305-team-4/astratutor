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
type Account struct {
	ID           string `json:"id" validate:"len=0"`
	Email        string `json:"email" validate:"nonzero"`
	Type         string `json:"type"`
	Password     string `json:"password,omitempty" validate:"passwd"`
	ParentsEmail string `json:"parents_email,omitempty"`
}

func handleAccounts(w http.ResponseWriter, r *http.Request) {
	account := &Account{}
	if !ParseBody(w, r, account) {
		return
	}
	accountType, err := services.ToAccountType(account.Type)
	if err != nil {
		httpError(w, r, err, http.StatusBadRequest)
		return
	}
	hash, err := services.NewPasswordHash(account.Password)
	if err != nil {
		httpError(w, r, err, http.StatusInternalServerError)
		return
	}
	serviceAccount := &services.Account{
		Email:        account.Email,
		Type:         accountType,
		PasswordHash: *hash,
	}
	if err = services.CreateAccount(serviceAccount); err != nil {
		httpError(w, r, err, http.StatusInternalServerError)
		return
	}
	outAccount := &Account{
		ID:    serviceAccount.ID.String(),
		Email: serviceAccount.Email,
		Type:  string(serviceAccount.Type),
	}
	if err = json.NewEncoder(w).Encode(outAccount); err != nil {
		httpError(w, r, err, http.StatusInternalServerError)
		return
	}
}

func handleAccountsVerify(w http.ResponseWriter, r *http.Request) {
	id, err := getUUID(r)
	if err != nil {
		httpError(w, r, err, http.StatusBadRequest)
	}
	fmt.Println(id)
}
