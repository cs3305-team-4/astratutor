package routes

import (
	"encoding/json"
	"net/http"

	"github.com/cs3305-team-4/api/pkg/services"
	"github.com/gorilla/mux"
)

func InjectTutorsRoutes(subrouter *mux.Router) {
	// Profile routes
	subrouter.HandleFunc("/{uuid}/profile", handleProfileGet).Methods("GET")
	subrouter.HandleFunc("/{uuid}/profile", handleProfilePost).Methods("POST")

	// Profile update routes
	subrouter.HandleFunc("/{uuid}/profile/avatar", handleProfileUpdateAvatar).Methods("POST")
	subrouter.HandleFunc("/{uuid}/profile/first_name", handleProfileUpdateFirstName).Methods("POST")
	subrouter.HandleFunc("/{uuid}/profile/last_name", handleProfileUpdateLastName).Methods("POST")
	subrouter.HandleFunc("/{uuid}/profile/city", handleProfileUpdateCity).Methods("POST")
	subrouter.HandleFunc("/{uuid}/profile/country", handleProfileUpdateCountry).Methods("POST")
	subrouter.HandleFunc("/{uuid}/profile/description", handleProfileUpdateDescription).Methods("POST")
	subrouter.HandleFunc("/{uuid}/profile/availability", handleTutorProfileAvailabilityPost).Methods("POST")

	// Lesson routes.
	subrouter.HandleFunc("/{uuid}/lessons", handleTutorsLessonsGet).Methods("GET")
	subrouter.HandleFunc("/{uuid}/lessons", handleTutorsLessonsPost).Methods("POST")
}

// UpdateAvailabilityDTO DTO.
type UpdateAvailabilityDTO struct {
	Value []bool `json:"value"`
}

func handleTutorProfileAvailabilityPost(w http.ResponseWriter, r *http.Request) {
	id, err := getUUID(r, "uuid")
	if err != nil {
		restError(w, r, err, http.StatusBadRequest)
		return
	}
	update := &UpdateAvailabilityDTO{}
	if !ParseBody(w, r, update) {
		return
	}
	value := services.Availability(update.Value)
	if err := validateUpdate("Availability", value, &ProfileDTO{}); err != nil {
		restError(w, r, err, http.StatusBadRequest)
		return
	}
	var profile *services.Profile
	if profile, err = services.UpdateProfileField(id, "availability", &value); err != nil {
		restError(w, r, err, http.StatusBadRequest)
		return
	}
	dto := dtoFromProfile(profile, services.Tutor)
	if err = json.NewEncoder(w).Encode(dto); err != nil {
		restError(w, r, err, http.StatusInternalServerError)
		return
	}
}

func handleTutorsLessonsGet(w http.ResponseWriter, r *http.Request) {
	id, err := getUUID(r, "uuid")
	if err != nil {
		restError(w, r, err, http.StatusBadRequest)
		return
	}

	lessons, err := services.ReadLessonsByTutorID(id)
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

func handleTutorsLessonsPost(w http.ResponseWriter, r *http.Request) {
	id, err := getUUID(r, "uuid")
	if err != nil {
		restError(w, r, err, http.StatusBadRequest)
		return
	}

	lessonRequest := &LessonRequestDTO{}
	if !ParseBody(w, r, lessonRequest) {
		return
	}

	tutor, err := services.ReadTutorByID(id, nil)
	if err != nil {
		restError(w, r, err, http.StatusBadRequest)
		return
	}

	student, err := services.ReadStudentByID(lessonRequest.RequesterID, nil)
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
