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
	subrouter.HandleFunc("/{uuid}/lessons", handleStudentsLessonsGet).Methods("GET")
	subrouter.HandleFunc("/{uuid}/lessons", handleStudentsLessonsPost).Methods("POST")
}

// Profile DTO.
type ProfileDTO struct {
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

func dtoFromProfile(p *services.Profile) *ProfileDTO {
	return &ProfileDTO{
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
	id, err := getUUID(r, "uuid")
	if err != nil {
		restError(w, r, err, http.StatusBadRequest)
		return
	}
	serviceProfile, err := services.ReadProfileByAccountID(id, nil)
	if err != nil {
		restError(w, r, err, http.StatusBadRequest)
		return
	}
	outProfile := dtoFromProfile(serviceProfile)
	if err = json.NewEncoder(w).Encode(outProfile); err != nil {
		restError(w, r, err, http.StatusInternalServerError)
		return
	}
}

func handleStudentsProfilePost(w http.ResponseWriter, r *http.Request) {
	id, err := getUUID(r, "uuid")
	if err != nil {
		restError(w, r, err, http.StatusBadRequest)
		return
	}
	profile := &ProfileDTO{}
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
		restError(w, r, err, http.StatusInternalServerError)
		return
	}
	outProfile := dtoFromProfile(serviceProfile)
	if err = json.NewEncoder(w).Encode(outProfile); err != nil {
		restError(w, r, err, http.StatusInternalServerError)
		return
	}
}

func handleStudentsLessonsGet(w http.ResponseWriter, r *http.Request) {
	id, err := getUUID(r, "uuid")
	if err != nil {
		restError(w, r, err, http.StatusBadRequest)
		return
	}

	lessons, err := services.ReadLessonsByStudentID(id)
	if err != nil {
		restError(w, r, err, http.StatusBadRequest)
		return
	}

	dtoLessons := dtoFromLessons(lessons)
	if err = json.NewEncoder(w).Encode(dtoLessons); err != nil {
		restError(w, r, err, http.StatusBadRequest)
		return
	}
}

func handleStudentsLessonsPost(w http.ResponseWriter, r *http.Request) {
	id, err := getUUID(r, "uuid")
	if err != nil {
		restError(w, r, err, http.StatusBadRequest)
		return
	}

	lessonRequest := &LessonRequestDTO{}
	if !ParseBody(w, r, lessonRequest) {
		return
	}

	tutor, err := services.ReadTutorByID(lessonRequest.RequesterID, nil)
	if err != nil {
		restError(w, r, err, http.StatusBadRequest)
		return
	}

	student, err := services.ReadStudentByID(id, nil)
	if err != nil {
		restError(w, r, err, http.StatusBadRequest)
		return
	}

	err = services.CreateLesson(student, tutor, lessonRequest.TimeStarts, lessonRequest.LessonDetail)
	if err != nil {
		restError(w, r, err, http.StatusBadRequest)
		return
	}
}
