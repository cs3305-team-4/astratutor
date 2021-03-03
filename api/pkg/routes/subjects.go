package routes

import (
	"encoding/json"
	"net/http"
	"strconv"
	"strings"

	"github.com/cs3305-team-4/api/pkg/services"
	"github.com/google/uuid"
	"github.com/gorilla/mux"
)

func InjectSubjectsRoutes(subrouter *mux.Router) {
	subrouter.HandleFunc("", handleSubjectsGet).Methods("GET")
	subrouter.HandleFunc("", handleRequestSubject).Nethods("POST")
	subrouter.HandleFunc("/tutors", handleSubjectTutorsGet).Methods("GET")
	subrouter.HandleFunc("/tutors/{tid}", handleGetSubjectsForTutor).Methods("GET")
}

//Subject DTO represents an existing subject
type SubjectResponseDTO struct {
	Name string    `json:"name" validate:"required"`
	Slug string    `json:"slug" validate:"len=0"`
	ID   uuid.UUID `json:"id" validate:"len=0"`
}

type SubjectRequestDTO struct {
	Name string `json:"name" validate:"required"`
}

// Represents a tutors subject
type SubjectTaughtDTO struct {
	ID          uuid.UUID `json:"id" validate:"len=0"`
	SubjectID   uuid.UUID `json:"subject_id" validate:"len=0"`
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

// SubjectTaughtRequestDTO represents a subject a Tutor wishes to teach
type SubjectTaughtRequestDTO struct {
	Description string  `json:"description"`
	Price       float32 `json:"price"`
}

// SubjectTaughtDescriptionUpdateRequestDTO represents a subject a Tutor wishes to update the description for
type SubjectTaughtDescriptionUpdateRequestDTO struct {
	Description string `json:"description"`
}

// SubjectTaughtPriceUpdateRequestDTO represents a subject a Tutor wishes to update the Price for
type SubjectTaughtPriceUpdateRequestDTO struct {
	Price float32 `json:"price"`
}

func ProfileToTutorSubjectsResponseDTO(profiles *[]services.Profile) *[]TutorSubjectsResponseDTO {
	tutorSubjectsResponse := []TutorSubjectsResponseDTO{}
	for _, profile := range *profiles {
		tutorSubjectsResponse = append(tutorSubjectsResponse, TutorSubjectsResponseDTO{
			ID:          profile.AccountID,
			FirstName:   profile.FirstName,
			LastName:    profile.LastName,
			Avatar:      profile.Avatar,
			Slug:        profile.Slug,
			Description: profile.Description,
			Subjects:    SubjectsTuaghtToDTO(&profile.Subjects),
		})
	}

	return &tutorSubjectsResponse
}

func SubjectTaughtToDTO(subjectTaught *services.SubjectTaught) *SubjectTaughtDTO {
	return &SubjectTaughtDTO{
		ID:          subjectTaught.ID,
		SubjectID:   subjectTaught.Subject.ID,
		Name:        subjectTaught.Subject.Name,
		Slug:        subjectTaught.Subject.Slug,
		Description: subjectTaught.Description,
		Price:       subjectTaught.Price,
	}
}

func SubjectsTuaghtToDTO(subjectsTaught *[]services.SubjectTaught) []SubjectTaughtDTO {
	subjectsTaughtDto := []SubjectTaughtDTO{}
	for _, subjectTaught := range *subjectsTaught {
		subjectsTaughtDto = append(subjectsTaughtDto, *SubjectTaughtToDTO(&subjectTaught))
	}
	return subjectsTaughtDto
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

	pageSize, err := strconv.Atoi(q.Get("page_size"))
	if err != nil || pageSize <= 0 {
		pageSize = 10
	}
	page, err := strconv.Atoi(q.Get("page"))
	if err != nil || page <= 0 {
		page = 1
	}

	if filter != "" {
		filtered, err := services.GetSubjectsBySlugs(strings.Split(filter, ","), nil)
		if err != nil {
			restError(w, r, err, http.StatusBadRequest)
			return
		}

		tutors, total_pages, err := services.GetTutorsBySubjectsPaginated(filtered, pageSize, page, nil)
		if err != nil {
			restError(w, r, err, http.StatusBadRequest)
			return
		}
		outTutors := ToPaginatedDTO(total_pages, ProfileToTutorSubjectsResponseDTO(&tutors))

		if err = json.NewEncoder(w).Encode(outTutors); err != nil {
			restError(w, r, err, http.StatusInternalServerError)
			return
		}

	} else {
		tutors, totalPages, err := services.GetAllTutorsPaginated(nil, pageSize, page)
		if err != nil {
			restError(w, r, err, http.StatusBadRequest)
			return
		}
		outTutors := ToPaginatedDTO(totalPages, ProfileToTutorSubjectsResponseDTO(&tutors))
		if err = json.NewEncoder(w).Encode(outTutors); err != nil {
			restError(w, r, err, http.StatusInternalServerError)
			return
		}
	}

}

func handleGetSubjectsForTutor(w http.ResponseWriter, r *http.Request) {
	tid, err := getUUID(r, "tid")
	if err != nil {
		restError(w, r, err, http.StatusInternalServerError)
		return
	}

	tutorProfile, err := services.ReadProfileByAccountID(tid, nil)
	if err != nil {
		restError(w, r, err, http.StatusInternalServerError)
		return
	}

	tutorSubjects, err := services.GetSubjectsTaughtByTutorID(tutorProfile.ID, nil, "Subject")
	if err != nil {
		restError(w, r, err, http.StatusInternalServerError)
		return
	}

	if err = json.NewEncoder(w).Encode(SubjectsTuaghtToDTO(&tutorSubjects)); err != nil {
		restError(w, r, err, http.StatusInternalServerError)
		return
	}
}

func handleRequestSubject(w http.ResponseWriter, r *http.Request) {
	subjectRequest := &SubjectRequestDTO{}
	if !ParseBody(w, r, subjectRequest) {
		return
	}

	authContext, err := ReadRequestAuthContext(r)
	if err != nil {
		restError(w, r, err, http.StatusBadRequest)
		return
	}

	if err = services.RequestSubject(authContext.Account, subjectRequest.Name); err != nil {
		restError(w, r, err, http.StatusBadRequest)
		return
	}
}
