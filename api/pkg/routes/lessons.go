package routes

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"mime"
	"net/http"
	"strings"
	"time"

	"github.com/cs3305-team-4/api/pkg/services"
	"github.com/google/uuid"
	"github.com/gorilla/mux"
	log "github.com/sirupsen/logrus"
)

// LessonResponseDTO represents an existing lesson
type LessonResponseDTO struct {
	ID uuid.UUID `json:"id"`

	// Time of the lesson
	StartTime time.Time `json:"start_time"`

	// Requester of the lesson
	RequesterID uuid.UUID `json:"requester_id"`

	// StudentID
	StudentID uuid.UUID `json:"student_id"`

	// TutorID represents the id of the tutor teaching the lesson
	TutorID uuid.UUID `json:"tutor_id"`

	// SubjectTaughtID
	SubjectID uuid.UUID `json:"subject_id"`

	// LessonDetail contains notes about what the student needs out of this lesson
	LessonDetail string `json:"lesson_detail"`

	// RequestStagedetermines what state of request the lesson is in
	RequestStage services.LessonRequestStage `json:"request_stage"`

	// RequestStageDetail contains a string related to the current request state
	RequestStageDetail string `json:"request_stage_detail"`

	// RequestStageChanger contains a reference to the account of the person who last changed the state of the lesson
	RequestStageChangerID uuid.UUID `json:"request_stage_changer_id"`

	Resources []ResourceMetadataDTO `json:"resources"`
}

// ResourceMetadataDTO represents a data transfer object for a resources metadata
type ResourceMetadataDTO struct {
	Name string `json:"name"`
	MIME string `json:"mime"`
}

// LessonRequestDTO represents a lesson that was first requested by an account
type LessonRequestDTO struct {
	// Time of the lesson
	StartTime time.Time `json:"start_time"`

	// The ID of who needs the lesson, can be a teacher or student
	// RequesterID must be a student if RequesteeID is a teacher and vice-versa
	StudentID uuid.UUID `json:"student_id"`

	// The ID of the person requesting the lesson, can be a teacher or a student
	// RequesterID must be a student if RequesteeID is a teacher and vice-versa
	SubjectTaughtID uuid.UUID `json:"subject_taught_id"`

	// LessonDetail contains info about what the lesson should be about
	LessonDetail string `json:"lesson_detail"`
}

// Represents a request to deny a lesson
type LessonDenyRequestDTO struct {
	Reason string `json:"reason"`
}

// Represents a request to cancel a lesson
type LessonCancelRequestDTO struct {
	Reason string `json:"reason"`
}

type LessonRescheduleRequestDTO struct {
	NewTime time.Time `json:"new_time"`
	Reason  string    `json:"reason"`
}

func dtoFromResourceMetadata(m *services.ResourceMetadata) *ResourceMetadataDTO {
	return &ResourceMetadataDTO{
		Name: m.Name,
		MIME: m.MIME,
	}
}

func dtoFromLesson(l *services.Lesson) *LessonResponseDTO {
	mds := []ResourceMetadataDTO{}

	for _, md := range l.Resources {
		mds = append(mds, *dtoFromResourceMetadata(&md))
	}

	return &LessonResponseDTO{
		ID:                    l.ID,
		StartTime:             l.StartTime,
		TutorID:               l.TutorID,
		StudentID:             l.StudentID,
		RequesterID:           l.RequesterID,
		SubjectID:             l.SubjectTaught.SubjectID,
		LessonDetail:          l.LessonDetail,
		RequestStage:          l.RequestStage,
		RequestStageDetail:    l.RequestStageDetail,
		RequestStageChangerID: l.RequestStageChangerID,
		Resources:             mds,
	}
}

func dtoFromLessons(lessons []services.Lesson) []LessonResponseDTO {
	dtoLessons := []LessonResponseDTO{}

	for _, l := range lessons {
		dtoLessons = append(dtoLessons, *dtoFromLesson(&l))
	}

	return dtoLessons
}

