package services

import (
	"fmt"
	"time"

	"github.com/cs3305-team-4/api/pkg/database"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type LessonRequestStage string

const (
	// The lesson has been requested by the sending party
	Requested LessonRequestStage = "requested"

	// The lesson has been accepted by the reciever party
	Accepted LessonRequestStage = "accepted"

	// One party has suggested a reschedule
	// SuggestReschedule LessonRequestStage = "reschedule"

	// The lesson (request) has been denied
	Denied LessonRequestStage = "denied"

	// The lesson has been cancelled
	Cancelled LessonRequestStage = "cancelled"

	// Lesson completed
	Completed LessonRequestStage = "completed"

	// No show from student
	NoShowStudent LessonRequestStage = "no-show-student"

	// No show from tutor
	NoShowTutor LessonRequestStage = "no-show-tutor"

	// Lesson request expired
	Expired LessonRequestStage = "expired"
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

	// RequestStagedetermines what state of request the lesson is in
	RequestStage LessonRequestStage

	// RequestStageDetail contains a string related to the current request state
	RequestStageDetail string

	// RequestStageChanger contains a reference to the account of the person who last changed the state of the lesson
	RequestStageChanger Account `gorm:"foreignKey:ID"`

	// Resources are
	Resources []Resource `gorm:"foreignKey:LessonID"`

	/*Subject Subject `gorm:"foreignKey:ID"`*/
}

type Resource struct {
	database.Model
	LessonID   uuid.UUID `gorm:"type:uuid"`
	Name       string
	MIME       string
	Base64Data string `gorm:"type:text"`
}

// LessonAtTime returns true if the account has a lesson at that time
func LessonAtTime(acc *Account, t time.Time) (bool, error) {
	db, err := database.Open()
	if err != nil {
		return false, err
	}

	var lessons []Lesson
	result := db.Find(&lessons, "id = ? AND time_start BETWEEN ? AND ?", acc.ID, t, t.Add(time.Minute*time.Duration(60)))

	if result.Error != nil {
		return false, nil
	}

	if result.RowsAffected > 0 {
		return true, nil
	}

	return false, nil
}

func CreateLesson(student *Account, tutor *Account, t time.Time /*subject *Subject*/, lessonDetail string) error {
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
			Student:             *student,
			Tutor:               *tutor,
			LessonDetail:        lessonDetail,
			RequestStage:        Requested,
			RequestStageDetail:  fmt.Sprintf("%s %s requested a lesson", student.Profile.FirstName, student.Profile.LastName),
			Resources:           []Resource{},
			RequestStageChanger: *student,
		}).Error

		if err != nil {
			tx.Rollback()
			return err
		}

		return nil
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

func ReadLessonsByTutorID(id uuid.UUID) ([]Lesson, error) {
	db, err := database.Open()
	if err != nil {
		return nil, err
	}

	var lessons []Lesson

	err = db.Where(&Lesson{
		Tutor: Account{
			Model: database.Model{
				ID: id,
			},
		},
	}).Find(&lessons).Error

	if err != nil {
		return nil, err
	}

	return lessons, nil
}

func ReadLessonsByStudentID(id uuid.UUID) ([]Lesson, error) {
	db, err := database.Open()
	if err != nil {
		return nil, err
	}

	var lessons []Lesson

	err = db.Where(&Lesson{
		Student: Account{
			Model: database.Model{
				ID: id,
			},
		},
	}).Find(&lessons).Error

	if err != nil {
		return nil, err
	}

	return lessons, nil
}

func (l *Lesson) CreateResource(name string, mime string, base64Data string) error {
	db, err := database.Open()
	if err != nil {
		return err
	}

	err = db.Create(&Resource{
		LessonID:   l.ID,
		Name:       name,
		Base64Data: base64Data,
	}).Error
	if err != nil {
		return err
	}

	return nil
}

func (r *Resource) DeleteResource() error {
	db, err := database.Open()
	if err != nil {
		return err
	}

	err = db.Delete(&r).Error
	if err != nil {
		return err
	}

	return nil
}

func (l *Lesson) ChangeLessonRequestStage(requestee *Account, newState LessonRequestStage) error {
	db, err := database.Open()
	if err != nil {
		return err
	}

	err = db.Transaction(func(tx *gorm.DB) error {
		// re-read the lesson, stops data races
		// lesson, err := ReadLessonByID(l.ID)
		// if err != nil {
		// 	tx.Rollback()
		// 	return err
		// }

		// switch newState {

		// }

		return nil
	})

	return err
}

func (l *Lesson) ReadResourceByID(rid uuid.UUID) (*Resource, error) {
	db, err := database.Open()
	if err != nil {
		return nil, err
	}

	var resource Resource
	err = db.First(&Resource{
		Model: database.Model{
			ID: rid,
		},
		LessonID: l.ID,
	}).Find(&resource).Error
	if err != nil {
		return nil, err
	}

	return &resource, nil
}
