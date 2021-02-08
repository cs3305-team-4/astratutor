package services

import (
	"fmt"

	"github.com/cs3305-team-4/api/pkg/database"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

//Subject contains information about a single subject
type Subject struct {
	database.Model
	Name  string `gorm:"unique;not null;"`
	Slug  string
	Image string
}

type SubjectTaughtError string

func (e SubjectTaughtError) Error() string {
	return string(e)
}

const (
	SubjectTaughtErrorDoesNotExist SubjectTaughtError = "This tutor subject relation does not exist"
)

func CreateSubject(name string, image string, slug string, db *gorm.DB) error {
	db, err := database.Open()
	if err != nil {
		return err
	}

	return db.Create(&Subject{Name: name, Image: image, Slug: slug}).Error
}

//Contains information on
type SubjectTaught struct {
	database.Model

	//Contains Subject bring taught
	Subject   Subject `gorm:"foreignKey:SubjectID"`
	SubjectID uuid.UUID
	//
	Tutor   Account `gorm:"foreignKey:TutorID"`
	TutorID uuid.UUID
	//Price that the Tutor wishes to charge per lesson
	Price uint
	//Description given by the tutor
	Description string
}

//gets all subjects in the DB
func GetSubjects(db *gorm.DB) ([]Subject, error) {
	if db == nil {
		var err error
		db, err = database.Open()
		if err != nil {
			return nil, err
		}
	}
	var subject []Subject
	return subject, db.Find(&subject).Error
}

// returns a subject when given a subject name
func GetSubjectByName(name string, db *gorm.DB) (*Subject, error) {
	if db == nil {
		var err error
		db, err = database.Open()
		if err != nil {
			return nil, err
		}
	}
	subject := &Subject{}
	return subject, db.Where("Name = ?", name).Find(&subject).Error
}

func GetSubjectByID(id uuid.UUID, db *gorm.DB) (*Subject, error) {
	if db == nil {
		var err error
		db, err = database.Open()
		if err != nil {
			return nil, err
		}
	}
	subject := &Subject{}
	return subject, db.Where("ID = ?", id).Find(&subject).Error
}

//Quries the DB for SubjectTaught where the subject ID is a match

func GetTutorsBySubjectID(sid uuid.UUID, db *gorm.DB) ([]SubjectTaught, error) {
	if db == nil {
		var err error
		db, err = database.Open()
		if err != nil {
			return nil, err
		}
	}
	var subjectTaught []SubjectTaught
	return subjectTaught, db.Where(&SubjectTaught{Subject: Subject{Model: database.Model{ID: sid}}}).Find(&subjectTaught).Error
}

//Quries the DB for SubjectTaught where the ID matches the SubjectTaught ID
func GetSubjectTaughtByID(stid uuid.UUID, db *gorm.DB) (*SubjectTaught, error) {
	if db == nil {
		var err error
		db, err = database.Open()
		if err != nil {
			return nil, err
		}
	}
	subjectTaught := &SubjectTaught{}
	return subjectTaught, db.Where(&SubjectTaught{Model: database.Model{ID: stid}}).Find(&subjectTaught).Error
}

//Returns all subjectTaught
func GetAllTutors(db *gorm.DB) ([]SubjectTaught, error) {
	if db == nil {
		var err error
		db, err = database.Open()
		if err != nil {
			return nil, err
		}
	}
	var subjectTaught []SubjectTaught
	return subjectTaught, db.Find(&subjectTaught).Error
}

//Returns subjectTaught for specific Tutors

func GetSubjectsByTutorID(tid uuid.UUID, db *gorm.DB) ([]SubjectTaught, error) {
	if db == nil {
		var err error
		db, err = database.Open()
		if err != nil {
			return nil, err
		}
	}
	var subjectTaught []SubjectTaught
	return subjectTaught, db.Where(&SubjectTaught{Tutor: Account{Model: database.Model{ID: tid}}}).Find(&subjectTaught).Error
}

func teachSubject(subject *Subject, tutor *Account, description string, price uint, db *gorm.DB) error {
	db, err := database.Open()
	if err != nil {
		return err
	}
	if !tutor.IsTutor() {
		return fmt.Errorf("account specified is not a tutor account and thus cannot teach")
	}
	if subject == nil {
		return fmt.Errorf("there must be a subject to teach")
	}

	err = db.Transaction(func(tx *gorm.DB) error {
		err = tx.Exec(`set transaction isolation level repeatable read`).Error
		if err != nil {
			tx.Rollback()
			return err
		}

		err = tx.Create(&SubjectTaught{
			Subject:     *subject,
			Tutor:       *tutor,
			Description: description,
			Price:       price,
		}).Error

		if err != nil {
			tx.Rollback()
			return err
		}

		return nil
	})

	return err

}

func UpdateCost(stid uuid.UUID, price uint, db *gorm.DB) (*SubjectTaught, error) {
	db, err := database.Open()
	if err != nil {
		return nil, err
	}
	var subjectTaught *SubjectTaught
	return subjectTaught, db.Transaction(func(tx *gorm.DB) error {
		dbSubjectTaught, err := GetSubjectTaughtByID(stid, nil)
		if err != nil {
			return err
		}
		if subjectTaught = dbSubjectTaught; subjectTaught == nil {
			return SubjectTaughtErrorDoesNotExist
		}
		return tx.Model(subjectTaught).Update("Price", price).Error
	})
}

func UpdateDescription(stid uuid.UUID, description string, db *gorm.DB) (*SubjectTaught, error) {
	db, err := database.Open()
	if err != nil {
		return nil, err
	}
	var subjectTaught *SubjectTaught
	return subjectTaught, db.Transaction(func(tx *gorm.DB) error {
		dbSubjectTaught, err := GetSubjectTaughtByID(stid, nil)
		if err != nil {
			return err
		}
		if subjectTaught = dbSubjectTaught; subjectTaught == nil {
			return SubjectTaughtErrorDoesNotExist
		}
		return tx.Model(subjectTaught).Update("Description", description).Error
	})
}

/*
func CreateSubjectTestAccounts() error {
	 db, err := database.Open()
	if err != nil {
		return err
	}

	hash, err := NewPasswordHash("grindshub")
	if err != nil {
		return err
	}

	english, err := GetSubjectByName("French", nil)
	teachSubject(english, &Account{Model: database.Model{ID: uuid.MustParse("deadlamb-cafe-badd-c0de-facadebadbad")},
		Email:         "tutor3@grindshub.localhost",
		EmailVerified: true,
		Type:          Tutor,
		Suspended:     false,
		PasswordHash:  *hash,
		Profile: &Profile{
			Avatar:         "",
			Slug:           "mike-tutor",
			FirstName:      "Mike",
			LastName:       "Tutor",
			City:           "Cork",
			Country:        "Ireland",
			Description:    "A tutor",
			Qualifications: []Qualification{},
			WorkExperience: []WorkExperience{},
		},
	}, "I hate gorm more", 67, nil)
	if err != nil {
		return err
	}

	return nil
}
*/
