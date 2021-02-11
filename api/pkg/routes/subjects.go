package routes

import (
	"encoding/json"
	"net/http"

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
	Name  string    `json:"name" validate:"required"`
	Slug  string    `json:"slug" validate:"len=0"`
	ID    uuid.UUID `json:"subject_id" validate:"len=0"`
	Image string    `json:"image" validate:"omitempty,base64"`
}

//TutorSubjectResponseDTO
type TutorSubjectResponseDTO struct {
	SubjectTaughtID uuid.UUID `json:"subject_taught_id" validate:"len=0"`
	SubjectName     string    `json:"subject_name" validate:"required"`
	SubjectID       uuid.UUID `json:"subject_id" validate:"len=0"`
	TutorFirstName  string    `json:"tutor_first_name" validate:"required"`
	TutorLastName   string    `json:"tutor_last_name" validate:"required"`
	TutorAvatar     string    `json:"tutor_avatar" validate:"omitempty,base64"`
	TutorAccountID  uuid.UUID `json:"tutor_id" validate:"required"`
	TutorSlug       string    `json:"tutor_slug" validate:"len=0"`
	Price           uint      `json:"price" validate:"required"`
	Description     string    `json:"description"`
}

func SingleTutorSubjectToDTO(subjectTaught *services.SubjectTaught) *TutorSubjectResponseDTO {

	tutor, err := services.ReadProfileByAccountID(subjectTaught.TutorID, nil)
	if err != nil {
		return nil
	}
	subject, err := services.GetSubjectByID(subjectTaught.SubjectID, nil)
	if err != nil {
		return nil
	}
	return &TutorSubjectResponseDTO{
		SubjectName:     subject.Name,
		SubjectTaughtID: subjectTaught.ID,
		SubjectID:       subjectTaught.SubjectID,
		Price:           subjectTaught.Price,
		Description:     subjectTaught.Description,
		TutorFirstName:  tutor.FirstName,
		TutorLastName:   tutor.LastName,
		TutorAvatar:     tutor.Avatar,
		TutorSlug:       tutor.Slug,
		TutorAccountID:  subjectTaught.TutorID,
	}

}

func TutorSubjectsToDTO(tutorSubjects []services.SubjectTaught) []TutorSubjectResponseDTO {
	tutorSubjectsDTO := []TutorSubjectResponseDTO{}
	for _, item := range tutorSubjects {
		tutorSubjectsDTO = append(tutorSubjectsDTO, *SingleTutorSubjectToDTO(&item))
	}
	return tutorSubjectsDTO
}

func SingleSubjectToDTO(subject *services.Subject) *SubjectResponseDTO {
	return &SubjectResponseDTO{
		Name:  subject.Name,
		Slug:  subject.Slug,
		ID:    subject.ID,
		Image: subject.Image,
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
		filtered, err := services.GetSubjectBySlug(filter, nil)
		if err != nil {
			restError(w, r, err, http.StatusBadRequest)
			return
		}

		tutors, err := services.GetTutorsBySubjectID(filtered.ID, nil)
		if err != nil {
			restError(w, r, err, http.StatusBadRequest)
			return
		}

		outTutors := TutorSubjectsToDTO(tutors)

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
		outTutors := TutorSubjectsToDTO(tutors)
		if err = json.NewEncoder(w).Encode(outTutors); err != nil {
			restError(w, r, err, http.StatusInternalServerError)
			return
		}
	}

}