func InjectLessonsRoutes(subrouter *mux.Router) {
	// User needs an account to do anything with lessons
	subrouter.Use(authRequired)

	// POST /
	subrouter.HandleFunc("", handleLessonsPost).Methods("POST")

	lessonResource := subrouter.PathPrefix("/{uuid}").Subrouter()

	// Only allow users access to lessons that concern them
	lessonResource.Use(authMiddleware(
		func(w http.ResponseWriter, r *http.Request, ac *AuthContext) error {
			id, err := getUUID(r, "uuid")
			if err != nil {
				return err
			}

			lesson, err := services.ReadLessonByID(id)
			if err != nil {
				return err
			}

			if ac.Account.ID == lesson.StudentID || ac.Account.ID == lesson.TutorID {
				return nil
			}

			return errors.New("can only operate on a lesson that you are a participant in")
		}, true,
	))

	// GET /{uuid}
	lessonResource.Path("").HandlerFunc(handleLessonsGet).Methods("GET")

	// POST /{uuid}/accept
	lessonResource.HandleFunc("/payment-required",
		handleLessonsMarkPaymentRequired,
	).Methods("POST")

	// POST /{uuid}/deny
	lessonResource.HandleFunc("/deny",
		handleLessonsDenyRequest,
	).Methods("POST")

	// POST /{uuid}/cancel
	lessonResource.HandleFunc("/cancel",
		handleLessonsCancelRequest,
	).Methods("POST")

	lessonResource.HandleFunc("/payment-intent-secret",
		handleLessonsBillingGetPaymentIntentSecret,
	).Methods("GET")

	lessonResource.HandleFunc("/schedule",
		handleLessonsMarkScheduled,
	).Methods("POST")

	// POST /{uuid}/cancel
	lessonResource.HandleFunc("/reschedule",
		handleLessonsRescheduleRequest,
	).Methods("POST")

	subrouter.HandleFunc("/resources",
		handleLessonsResourcesPost).Methods("POST")

	subrouter.HandleFunc("/resources/{rid}",
		handleLessonsResourceGet).Methods("GET")
}

// type LessonCheckoutSessionResponseDTO struct {
// 	ID string `json:"id"`
// }

// func handleLessonsCreateCheckoutSession(w http.ResponseWriter, r *http.Request) {
// 	id, err := getUUID(r, "uuid")
// 	if err != nil {
// 		restError(w, r, err, http.StatusBadRequest)
// 		return
// 	}

// 	lesson, err := services.ReadLessonByID(id)
// 	if err != nil {
// 		restError(w, r, err, http.StatusBadRequest)
// 		return
// 	}

// 	sid, err := lesson.CreateCheckoutSession()
// 	if err != nil {
// 		restError(w, r, err, http.StatusInternalServerError)
// 		return
// 	}

// 	WriteBody(w, r, &LessonCheckoutSessionResponseDTO{
// 		ID: sid,
// 	})

// 	return
// }

type LessonBillingPaymentIntentSecretDTO struct {
	ID string `json:"id"`
}

func handleLessonsBillingGetPaymentIntentSecret(w http.ResponseWriter, r *http.Request) {
	id, err := getUUID(r, "uuid")
	if err != nil {
		restError(w, r, err, http.StatusBadRequest)
		return
	}

	lesson, err := services.ReadLessonByID(id)
	if err != nil {
		restError(w, r, err, http.StatusBadRequest)
		return
	}

	pid, err := lesson.GetPaymentIntentClientSecret()
	if err != nil {
		restError(w, r, err, http.StatusInternalServerError)
		return
	}

	WriteBody(w, r, &LessonBillingPaymentIntentSecretDTO{
		ID: pid,
	})

	return
}

