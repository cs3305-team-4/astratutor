package services

import (
	"fmt"
	"time"

	"github.com/cs3305-team-4/api/pkg/database"
	"github.com/google/uuid"
	"gorm.io/gorm"
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

	Expired LessonRequestState = "expired"
)

// Lesson contains information about a lesson
type Lesson struct {
	database.Model

	// Time of the lesson
	TimeStarts time.Time

	// Length of the lesson time
	// TimeStar

	// Tutor of the lesson
	Tutor Account `gorm:"foreignKey:ID"`

	// Student of the lesson
	Student Account `gorm:"foreignKey:ID"`

	// LessonDetail contains notes about what the student needs out of this lesson
	LessonDetail string

	// RequestState determines what state of request the lesson is in
	RequestState LessonRequestState

	// RequestStateDetail contains a string related to the current request state
	RequestStateDetail string

	// RequestStateChanger contains a reference to the account of the person who last changed the state of the lesson
	RequestStateChanger Account `gorm:"foreignKey:ID"`

	// Resources are
	Resources []Resource `gorm:"foreignKey:LessonID"`

	/*Subject Subject `gorm:"foreignKey:ID"`*/
}

type Resource struct {
	database.Model
	LessonID   uuid.UUID `gorm:"type:uuid"`
	Name       string
	Base64Data string
}

// LessonAtTime returns true if the account has a lesson at that time
func LessonAtTime(acc Account, t time.Time) (bool, error) {
	db, err := database.Open()
	if err != nil {
		return false, err
	}

	var lessons []Lesson
	db.Find(&lessons, "id = ? AND time_start BETWEEN ? AND ?", acc.ID, t, t.Add(time.Minute*time.Duration(60)))

	if len(lessons) > 0 {
		return true, nil
	}

	return false, nil
}

func CreateLesson(student Account, tutor Account, t time.Time /*subject Subject*/, lessonDetail string) error {
	if !student.IsStudent() {
		return fmt.Errorf("account specified as student is not a student!")
	}

	if !tutor.IsTutor() {
		return fmt.Errorf("account specified as student is not a student!")
	}

	if !t.After(time.Now()) {
		return fmt.Errorf("can't request a lesson in the past")
	}

	db, err := database.Open()
	if err != nil {
		return err
	}

	err = db.Transaction(func(tx *gorm.DB) error {
		err = tx.Exec(`set transaction isolation level repeatable read`).Error
		if err != nil {
			tx.Rollback()
			return err
		}

		lat, err := LessonAtTime(student, t)
		if err != nil {
			tx.Rollback()
			return err
		}

		if lat == true {
			return fmt.Errorf("cannot create lesson: the student has a lesson at that time")
		}

		lat, err = LessonAtTime(tutor, t)
		if err != nil {
			tx.Rollback()
			return err
		}

		if lat == true {
			return fmt.Errorf("cannot create lesson: the tutor has a lesson at that time")
		}

		err = tx.Create(&Lesson{
			TimeStarts:          t,
			Student:             student,
			Tutor:               tutor,
			LessonDetail:        lessonDetail,
			RequestState:        Requested,
			RequestStateDetail:  fmt.Sprintf("%s %s requested a lesson", student.Profile.FirstName, student.Profile.LastName),
			Resources:           []Resource{},
			RequestStateChanger: student,
		}).Error

		if err != nil {
			tx.Rollback()
			return err
		}
	})

	return err

}

func ReadLessonByID(id uuid.UUID) (*Lesson, error) {
	db, err := database.Open()
	if err != nil {
		return nil, err
	}

	var lesson Lesson
	if err := db.First(&lesson, id).Error; err != nil {
		return nil, fmt.Errorf("Lesson not found")
	}

	return &lesson, nil
}

func ReadLessonByAccount(tutor Account) (*Lesson, error) {
	db, err := database.Open()
	if err != nil {
		return nil, err
	}

	var lesson Lesson
	if err := db.First(&lesson, id).Error; err != nil {
		return nil, fmt.Errorf("Lesson not found")
	}

	return &lesson, nil
}

func (l *Lesson) CreateLessonResource(name string, data string) error {
	db, err := database.Open()
	if err != nil {
		return err
	}

	err = db.Create(&Resource{
		LessonID:   l.ID,
		Name:       name,
		Base64Data: data,
	}).Error
	if err != nil {
		return err
	}
}

func (l *Resource) DeleteResource() {

}

// func (l *Lesson) ChangeLessonRequestState(requestee Account, newState LessonRequestState) error {
// 	db, err := database.Open()
// 	if err != nil {
// 		return err
// 	}

// 	err := db.Transaction(func(tx *gorm.DB) error {
// 		err = tx.Exec(`set transaction isolation level repeatable read`).Error
// 		if err != nil {
// 			tx.Rollback()
// 			return err
// 		}

// 		// re-read the lesson, stops data races
// 		lesson, err := GetLessonByID(l)
// 		if err != nil {
// 			tx.Rollback()
// 			return err
// 		}

// 		switch newState {

// 		}
// 	})

// 	return err
// }
