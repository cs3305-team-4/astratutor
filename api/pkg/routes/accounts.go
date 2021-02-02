package routes

import (
	"encoding/json"
	"net/http"

	"github.com/cs3305-team-4/api/pkg/services"
	"github.com/gorilla/mux"
)

func InjectAccountsRoutes(subrouter *mux.Router) {
	subrouter.HandleFunc("/", handleAccounts).Methods("POST")
}

type Account struct {
	ID           string `json:"id"`
	Email        string `json:"email"`
	Type         string `json:"type"`
	Password     string `json:"password"`
	ParentsEmail string `json:"parents:email"`
}

func handleAccounts(w http.ResponseWriter, r *http.Request) {
	account := &Account{}
	if err := json.NewDecoder(r.Body).Decode(account); err != nil {
		httpError(w, r, err, http.StatusBadRequest)
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
	if err = services.CreateAccount(&services.Account{
		Email:        account.Email,
		Type:         accountType,
		PasswordHash: *hash,
	}); err != nil {
		httpError(w, r, err, http.StatusInternalServerError)
		return
	}
	if err = json.NewEncoder(w).Encode(account); err != nil {
		httpError(w, r, err, http.StatusInternalServerError)
		return
	}
}