func handleLessonsMarkScheduled(w http.ResponseWriter, r *http.Request) {
	authContext, err := ReadRequestAuthContext(r)
	if err != nil {
		restError(w, r, err, http.StatusBadRequest)
		return
	}

	id, err := getUUID(r, "uuid")
	if err != nil {
		restError(w, r, err, http.StatusBadRequest)
		return
	}

	lesson, err := services.ReadLessonByID(id)
	if err != nil {
		restError(w, r, err, http.StatusBadRequest)
		return
	}

	err = lesson.MarkScheduled(authContext.Account)
	if err != nil {
		restError(w, r, err, http.StatusBadRequest)
		return
	}
}

func handleLessonsGet(w http.ResponseWriter, r *http.Request) {
	id, err := getUUID(r, "uuid")
	if err != nil {
		restError(w, r, err, http.StatusBadRequest)
		return
	}

	lesson, err := services.ReadLessonByID(id)
	if err != nil {
		restError(w, r, err, http.StatusBadRequest)
		return
	}

	dtoLesson := dtoFromLesson(lesson)
	if err = json.NewEncoder(w).Encode(dtoLesson); err != nil {
		restError(w, r, err, http.StatusBadRequest)
		return
	}
}

func handleLessonsPost(w http.ResponseWriter, r *http.Request) {
	lessonRequest := &LessonRequestDTO{}
	if !ParseBody(w, r, lessonRequest) {
		return
	}

	authContext, err := ReadRequestAuthContext(r)
	if err != nil {
		restError(w, r, err, http.StatusBadRequest)
		return
	}

	subjectTaught, err := services.GetSubjectTaughtByID(lessonRequest.SubjectTaughtID, nil)
	if err != nil {
		restError(w, r, err, http.StatusBadRequest)
		return
	}

	if !(authContext.Account.ID == lessonRequest.StudentID || authContext.Account.ID == subjectTaught.TutorProfile.AccountID) {
		restError(w, r, errors.New("only allowed operate on a lesson that has your account as a participant"), http.StatusBadRequest)
		return
	}

	student, err := services.ReadAccountByID(lessonRequest.StudentID, nil)
	if err != nil {
		restError(w, r, err, http.StatusBadRequest)
		return
	}

	// tutor, err := services.ReadAccountByID(subjectTaught.TutorProfile.AccountID, nil)
	// if err != nil {
	// 	restError(w, r, err, http.StatusBadRequest)
	// 	return
	// }

	err = services.RequestLesson(authContext.Account, student, subjectTaught, lessonRequest.StartTime, lessonRequest.LessonDetail)
	if err != nil {
		restError(w, r, err, http.StatusBadRequest)
		return
	}
}

func handleLessonsMarkPaymentRequired(w http.ResponseWriter, r *http.Request) {
	authContext, err := ReadRequestAuthContext(r)
	if err != nil {
		restError(w, r, err, http.StatusBadRequest)
		return
	}

	id, err := getUUID(r, "uuid")
	if err != nil {
		restError(w, r, err, http.StatusBadRequest)
		return
	}

	lesson, err := services.ReadLessonByID(id)
	if err != nil {
		restError(w, r, err, http.StatusBadRequest)
		return
	}

	err = lesson.MarkPaymentRequired(authContext.Account)
	if err != nil {
		restError(w, r, err, http.StatusBadRequest)
		return
	}
}

func handleLessonsDenyRequest(w http.ResponseWriter, r *http.Request) {
	denyRequest := &LessonDenyRequestDTO{}
	if !ParseBody(w, r, denyRequest) {
		return
	}

	authContext, err := ReadRequestAuthContext(r)
	if err != nil {
		restError(w, r, err, http.StatusBadRequest)
		return
	}

	id, err := getUUID(r, "uuid")
	if err != nil {
		restError(w, r, err, http.StatusBadRequest)
		return
	}

	lesson, err := services.ReadLessonByID(id)
	if err != nil {
		restError(w, r, err, http.StatusBadRequest)
		return
	}

	err = lesson.MarkDenied(authContext.Account, denyRequest.Reason)
	if err != nil {
		restError(w, r, err, http.StatusBadRequest)
		return
	}
}

