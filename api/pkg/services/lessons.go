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
	Accepted = "accepted"

	// One party has suggested a reschedule
	SuggestReschedule = "reschedule"

	// The lesson (request) has been denied
	Denied = "denied"

	// The lesson has been cancelled
	Cancelled = "cancelled"

	// Lesson completed
	Completed = "completed"

	NoShowStudent = "no-show-student"

	NoShowTeacher = "no-show-teacher"
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
