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

// LessonDTO represents an existing lesson
type LessonDTO struct {
	ID uuid.UUID `json:"id"`

	// Time of the lesson
	StartTime time.Time `json:"start_time"`

	// Requester of the lesson
	RequesterID uuid.UUID `json:"requester_id"`

	// StudentID
	StudentID uuid.UUID `json:"student_id"`

	// TutorID
	TutorID uuid.UUID `json:"tutor_id"`

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
	TutorID uuid.UUID `json:"tutor_id"`

	// LessonDetail contains info about what the lesson should be about
	LessonDetail string `json:"lesson_detail"`
}

// LessonRequesStagetDTO represents a change in request stage for a lesson
// i.e confirmed, expired, etc
type LessonStageChangeDTO struct {
	// The ID of who wants the stage change, can be a teacher or student
	// i.e if they request the tutors lesson create endpoint, it'l expect RequesterID to be a student
	RequesterID uuid.UUID `json:"requester_id"`

	LessonDetail string `json:"lesson_detail"`
}

func dtoFromResourceMetadata(m *services.ResourceMetadata) *ResourceMetadataDTO {
	return &ResourceMetadataDTO{
		Name: m.Name,
		MIME: m.MIME,
	}
}

func dtoFromLesson(l *services.Lesson) *LessonDTO {
	mds := []ResourceMetadataDTO{}

	for _, md := range l.Resources {
		mds = append(mds, *dtoFromResourceMetadata(&md))
	}

	return &LessonDTO{
		ID:                    l.ID,
		StartTime:             l.StartTime,
		TutorID:               l.TutorID,
		StudentID:             l.StudentID,
		RequesterID:           l.RequesterID,
		LessonDetail:          l.LessonDetail,
		RequestStage:          l.RequestStage,
		RequestStageDetail:    l.RequestStageDetail,
		RequestStageChangerID: l.RequestStageChangerID,
		Resources:             mds,
	}
}

func dtoFromLessons(lessons []services.Lesson) []LessonDTO {
	dtoLessons := []LessonDTO{}

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
	lessonResource.PathPrefix("").HandlerFunc(handleLessonsGet).Methods("GET")

	// POST /{uuid}/accept
	lessonResource.HandleFunc("/accept",
		handleLessonsRequestStageChangeClosure(services.Accepted),
	).Methods("POST")

	// POST /{uuid}/deny
	lessonResource.HandleFunc("/deny",
		handleLessonsRequestStageChangeClosure(services.Denied),
	).Methods("POST")

	// POST /{uuid}/cancel
	lessonResource.HandleFunc("/cancel",
		handleLessonsRequestStageChangeClosure(services.Cancelled),
	).Methods("POST")

	subrouter.HandleFunc("/resources",
		handleLessonsResourcesPost).Methods("POST")

	subrouter.HandleFunc("/resources/{rid}",
		handleLessonsResourceGet).Methods("GET")
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

	if !(authContext.Account.ID == lessonRequest.StudentID || authContext.Account.ID == lessonRequest.TutorID) {
		restError(w, r, errors.New("only allowed operate on a lesson that has your account as a participant"), http.StatusBadRequest)
		return
	}

	student, err := services.ReadAccountByID(lessonRequest.StudentID, nil)
	if err != nil {
		restError(w, r, err, http.StatusBadRequest)
		return
	}

	tutor, err := services.ReadAccountByID(lessonRequest.TutorID, nil)
	if err != nil {
		restError(w, r, err, http.StatusBadRequest)
		return
	}

	err = services.CreateLesson(authContext.Account, student, tutor, lessonRequest.StartTime, lessonRequest.LessonDetail)
	if err != nil {
		restError(w, r, err, http.StatusBadRequest)
		return
	}
}

func handleLessonsRequestStageChange(w http.ResponseWriter, r *http.Request, stage services.LessonRequestStage) {
	var stageRequest LessonStageChangeDTO
	if !ParseBody(w, r, &stageRequest) {
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

	err = lesson.UpdateRequestStageByAccount(authContext.Account, stage, stageRequest.LessonDetail)
	if err != nil {
		restError(w, r, err, http.StatusBadRequest)
		return
	}
}

func handleLessonsRequestStageChangeClosure(stage services.LessonRequestStage) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		handleLessonsRequestStageChange(w, r, stage)
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
