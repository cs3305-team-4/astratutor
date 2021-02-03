package routes

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"mime"
	"net/http"
	"strings"
	"time"

	"github.com/cs3305-team-4/api/pkg/services"
	"github.com/google/uuid"
	"github.com/gorilla/mux"
)

// LessonDTO represents an existing lesson
type LessonDTO struct {
	// Time of the lesson
	TimeStarts time.Time `json:"time"`

	// Tutor of the lesson
	TutorID uuid.UUID `json:"tutor_id"`

	// Student of the lesson
	StudentID uuid.UUID `json:"student_id"`

	// LessonDetail contains notes about what the student needs out of this lesson
	LessonDetail string `json:"lesson_detail"`

	// RequestStagedetermines what state of request the lesson is in
	RequestStage services.LessonRequestStage `json:"request_stage"`

	// RequestStageDetail contains a string related to the current request state
	RequestStageDetail string `json:"request_stage_detail"`

	// RequestStageChanger contains a reference to the account of the person who last changed the state of the lesson
	RequestStageChangerID uuid.UUID `json:"request_stage_change_id"`

	Resources []LessonResourceDTO `json:"resources"`
}

type LessonResourceDTO struct {
	Name string `json:"name"`
	MIME string `json:"mime"`
}

type LessonResourceFileDTO struct {
	Name       string `json:"name"`
	MIME       string `json:"mime"`
	Base64Data string `json:"base64_data"`
}

// LessonRequestDTO represents a lesson that was first requested by an account
type LessonRequestDTO struct {
	// Time of the lesson
	TimeStarts time.Time `json:"time"`

	// The ID of who needs the lesson, can be a teacher or student
	// i.e if they request the tutors lesson create endpoint, it'l expect RequesterID to be a student
	RequesterID uuid.UUID `json:"requester_id"`

	LessonDetail string `json:"lesson_detail"`
}

// LessonRequesStagetDTO represents a change in request stage for a lesson
// i.e confirmed, expired, etc
type LessonRequestStageChangeDTO struct {
	// Time of the lesson
	TimeStarts time.Time `json:"time"`

	// The ID of who requested the lesson, can be a teacher or student
	// i.e if they request the tutors lesson create endpoint, it'l expect RequesterID to be a student
	RequesterID uuid.UUID `json:"tutor_id"`

	LessonDetail string `json:"lesson_detail"`
}

func dtoFromLesson(l *services.Lesson) *LessonDTO {
	return &LessonDTO{
		TimeStarts:            l.TimeStarts,
		TutorID:               l.Tutor.ID,
		StudentID:             l.Student.ID,
		LessonDetail:          l.LessonDetail,
		RequestStage:          l.RequestStage,
		RequestStageDetail:    l.RequestStageDetail,
		RequestStageChangerID: l.RequestStageChanger.ID,
		Resources:             []LessonResourceDTO{},
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
	subrouter.PathPrefix("/{uuid}").HandlerFunc(handleLessonsGet).Methods("GET")

	subrouter.PathPrefix("/{uuid}/accept").HandlerFunc(
		handleLessonsRequestStageChangeClosure(services.Accepted),
	).Methods("POST")

	subrouter.PathPrefix("/{uuid}/deny").HandlerFunc(
		handleLessonsRequestStageChangeClosure(services.Denied),
	).Methods("POST")

	subrouter.PathPrefix("/{uuid}/cancel").HandlerFunc(
		handleLessonsRequestStageChangeClosure(services.Cancelled),
	).Methods("POST")

	subrouter.PathPrefix("/{uuid}/complete").HandlerFunc(
		handleLessonsRequestStageChangeClosure(services.Completed),
	).Methods("POST")

	subrouter.PathPrefix("/{uuid}/resources").HandlerFunc(handleLessonsResourcesPost).Methods("POST")
	subrouter.PathPrefix("/{uuid}/resources/{rid}").HandlerFunc(handleLessonsResourceGet).Methods("GET")
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

func handleLessonsRequestStageChange(w http.ResponseWriter, r *http.Request, stage services.LessonRequestStage) {
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

	err = lesson.ChangeLessonRequestStage(nil, stage)
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

	err = lesson.ReadResourceByID(rid)
	if err != nil {
		restError(w, r, err, http.StatusBadRequest)
		return
	}

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

	r.ParseMultipartForm(16777216) // 16mb

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

	b64 := base64.StdEncoding.EncodeToString(buf.Bytes())

	err = lesson.CreateResource(filenameExclExt, mimeType, b64)
	if err != nil {
		restError(w, r, err, http.StatusBadRequest)
		return
	}
}
