package services

import (
	"time"

	"github.com/cs3305-team-4/api/pkg/db"
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

type Lesson struct {
	db.Model
	Tutor   Account   `gorm:"foreignKey:TutorID"`
	TutorID uuid.UUID `gorm:"type:uuid"`

	Student   Account   `gorm:"foreignKey:StudentID"`
	StudentID uuid.UUID `gorm:"type:uuid"`

	SubjectTaught SubjectTaught // `gorm:"foreignKey:ID"`

	Time time.Time

	RequestState        LessonRequestState
	RequestStateDetail  string
	RequestStateChanger AccountType

	Resources []Resource `gorm:"foreignKey:LessonID"`
}

type Resource struct {
	db.Model
	LessonID   uuid.UUID `gorm:"type:uuid"`
	Base64Data string
}
