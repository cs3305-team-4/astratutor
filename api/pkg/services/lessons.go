package services

import (
	"errors"
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

	// The lesson request has been accepted by the reciever party and now the student needs to pay
	PaymentRequired LessonRequestStage = "payment-required"

	// The lesson (request) has been denied
	Denied LessonRequestStage = "denied"

	// The student has paid for the lesson and the lesson is now scheduled
	Scheduled LessonRequestStage = "scheduled"

	// The lesson has been rescheduled
	Rescheduled LessonRequestStage = "rescheduled"

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
	StartTime time.Time

	EndTime time.Time

	Student   Account `gorm:"foreignKey:StudentID"`
	StudentID uuid.UUID

	SubjectTaught   SubjectTaught `gorm:"foreignKey:SubjectTaughtID"`
	SubjectTaughtID uuid.UUID

	PaymentIntentID string

	// Paid status, used to minimize the number of API Calls towards stripe
	Paid bool

	// PaidOut - has the tutor been paid for this lesson
	PaidOut bool

	// DatePaidOut - date the lesson was paid out
	DatePaidOut *time.Time

	// Approximate time the lesson was paid for
	DatePaid *time.Time

	// PriceAmount is the cost of the lesson
	PriceAmount int64

	// PayoutAmount is the amount the tutor will earn on this lesson
	PayoutAmount int64

	// Refunded status
	Refunded bool

	Tutor   Account `gorm:"foreignKey:TutorID"`
	TutorID uuid.UUID

	// LessonDetail contains notes about what the student needs out of this lesson
	LessonDetail string

	// RequestStagedetermines what state of request the lesson is in
	RequestStage LessonRequestStage

	// RequestStageDetail contains a string related to the current request state
	RequestStageDetail string

	// Requester of the lesson, should be a student or a tutor
	Requester   Account `gorm:"foreignKey:RequesterID"`
	RequesterID uuid.UUID

	// RequestStageChanger contains a reference to the account of the person who last changed the state of the lesson
	RequestStageChanger   Account `gorm:"foreignKey:RequestStageChangerID"`
	RequestStageChangerID uuid.UUID

	// Resources are
	Resources []ResourceMetadata `gorm:"foreignKey:LessonID"`
}

// ResourceMetadata contains metadata about a resource
type ResourceMetadata struct {
	database.Model
	LessonID       uuid.UUID `gorm:"type:uuid"`
	Name           string
	MIME           string
	ResourceData   ResourceData `gorm:"foreignKey:ResourceDataID"`
	ResourceDataID uuid.UUID
}

// ResourceData contains the data of an actual resource
type ResourceData struct {
	database.Model
	Data []byte `gorm:"type:bytea"`
}

// LessonAtTime returns true if the account has a lesson at that time
func LessonAtTime(acc *Account, startTime time.Time, endTime time.Time) (bool, error) {
	db, err := database.Open()
	if err != nil {
		return false, err
	}

	var lessons []Lesson

	result := db.Where(
		"(student_id = ? OR tutor_id = ?) AND (end_time > ? AND start_time < ?)",
		acc.ID, acc.ID, startTime, endTime,
	).Find(&lessons)

	if result.Error != nil {
		return false, result.Error
	}

	if len(lessons) > 1 {
		return true, nil
	}

	return false, nil
}