func handleLessonsCancelRequest(w http.ResponseWriter, r *http.Request) {
	cancelRequest := &LessonCancelRequestDTO{}
	if !ParseBody(w, r, cancelRequest) {
		return
	}

	authContext, err := ReadRequestAuthContext(r)
	if err != nil {
		restError(w, r, err, http.StatusBadRequest)
		return
	}

	id, err := getUUID(r, "uuid")
	if err != nil {
		restError(w, r, err, http.StatusBadRequest)
		return
	}

	lesson, err := services.ReadLessonByID(id)
	if err != nil {
		restError(w, r, err, http.StatusBadRequest)
		return
	}

	err = lesson.MarkCancelled(authContext.Account, cancelRequest.Reason)
	if err != nil {
		restError(w, r, err, http.StatusBadRequest)
		return
	}
}

func handleLessonsRescheduleRequest(w http.ResponseWriter, r *http.Request) {
	rescheduleRequest := &LessonRescheduleRequestDTO{}
	if !ParseBody(w, r, rescheduleRequest) {
		return
	}

	authContext, err := ReadRequestAuthContext(r)
	if err != nil {
		restError(w, r, err, http.StatusBadRequest)
		return
	}

	id, err := getUUID(r, "uuid")
	if err != nil {
		restError(w, r, err, http.StatusBadRequest)
		return
	}

	lesson, err := services.ReadLessonByID(id)
	if err != nil {
		restError(w, r, err, http.StatusBadRequest)
		return
	}

	err = lesson.MarkRescheduled(authContext.Account, rescheduleRequest.NewTime, rescheduleRequest.Reason)
	if err != nil {
		restError(w, r, err, http.StatusBadRequest)
		return
	}
}

func resourceAllowedMIME(mime string) bool {
	return (mime == "application/pdf")
}

func handleLessonsResourceGet(w http.ResponseWriter, r *http.Request) {
	lid, err := getUUID(r, "uuid")
	if err != nil {
		restError(w, r, err, http.StatusBadRequest)
		return
	}

	rid, err := getUUID(r, "rid")
	if err != nil {
		restError(w, r, err, http.StatusBadRequest)
		return
	}

	lesson, err := services.ReadLessonByID(lid)
	if err != nil {
		restError(w, r, err, http.StatusBadRequest)
		return
	}

	resource, err := lesson.ReadResourceByID(rid)
	if err != nil {
		restError(w, r, err, http.StatusBadRequest)
		return
	}

	data, err := resource.GetData()
	w.Header().Add("Content-Type", resource.MIME)
	w.Header().Add("Content-Length", fmt.Sprintf("%d", len(data)))
	_, err = w.Write(data)

	if err != nil {
		log.Error(fmt.Errorf("error writing data, %s", err))
	}

	return
}

func handleLessonsResourcesPost(w http.ResponseWriter, r *http.Request) {
	id, err := getUUID(r, "uuid")
	if err != nil {
		restError(w, r, err, http.StatusBadRequest)
		return
	}

	lesson, err := services.ReadLessonByID(id)
	if err != nil {
		restError(w, r, err, http.StatusBadRequest)
		return
	}

	r.ParseMultipartForm(16777216) // 16mb max

	file, header, err := r.FormFile("file")
	if err != nil {
		restError(w, r, err, http.StatusBadRequest)
	}
	defer file.Close()

	nameSplit := strings.Split(header.Filename, ".")
	filenameExclExt := nameSplit[0]
	fileExt := nameSplit[len(nameSplit)-1]

	mimeType := mime.TypeByExtension(fileExt)
	if !resourceAllowedMIME(mimeType) {
		restError(w, r, fmt.Errorf("%s is not an allowed resource type", mimeType), http.StatusBadRequest)
		return
	}

	buf := bytes.NewBuffer(nil)
	if _, err := io.Copy(buf, file); err != nil {
		restError(w, r, err, http.StatusBadRequest)
	}

	err = lesson.CreateResource(filenameExclExt, mimeType, buf.Bytes())
	if err != nil {
		restError(w, r, err, http.StatusBadRequest)
		return
	}
}
