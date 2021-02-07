package routes

import (
	"errors"
	"net/http"

	"github.com/cs3305-team-4/api/pkg/services"
	"github.com/gorilla/mux"
)

func InjectTutorsRoutes(subrouter *mux.Router) {
	// Profile routes
	subrouter.HandleFunc("/{uuid}/profile", handleProfileGet).Methods("GET")

	accountResource := subrouter.PathPrefix("/{uuid}").Subrouter()
	accountResource.Use(authAccount())

	accountResource.HandleFunc("/profile", handleProfilePost).Methods("POST")

	// Profile update routes
	accountResource.HandleFunc("/profile/avatar", handleProfileUpdateAvatar).Methods("POST")
	accountResource.HandleFunc("/profile/first-name", handleProfileUpdateFirstName).Methods("POST")
	accountResource.HandleFunc("/profile/last-name", handleProfileUpdateLastName).Methods("POST")
	accountResource.HandleFunc("/profile/city", handleProfileUpdateCity).Methods("POST")
	accountResource.HandleFunc("/profile/country", handleProfileUpdateCountry).Methods("POST")
	accountResource.HandleFunc("/profile/description", handleProfileUpdateDescription).Methods("POST")
	accountResource.HandleFunc("/profile/availability", handleTutorProfileAvailabilityPost).Methods("POST")

	accountResource.HandleFunc("/profile/qualifications", handleTutorProfileQualificationsPost).Methods("POST")
	accountResource.HandleFunc("/profile/qualifications/{qid}", handleTutorProfileQualificationsDelete).Methods("DELETE")
	accountResource.HandleFunc("/profile/work-experience", handleTutorProfileWorkExperiencePost).Methods("POST")
	accountResource.HandleFunc("/profile/work-experience/{wid}", handleTutorProfileWorkExperienceDelete).Methods("DELETE")

}

func handleTutorProfileQualificationsPost(w http.ResponseWriter, r *http.Request) {
	userID, err := getUUID(r, "uuid")
	if err != nil {
		restError(w, r, err, http.StatusBadRequest)
		return
	}
	dto := &QualificationsRequestDTO{}
	if !ParseBody(w, r, dto) {
		return
	}
	qualifications := &services.Qualification{
		Field:    dto.Field,
		Degree:   dto.Degree,
		School:   dto.School,
		Verified: false,
	}
	profile, err := qualifications.SetOnProfileByAccountID(userID)
	if err != nil {
		restError(w, r, err, http.StatusInternalServerError)
		return
	}
	if ok, err := profile.IsAccountType(services.Tutor); err != nil {
		restError(w, r, err, http.StatusInternalServerError)
		return
	} else if !ok {
		restError(w, r, errors.New("Account type does not match endpoint."), http.StatusBadRequest)
		return
	}
	profileDto := dtoFromProfile(profile, services.Tutor)
	WriteBody(w, r, profileDto)
}

func handleTutorProfileQualificationsDelete(w http.ResponseWriter, r *http.Request) {
	userID, err := getUUID(r, "uuid")
	if err != nil {
		restError(w, r, err, http.StatusBadRequest)
		return
	}
	qualificationID, err := getUUID(r, "qid")
	if err != nil {
		restError(w, r, err, http.StatusBadRequest)
		return
	}
	profile, err := services.ReadProfileByAccountID(userID, nil)
	if err != nil {
		restError(w, r, err, http.StatusBadRequest)
		return
	}
	if err = profile.RemoveQualificationByID(qualificationID); err != nil {
		restError(w, r, err, http.StatusBadRequest)
		return
	}
	if ok, err := profile.IsAccountType(services.Tutor); err != nil {
		restError(w, r, err, http.StatusInternalServerError)
		return
	} else if !ok {
		restError(w, r, errors.New("Account type does not match endpoint."), http.StatusBadRequest)
		return
	}
	profileDto := dtoFromProfile(profile, services.Tutor)
	WriteBody(w, r, profileDto)
}

func handleTutorProfileWorkExperiencePost(w http.ResponseWriter, r *http.Request) {
	userID, err := getUUID(r, "uuid")
	if err != nil {
		restError(w, r, err, http.StatusBadRequest)
		return
	}
	dto := &WorkExperienceRequestDTO{}
	if !ParseBody(w, r, dto) {
		return
	}
	exp := &services.WorkExperience{
		Role:        dto.Role,
		YearsExp:    dto.YearsExp,
		Description: dto.Description,
	}
	profile, err := exp.SetOnProfileByAccountID(userID)
	if err != nil {
		restError(w, r, err, http.StatusBadRequest)
		return
	}
	if ok, err := profile.IsAccountType(services.Tutor); err != nil {
		restError(w, r, err, http.StatusInternalServerError)
		return
	} else if !ok {
		restError(w, r, errors.New("Account type does not match endpoint."), http.StatusBadRequest)
		return
	}
	profileDto := dtoFromProfile(profile, services.Tutor)
	WriteBody(w, r, profileDto)
}

func handleTutorProfileWorkExperienceDelete(w http.ResponseWriter, r *http.Request) {
	userID, err := getUUID(r, "uuid")
	if err != nil {
		restError(w, r, err, http.StatusBadRequest)
		return
	}
	expID, err := getUUID(r, "wid")
	if err != nil {
		restError(w, r, err, http.StatusBadRequest)
		return
	}
	profile, err := services.ReadProfileByAccountID(userID, nil)
	if err != nil {
		restError(w, r, err, http.StatusBadRequest)
		return
	}
	if ok, err := profile.IsAccountType(services.Tutor); err != nil {
		restError(w, r, err, http.StatusInternalServerError)
		return
	} else if !ok {
		restError(w, r, errors.New("Account type does not match endpoint."), http.StatusBadRequest)
		return
	}
	if err = profile.RemoveWorkExperienceByID(expID); err != nil {
		restError(w, r, err, http.StatusBadRequest)
		return
	}
	profileDto := dtoFromProfile(profile, services.Tutor)
	WriteBody(w, r, profileDto)
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
	if err := validateUpdate("Availability", value, &TutorResponseDTO{}); err != nil {
		restError(w, r, err, http.StatusBadRequest)
		return
	}
	var profile *services.Profile
	if profile, err = services.UpdateProfileField(id, "availability", &value); err != nil {
		restError(w, r, err, http.StatusBadRequest)
		return
	}
	if ok, err := profile.IsAccountType(services.Tutor); err != nil {
		restError(w, r, err, http.StatusInternalServerError)
		return
	} else if !ok {
		restError(w, r, errors.New("Account type does not match endpoint."), http.StatusBadRequest)
		return
	}
	profileDto := dtoFromProfile(profile, services.Tutor)
	WriteBody(w, r, profileDto)
}
