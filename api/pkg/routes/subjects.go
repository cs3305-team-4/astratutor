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
	subrouter.HandleFunc("/tutors", handleSubjectTutorsGet).Methods("GET") //.Queries("filter", "{filter}")
}

//Subject DTO represents an existing subject
type Subject struct {
	Name  string    `json:"name"`
	Slug  string    `json:"slug"`
	ID    uuid.UUID `json:"subject_id"`
	Image string
}

//TutorSubjectDTO
type TutorSubject struct {
	SubjectName    string    `json:"subject_name"`
	SubjectID      uuid.UUID `json:"subject_id"`
	TutorFirstName string    `json:"tutor_first_name"`
	TutorLastName  string    `json:"tutor_last_name"`
	TutorAccountID uuid.UUID `json:"tutor_id"`
	Price          uint      `json:"price"`
	Description    string    `json:"description"`
	Image          string    `json:"image"`
}

func SingleTutorSubjectToDTO(subjectTaught *services.SubjectTaught) *TutorSubject {
	return &TutorSubject{
		SubjectName:    subjectTaught.Subject.Name,
		SubjectID:      subjectTaught.Subject.ID,
		Price:          subjectTaught.Price,
		Description:    subjectTaught.Description,
		TutorFirstName: subjectTaught.Tutor.Profile.FirstName,
		TutorLastName:  subjectTaught.Tutor.Profile.LastName,
		TutorAccountID: subjectTaught.Tutor.ID,
	}

}

func TutorSubjectsToDTO(tutorSubjects []services.SubjectTaught) []TutorSubject {
	tutorSubjectsDTO := []TutorSubject{}
	for _, item := range tutorSubjects {
		tutorSubjectsDTO = append(tutorSubjectsDTO, *SingleTutorSubjectToDTO(&item))
	}
	return tutorSubjectsDTO
}

func SingleSubjectToDTO(subject *services.Subject) *Subject {
	return &Subject{
		Name:  subject.Name,
		Slug:  subject.Slug,
		ID:    subject.ID,
		Image: subject.Image,
	}

}

func SubjectsToDTO(subjects []services.Subject) []Subject {
	SubjectsDTO := []Subject{}
	for _, subject := range subjects {
		SubjectsDTO = append(SubjectsDTO, *SingleSubjectToDTO(&subject))
	}
	return SubjectsDTO
}

//Handler functions:

//
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
		filtered, err := services.GetSubjectByName(filter, nil)
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
