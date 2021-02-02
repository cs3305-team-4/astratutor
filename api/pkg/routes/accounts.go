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
	subrouter.HandleFunc("/test", func(w http.ResponseWriter, r *http.Request) {
		fmt.Println("test")
	})
}

type Account struct {
	ID           string `json:"id"`
	Email        string `json:"email"`
	Type         string `json:"type"`
	Password     string `json:"password,omitempty"`
	ParentsEmail string `json:"parents_email,omitempty"`
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