func RequestLesson(requester *Account, student *Account, subjectTaught *SubjectTaught, startTime time.Time, lessonDetail string) error {
	if !startTime.After(time.Now()) {
		return fmt.Errorf("can't request a lesson in the past")
	}

	tutor, err := ReadAccountByID(subjectTaught.TutorID, nil)
	if err != nil {
		return err
	}

	if !(requester.ID == student.ID || requester.ID == tutor.ID) {
		return errors.New("account requesting the lesson must be involved in the lesson")
	}

	if student.Type != Student {
		return fmt.Errorf("specified student account is not a student")
	}

	if tutor.Type != Tutor {
		return fmt.Errorf("specified tutor account is not a tutor")
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

		endTime := startTime.Add(time.Minute*time.Duration(59) + time.Second*time.Duration(59))

		lat, err := LessonAtTime(student, startTime, endTime)
		if err != nil {
			tx.Rollback()
			return err
		}

		if lat == true {
			return fmt.Errorf("cannot create lesson: the student has a lesson at that time")
		}

		lat, err = LessonAtTime(tutor, startTime, endTime)
		if err != nil {
			tx.Rollback()
			return err
		}

		if lat == true {
			return fmt.Errorf("cannot create lesson: the teacher has a lesson at that time")
		}

		// First create stripe invoice for the lesson

		l := &Lesson{
			StartTime:           startTime,
			EndTime:             endTime,
			RequesterID:         requester.ID,
			Student:             *student,
			Tutor:               *tutor,
			StudentID:           student.ID,
			TutorID:             tutor.ID,
			SubjectTaught:       *subjectTaught,
			SubjectTaughtID:     subjectTaught.ID,
			Paid:                false,
			LessonDetail:        lessonDetail,
			RequestStage:        Requested,
			RequestStageDetail:  lessonDetail,
			Resources:           []ResourceMetadata{},
			RequestStageChanger: *requester,
		}

		err = l.SetupPaymentIntent()
		if err != nil {
			tx.Rollback()
			return err
		}

		err = tx.Create(l).Error

		if err != nil {
			tx.Rollback()
			return err
		}

		return nil
	})

	return err

}

func ReadLessonByID(id uuid.UUID, preloads ...string) (*Lesson, error) {
	db, err := database.Open()
	if err != nil {
		return nil, err
	}

	for _, preload := range preloads {
		db = db.Preload(preload)
	}

	var lesson Lesson
	if err := db.First(&lesson, id).Error; err != nil {
		return nil, fmt.Errorf("lesson not found")
	}

	return &lesson, nil
}

func ReadLessonsByAccountID(id uuid.UUID) ([]Lesson, error) {
	db, err := database.Open()
	if err != nil {
		return nil, err
	}

	var lessons []Lesson

	err = db.Where("student_id = ? OR tutor_id = ?", id, id).Find(&lessons).Error

	if err != nil {
		return nil, err
	}

	return lessons, nil
}

func (l *Lesson) CreateResource(name string, mime string, data []byte) error {
	db, err := database.Open()
	if err != nil {
		return err
	}

	err = db.Create(&ResourceMetadata{
		LessonID: l.ID,
		Name:     name,
		ResourceData: ResourceData{
			Data: data,
		},
	}).Error
	if err != nil {
		return err
	}

	return nil
}

// ChangeRequestStage changes the stage the lesson is at
// i.e a requester can request the lesson move from the Requested state to the Acceptd state to confirm that the lesson will take place
// func (l *Lesson) UpdateRequestStageByAccount(stageRequester *Account, newStage LessonRequestStage, detail string) error {
// 	db, err := database.Open()
// 	if err != nil {
// 		return err
// 	}

// 	err = db.Transaction(func(tx *gorm.DB) error {
// 		// re-read the lesson, stops data races
// 		lesson, err := ReadLessonByID(l.ID)
// 		if err != nil {
// 			tx.Rollback()
// 			return err
// 		}

// 		// Make a decision based off the current stage the lesson is at
// 		switch lesson.RequestStage {
// 		// If the lesson stage was 'requested'
// 		case Requested:
// 			switch newStage {
// 			case Scheduled:
// 				if stageRequester.ID == lesson.RequesterID {
// 					return errors.New("you can not mark a lesson as scheduled if you were the one who created the lesson")
// 				}

// 			case Denied:
// 				if stageRequester.ID == lesson.RequesterID {
// 					return errors.New("you can not deny a lesson if you were the one who created the lesson")
// 				}

// 			case Cancelled:
// 				if stageRequester.ID != lesson.RequesterID {
// 					return errors.New("only the person who requested the lesson can cancel the request")
// 				}

// 			default:
// 				return fmt.Errorf("unsupported stage %s from %s", newStage, lesson.RequestStage)
// 			}

// 		case Scheduled:
// 			switch newStage {
// 			case Cancelled:

