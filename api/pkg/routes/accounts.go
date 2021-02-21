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
	subrouter.HandleFunc("", handleAccountsPost).Methods("POST")

	// Only allow users to access routes relevant to their own account
	accountResource := subrouter.PathPrefix("/{uuid}").Subrouter()
	accountResource.Use(authAccount())
	accountResource.HandleFunc("", handleAccountsGet).Methods("GET")
	accountResource.HandleFunc("/verify", handleAccountsVerify).Methods("POST")
	accountResource.HandleFunc("/email", handleAccountsUpdateEmail).Methods("POST")
	accountResource.HandleFunc("/password", handleAccountsUpdatePassword).Methods("POST")
	accountResource.HandleFunc("/lessons", handleAccountsLessonsGet).Methods("GET")

	accountResource.HandleFunc("/billing/tutor-onboard", handleTutorBillingGetOnboard).Methods("GET")
	accountResource.HandleFunc("/billing/tutor-onboard-url", handleTutorBillingGetOnboardURL).Methods("GET")
	accountResource.HandleFunc("/billing/tutor-requirements-met", handleTutorBillingGetRequirementsMet).Methods("GET")
	accountResource.HandleFunc("/billing/tutor-panel-url", handleTutorBillingGetPanelURL).Methods("GET")
	accountResource.HandleFunc("/billing/payouts", handleTutorBillingCreatePayout).Methods("POST")
}

// AccountDTO return DTO.
type AccountResponseDTO struct {
	ID           string `json:"id" validate:"required,uuid"`
	Email        string `json:"email" validate:"required,email"`
	Type         string `json:"type" validate:"required"`
	ParentsEmail string `json:"parents_email,omitempty" validate:"omitempty,email"`
}

// AccountRequestDTO request DTO.
type AccountRequestDTO struct {
	Email        string `json:"email" validate:"required,email"`
	Type         string `json:"type" validate:"required"`
	Password     string `json:"password" validate:"required,passwd"`
	ParentsEmail string `json:"parents_email,omitempty" validate:"omitempty,email"`
}

func dtoFromAccount(a *services.Account) *AccountResponseDTO {
	return &AccountResponseDTO{
		ID:    a.ID.String(),
		Email: a.Email,
		Type:  string(a.Type),
	}
}

func handleTutorBillingGetOnboard(w http.ResponseWriter, r *http.Request) {
	id, err := getUUID(r, "uuid")
	if err != nil {
		restError(w, r, err, http.StatusBadRequest)
		return
	}
	serviceAccount, err := services.ReadAccountByID(id, nil)
	if err != nil {
		restError(w, r, err, http.StatusNotFound)
		return
	}

	ready, err := serviceAccount.IsTutorBillingOnboarded()
	if err != nil {
		restError(w, r, err, http.StatusInternalServerError)
		return
	}

	if ready {
		w.WriteHeader(http.StatusOK)
		return
	}

	w.WriteHeader(http.StatusNotFound)
	return
}

func handleTutorBillingCreatePayout(w http.ResponseWriter, r *http.Request) {
	id, err := getUUID(r, "uuid")
	if err != nil {
		restError(w, r, err, http.StatusBadRequest)
		return
	}
	serviceAccount, err := services.ReadAccountByID(id, nil)
	if err != nil {
		restError(w, r, err, http.StatusNotFound)
		return
	}

	ready, err := serviceAccount.IsTutorBillingOnboarded()
	if err != nil {
		restError(w, r, err, http.StatusInternalServerError)
		return
	}

	if ready {
		w.WriteHeader(http.StatusOK)
		return
	}

	w.WriteHeader(http.StatusNotFound)
	return
}

