package routes

import (
	"encoding/json"
	"net/http"
	"strings"

	"github.com/cs3305-team-4/api/pkg/services"
	"github.com/google/uuid"
	"github.com/gorilla/mux"
)

func InjectSubjectsRoutes(subrouter *mux.Router) {
	subrouter.HandleFunc("", handleSubjectsGet).Methods("GET")
	subrouter.HandleFunc("/tutors", handleSubjectTutorsGet).Methods("GET")
}

//Subject DTO represents an existing subject
type SubjectResponseDTO struct {
	Name string    `json:"name" validate:"required"`
	Slug string    `json:"slug" validate:"len=0"`
	ID   uuid.UUID `json:"subject_id" validate:"len=0"`
}

// Represents a tutors subject
type SubjectTaughtDTO struct {
	ID          uuid.UUID `json:"id" validate:"len=0"`
	Name        string    `json:"name" validate:"required"`
	Slug        string    `json:"slug" validate:"required"`
	Description string    `json:"description"`
	Price       float32   `json:"price" validate:"required"`
}

// Represents a Tutor and their subjects
type TutorSubjectsResponseDTO struct {
	ID          uuid.UUID          `json:"id" validate:"len=0"`
	FirstName   string             `json:"first_name" validate:"required"`
	LastName    string             `json:"last_name" validate:"required"`
	Avatar      string             `json:"avatar" validate:"omitempty,base64"`
	Slug        string             `json:"slug" validate:"len=0"`
	Description string             `json:"description"`
	Subjects    []SubjectTaughtDTO `json:"subjects"`
}

func ProfileToTutorSubjectsResponseDTO(profiles *[]services.Profile) *[]TutorSubjectsResponseDTO {
	tutorSubjectsResponse := []TutorSubjectsResponseDTO{}
	for _, profile := range *profiles {
		subjects := []SubjectTaughtDTO{}
		for _, subjectTaught := range profile.Subjects {
			subjects = append(subjects, SubjectTaughtDTO{
				ID:          subjectTaught.Subject.ID,
				Name:        subjectTaught.Subject.Name,
				Slug:        subjectTaught.Subject.Slug,
				Description: subjectTaught.Description,
				Price:       subjectTaught.Price,
			})
		}

		tutorSubjectsResponse = append(tutorSubjectsResponse, TutorSubjectsResponseDTO{
			ID:          profile.AccountID,
			FirstName:   profile.FirstName,
			LastName:    profile.LastName,
			Avatar:      profile.Avatar,
			Slug:        profile.Slug,
			Description: profile.Description,
			Subjects:    subjects,
		})
	}

	return &tutorSubjectsResponse
}

func SingleSubjectToDTO(subject *services.Subject) *SubjectResponseDTO {
	return &SubjectResponseDTO{
		Name: subject.Name,
		Slug: subject.Slug,
		ID:   subject.ID,
	}

}

func SubjectsToDTO(subjects []services.Subject) []SubjectResponseDTO {
	SubjectsDTO := []SubjectResponseDTO{}
	for _, subject := range subjects {
		SubjectsDTO = append(SubjectsDTO, *SingleSubjectToDTO(&subject))
	}
	return SubjectsDTO
}

//Handler functions:

//returns all subjects
func handleSubjectsGet(w http.ResponseWriter, r *http.Request) {
	serviceSubjects, err := services.GetSubjects(nil)
	if err != nil {
		restError(w, r, err, http.StatusBadRequest)
		return
	}
	outSubjects := SubjectsToDTO(serviceSubjects)
	if err = json.NewEncoder(w).Encode(outSubjects); err != nil {
		restError(w, r, err, http.StatusInternalServerError)
		return
	}

}

//checks for a filter and if none returns all tutors and if a filter returns tutors for that subject
func handleSubjectTutorsGet(w http.ResponseWriter, r *http.Request) {
	q := r.URL.Query()
	filter := q.Get("filter")
	if filter != "" {
		filtered, err := services.GetSubjectsBySlugs(strings.Split(filter, ","), nil)
		if err != nil {
			restError(w, r, err, http.StatusBadRequest)
			return
		}

		tutors, err := services.GetTutorsBySubjects(filtered, nil)
		if err != nil {
			restError(w, r, err, http.StatusBadRequest)
			return
		}
		outTutors := ProfileToTutorSubjectsResponseDTO(&tutors)

		if err = json.NewEncoder(w).Encode(outTutors); err != nil {
			restError(w, r, err, http.StatusInternalServerError)
			return
		}

	} else {
		tutors, err := services.GetAllTutors(nil)
		if err != nil {
			restError(w, r, err, http.StatusBadRequest)
			return
		}
		outTutors := ProfileToTutorSubjectsResponseDTO(&tutors)
		if err = json.NewEncoder(w).Encode(outTutors); err != nil {
			restError(w, r, err, http.StatusInternalServerError)
			return
		}
	}

}
