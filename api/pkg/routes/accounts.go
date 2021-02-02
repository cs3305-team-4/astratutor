package routes

import (
	"encoding/json"
	"net/http"

	"github.com/gorilla/mux"
)

func InjectAccountsRoutes(subrouter *mux.Router) {
	subrouter.HandleFunc("/", handleAccounts).Methods("POST")
}

type Account struct {
	ID       string `json:"id"`
	Email    string `json:"email"`
	Type     string `json:"type"`
	Password string `json:"password"`
}

func handleAccounts(w http.ResponseWriter, r *http.Request) {
	account := &Account{}
	if err := json.NewDecoder(r.Body).Decode(account); err != nil {
		httpError(w, r, err, http.StatusBadRequest)
		return
	}
}