func handleTutorBillingGetRequirementsMet(w http.ResponseWriter, r *http.Request) {
	id, err := getUUID(r, "uuid")
	if err != nil {
		restError(w, r, err, http.StatusBadRequest)
		return
	}
	serviceAccount, err := services.ReadAccountByID(id, nil)
	if err != nil {
		restError(w, r, err, http.StatusNotFound)
		return
	}

	ready, err := serviceAccount.IsTutorBillingRequirementsMet()
	if err != nil {
		restError(w, r, err, http.StatusInternalServerError)
		return
	}

	if ready {
		w.WriteHeader(http.StatusOK)
		return
	}

	w.WriteHeader(http.StatusNotFound)
	return
}

type BillingOnboardURLResponseDTO struct {
	URL string `json:"url" validate:"required"`
}

type BillingPanelURLResponseDTO struct {
	URL string `json:"url" validate:"required"`
}

func handleTutorBillingGetOnboardURL(w http.ResponseWriter, r *http.Request) {
	id, err := getUUID(r, "uuid")
	if err != nil {
		restError(w, r, err, http.StatusBadRequest)
		return
	}

	serviceAccount, err := services.ReadAccountByID(id, nil)
	if err != nil {
		restError(w, r, err, http.StatusNotFound)
		return
	}

	url, err := serviceAccount.GetTutorBillingOnboardURL()
	if err != nil {
		restError(w, r, err, http.StatusInternalServerError)
		return
	}

	WriteBody(w, r, &BillingOnboardURLResponseDTO{
		URL: url,
	})
	return
}

func handleTutorBillingGetPanelURL(w http.ResponseWriter, r *http.Request) {
	id, err := getUUID(r, "uuid")
	if err != nil {
		restError(w, r, err, http.StatusBadRequest)
		return
	}

	serviceAccount, err := services.ReadAccountByID(id, nil)
	if err != nil {
		restError(w, r, err, http.StatusNotFound)
		return
	}

	url, err := serviceAccount.GetTutorBillingPanelURL()
	if err != nil {
		restError(w, r, err, http.StatusInternalServerError)
		return
	}

	WriteBody(w, r, &BillingOnboardURLResponseDTO{
		URL: url,
	})
	return
}

func handleAccountsGet(w http.ResponseWriter, r *http.Request) {
	id, err := getUUID(r, "uuid")
	if err != nil {
		restError(w, r, err, http.StatusBadRequest)
		return
	}
	serviceAccount, err := services.ReadAccountByID(id, nil)
	if err != nil {
		restError(w, r, err, http.StatusNotFound)
		return
	}
	outAccount := dtoFromAccount(serviceAccount)
	WriteBody(w, r, outAccount)
}

func handleAccountsUpdateEmail(w http.ResponseWriter, r *http.Request) {
	id, err := getUUID(r, "uuid")
	if err != nil {
		restError(w, r, err, http.StatusBadRequest)
		return
	}
	// TODO(ericm): Add email verification.
	field := ParseUpdateString(w, r)
	if err = validateUpdate("Email", field, &AccountRequestDTO{}); err != nil {
		restError(w, r, err, http.StatusBadRequest)
		return
	}
	var account *services.Account
	if account, err = services.UpdateAccountField(id, "email", field); err != nil {
		restError(w, r, err, http.StatusBadRequest)
		return
	}
	outAccount := dtoFromAccount(account)
	WriteBody(w, r, outAccount)
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
	oldPassword := update.Value.OldPassword
	newPassword := update.Value.NewPassword

	// Validate old password.
	oldPasswordHash, err := services.ReadPasswordHashByAccountID(id)
	if err != nil {
		restError(w, r, err, http.StatusBadRequest)
		return
	}
	if !oldPasswordHash.ValidMatch(oldPassword) {
		restError(w, r, errors.New("Old password provided is incorrect."), http.StatusForbidden)
		return
	}

	if err = validateUpdate("Password", newPassword, &AccountResponseDTO{}); err != nil {
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
	WriteBody(w, r, outAccount)
}

func handleAccountsPost(w http.ResponseWriter, r *http.Request) {
	account := &AccountRequestDTO{}
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
	WriteBody(w, r, outAccount)
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