// 			default:
// 				return fmt.Errorf("unsupported stage %s from %s", newStage, lesson.RequestStage)
// 			}

// 		default:
// 			return fmt.Errorf("unsupported stage %s from %s", newStage, lesson.RequestStage)
// 		}

// 		db.Model(&lesson).Updates(&Lesson{
// 			RequestStage:          newStage,
// 			RequestStageDetail:    detail,
// 			RequestStageChangerID: stageRequester.ID,
// 		})
// 		return nil
// 	})

// 	return err
// }

func (l *Lesson) MarkScheduled(requester *Account) error {
	db, err := database.Open()
	if err != nil {
		return err
	}

	err = l.RefreshPaidStatus()
	if err != nil {
		return err
	}

	if l.Paid == false {
		return errors.New("cannot mark a lesson as scheduled if the lesson has not been paid for by the student")
	}

	err = db.Transaction(func(tx *gorm.DB) error {
		// re-read the lesson, stops data races
		lesson, err := ReadLessonByID(l.ID)
		if err != nil {
			tx.Rollback()
			return err
		}

		// Make a decision based off the current stage the lesson is at
		switch lesson.RequestStage {
		case PaymentRequired:
			// Can only mark a lesson as scheduled if it previously required payment

		default:
			return fmt.Errorf("unsupported stage %s from %s", Scheduled, lesson.RequestStage)
		}

		db.Model(&lesson).Updates(&Lesson{
			RequestStage:          Scheduled,
			RequestStageChangerID: requester.ID,
		})
		return nil
	})

	return err
}

func (l *Lesson) MarkPaymentRequired(acceptor *Account) error {
	db, err := database.Open()
	if err != nil {
		return err
	}

	err = db.Transaction(func(tx *gorm.DB) error {
		// re-read the lesson, stops data races
		lesson, err := ReadLessonByID(l.ID)
		if err != nil {
			tx.Rollback()
			return err
		}

		// Make a decision based off the current stage the lesson is at
		switch lesson.RequestStage {
		case Requested:
			if acceptor.ID == lesson.RequesterID {
				return errors.New("you can not mark a lesson as scheduled if you were the one who created the lesson")
			}

		case Rescheduled:
			if acceptor.ID == lesson.RequestStageChangerID {
				return errors.New("you can not mark a lesson as scheduled if you were the one who rescheduled the lesson")
			}

		default:
			return fmt.Errorf("unsupported stage %s from %s", Scheduled, lesson.RequestStage)
		}

		db.Model(&lesson).Updates(&Lesson{
			RequestStage:          PaymentRequired,
			RequestStageChangerID: acceptor.ID,
		})
		return nil
	})

	return err
}

func (l *Lesson) MarkDenied(denier *Account, reason string) error {
	db, err := database.Open()
	if err != nil {
		return err
	}

	err = db.Transaction(func(tx *gorm.DB) error {
		// re-read the lesson, stops data races
		lesson, err := ReadLessonByID(l.ID)
		if err != nil {
			tx.Rollback()
			return err
		}

		// Make a decision based off the current stage the lesson is at
		switch lesson.RequestStage {
		case Requested:
			if denier.ID == lesson.RequesterID {
				return errors.New("you can not mark a lesson as denied if you were the one who created the lesson")
			}

		case Rescheduled:
			if denier.ID == lesson.RequestStageChangerID {
				return errors.New("you can not mark a lesson as denied if you were the one who rescheduled the lesson")
			}

		default:
			return fmt.Errorf("unsupported stage %s from %s", Denied, lesson.RequestStage)
		}

		db.Model(&lesson).Updates(&Lesson{
			RequestStage:          Denied,
			RequestStageDetail:    reason,
			RequestStageChangerID: denier.ID,
		})
		return nil
	})

	return err
}

