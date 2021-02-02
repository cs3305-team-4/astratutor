package services

import (
	"time"

	"github.com/cs3305-team-4/api/pkg/database"
	"github.com/google/uuid"
)

type LessonRequestState string

const (
	// The lesson has been requested by the sending party
	Requested LessonRequestState = "requested"

	// The lesson has been accepted by the reciever party
	Accepted LessonRequestState = "accepted"

	// One party has suggested a reschedule
	SuggestReschedule LessonRequestState = "reschedule"

	// The lesson (request) has been denied
	Denied LessonRequestState = "denied"

	// The lesson has been cancelled
	Cancelled LessonRequestState = "cancelled"

	// Lesson completed
	Completed LessonRequestState = "completed"

	NoShowStudent LessonRequestState = "no-show-student"

	NoShowTeacher LessonRequestState = "no-show-teacher"
)

// Lesson contains information about a lesson
type Lesson struct {
	database.Model

	// ScheduledTime of the lesson
	ScheduledTime time.Time

	// Tutor of the lesson
	Tutor Account `gorm:"foreignKey:ID"`

	// Student of the lesson
	Student Account `gorm:"foreignKey:ID"`

	// RequestState determines what state of request the lesson is in
	RequestState LessonRequestState

	// RequestStateDetail contains a string related to the current request state
	RequestStateDetail string

	// RequestStateChanger contains a reference to the account of the person who last changed the state of the lesson
	RequestStateChanger Account `gorm:"foreignKey:ID"`

	// Resources are
	Resources []Resource `gorm:"foreignKey:LessonID"`
}

type Resource struct {
	database.Model
	LessonID   uuid.UUID `gorm:"type:uuid"`
	Base64Data string
}
