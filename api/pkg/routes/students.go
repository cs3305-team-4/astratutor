package routes

import (
	"encoding/json"
	"net/http"

	"github.com/cs3305-team-4/api/pkg/services"
	"github.com/gorilla/mux"
)

func InjectStudentsRoutes(subrouter *mux.Router) {
	subrouter.HandleFunc("/{uuid}/profile", handleStudentsProfileGet).Methods("GET")
	subrouter.HandleFunc("/{uuid}/profile", handleStudentsProfilePost).Methods("POST")
}

// Profile DTO.
type Profile struct {
	AccountID   string `json:"account_id" validate:"len=0"`
	ID          string `json:"id" validate:"len=0"`
	Avatar      string `json:"avatar"`
	Slug        string `json:"slug" validate:"len=0"`
	FirstName   string `json:"first_name" validate:"nonzero"`
	LastName    string `json:"last_name" validate:"nonzero"`
	City        string `json:"city" validate:"nonzero"`
	Country     string `json:"country" validate:"nonzero"`
	Description string `json:"description"`
}

func mapOutProfile(p *services.Profile) *Profile {
	return &Profile{
		AccountID:   p.AccountID.String(),
		ID:          p.ID.String(),
		Avatar:      p.Avatar,
		Slug:        p.Slug,
		FirstName:   p.FirstName,
		LastName:    p.LastName,
		City:        p.City,
		Country:     p.Country,
		Description: p.Description,
	}
}

func handleStudentsProfileGet(w http.ResponseWriter, r *http.Request) {
	id, err := getUUID(r)
	if err != nil {
		httpError(w, r, err, http.StatusBadRequest)
		return
	}
	serviceProfile, err := services.GetProfileByAccountID(id, nil)
	if err != nil {
		httpError(w, r, err, http.StatusBadRequest)
		return
	}
	outProfile := mapOutProfile(serviceProfile)
	if err = json.NewEncoder(w).Encode(outProfile); err != nil {
		httpError(w, r, err, http.StatusInternalServerError)
		return
	}
}

func handleStudentsProfilePost(w http.ResponseWriter, r *http.Request) {
	id, err := getUUID(r)
	if err != nil {
		httpError(w, r, err, http.StatusBadRequest)
		return
	}
	profile := &Profile{}
	if !ParseBody(w, r, profile) {
		return
	}
	serviceProfile := &services.Profile{
		AccountID:   id,
		Avatar:      profile.Avatar,
		Slug:        profile.Slug,
		FirstName:   profile.FirstName,
		LastName:    profile.LastName,
		City:        profile.City,
		Country:     profile.Country,
		Description: profile.Description,
	}
	if err := services.CreateProfile(serviceProfile); err != nil {
		httpError(w, r, err, http.StatusInternalServerError)
		return
	}
	outProfile := mapOutProfile(serviceProfile)
	if err = json.NewEncoder(w).Encode(outProfile); err != nil {
		httpError(w, r, err, http.StatusInternalServerError)
		return
	}
}