func (l *Lesson) MarkCancelled(cancelee *Account, reason string) error {
	db, err := database.Open()
	if err != nil {
		return err
	}

	err = db.Transaction(func(tx *gorm.DB) error {
		// re-read the lesson, stops data races
		lesson, err := ReadLessonByID(l.ID)
		if err != nil {
			tx.Rollback()
			return err
		}

		// Make a decision based off the current stage the lesson is at
		switch lesson.RequestStage {
		case Scheduled:

		case Requested:
			if cancelee.ID != lesson.RequesterID {
				return errors.New("only the person who requested the lesson can cancel the request")
			}

		case PaymentRequired:
			if cancelee.ID != lesson.RequesterID {
				return errors.New("only the person who requested the lesson can cancel the request")
			}

		case Rescheduled:
			if cancelee.ID != lesson.RequestStageChangerID {
				return errors.New("only the person who rescheduled the lesson can cancel the request")
			}

		default:
			return fmt.Errorf("unsupported stage %s from %s", Cancelled, lesson.RequestStage)
		}

		if l.Paid == true {
			err = l.Refund()
			if err != nil {
				tx.Rollback()
				return err
			}
		}

		db.Model(&lesson).Updates(&Lesson{
			RequestStage:          Cancelled,
			RequestStageDetail:    reason,
			RequestStageChangerID: cancelee.ID,
		})
		return nil
	})

	return err
}

func (l *Lesson) MarkRescheduled(reschedulee *Account, newTime time.Time, reason string) error {
	db, err := database.Open()
	if err != nil {
		return err
	}

	if !newTime.After(time.Now()) {
		return fmt.Errorf("can't reschedule a lesson to the past")
	}

	err = db.Transaction(func(tx *gorm.DB) error {
		// re-read the lesson, stops data races
		lesson, err := ReadLessonByID(l.ID, "Student", "Tutor")
		if err != nil {
			tx.Rollback()
			return err
		}

		// Make a decision based off the current stage the lesson is at
		switch lesson.RequestStage {
		case Scheduled:
			break
		case Requested:
			break
		case Rescheduled:
			break

		default:
			return fmt.Errorf("unsupported stage %s from %s", Rescheduled, lesson.RequestStage)
		}

		endTime := newTime.Add(time.Minute*time.Duration(59) + time.Second*time.Duration(59))

		lat, err := LessonAtTime(&lesson.Student, newTime, endTime)
		if err != nil {
			tx.Rollback()
			return err
		}

		if lat == true {
			return fmt.Errorf("cannot create lesson: the student has a lesson at that time")
		}

		lat, err = LessonAtTime(&lesson.Tutor, newTime, endTime)
		if err != nil {
			tx.Rollback()
			return err
		}

		if lat == true {
			return fmt.Errorf("cannot create lesson: the teacher has a lesson at that time")
		}

		db.Model(&lesson).Updates(&Lesson{
			StartTime:             newTime,
			EndTime:               endTime,
			RequestStage:          Rescheduled,
			RequestStageDetail:    reason,
			RequestStageChangerID: reschedulee.ID,
		})
		return nil
	})

	return err
}

func (l *Lesson) ReadResourceByID(rid uuid.UUID) (*ResourceMetadata, error) {
	db, err := database.Open()
	if err != nil {
		return nil, err
	}

	var resource ResourceMetadata
	err = db.First(&ResourceMetadata{
		Model: database.Model{
			ID: rid,
		},
	}).Find(&resource).Error

	if err != nil {
		return nil, err
	}

	return &resource, nil
}

func (r *ResourceMetadata) GetData() ([]byte, error) {
	db, err := database.Open()
	if err != nil {
		return nil, err
	}

	var resourceData ResourceData
	err = db.First(&resourceData, r.ResourceDataID).Error
	if err != nil {
		return nil, err
	}

	return resourceData.Data, nil
}

func HaveCompletedLesson(student uuid.UUID, tutor uuid.UUID) (bool, error) {
	db, err := database.Open()
	if err != nil {
		return false, err
	}

	res := db.Model(&Lesson{}).Where(&Lesson{
		StudentID:    student,
		TutorID:      tutor,
		RequestStage: Completed,
	}).First(&Lesson{})

	if res.RowsAffected == 0 {
		return false, ReviewErrorNoCompletedLesson
	}

	return true, nil
}
