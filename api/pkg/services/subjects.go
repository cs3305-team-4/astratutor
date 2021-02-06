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
	Name string `gorm:"unique;not null;"`
	Slug string
}

type SubjectTaughtError string

func (e SubjectTaughtError) Error() string {
	return string(e)
}

const (
	SubjectTaughtErrorDoesNotExist SubjectTaughtError = "This tutor subject relation does not exist"
)

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

	return db.Create(&SubjectTaught{Subject: *subject, Tutor: *tutor, Price: price, Description: description}).Error

}

func updateCost(stid uuid.UUID, price uint, db *gorm.DB) (*SubjectTaught, error) {
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

func updateDescription(stid uuid.UUID, description string, db *gorm.DB) (*SubjectTaught, error) {
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

func CreateSubjectTestAccounts() error {
	db, err := database.Open()
	db.Create(&Subject{Name: "English"})
	db.Create(&Subject{Name: "Maths"})
	if err != nil {
		return err
	}

	hash, err := NewPasswordHash("grindshub")
	if err != nil {
		return err
	}

	db.FirstOrCreate(&Account{
		Model: database.Model{
			ID: uuid.MustParse("deadbeef-cafe-badd-c0de-facadebadbad"),
		},
		Email:         "tutor@grindshub.localhost",
		EmailVerified: true,
		Type:          Tutor,
		Suspended:     false,
		PasswordHash:  *hash,
		Profile: &Profile{
			Avatar:         "",
			Slug:           "Mike-tutor",
			FirstName:      "Mike",
			LastName:       "Tutor",
			City:           "Cork",
			Country:        "Ireland",
			Description:    "A tutor",
			Qualifications: []Qualification{},
			WorkExperience: []WorkExperience{},
		},
	})
	if err != nil {
		return err
	}

	db.FirstOrCreate(&Account{
		Model: database.Model{
			ID: uuid.MustParse("deadbeef-cafe-badd-c0de-facadebadbad"),
		},
		Email:         "tutor2@grindshub.localhost",
		EmailVerified: true,
		Type:          Tutor,
		Suspended:     false,
		PasswordHash:  *hash,
		Profile: &Profile{
			Avatar:         "",
			Slug:           "john-tutor",
			FirstName:      "John",
			LastName:       "Tutor",
			City:           "Cork",
			Country:        "Ireland",
			Description:    "A tutor",
			Qualifications: []Qualification{},
			WorkExperience: []WorkExperience{},
		},
	})
	if err != nil {
		return err
	}

	John, err := ReadAccountByEmail("tutor@grindshub.localhost2", nil)
	Mike, err := ReadAccountByEmail("tutor@grindshub.localhost", nil)
	English, err := GetSubjectByName("English", nil)
	Maths, err := GetSubjectByName("Maths", nil)
	teachSubject(Maths, John, "John teaches stuff", 70, nil)
	teachSubject(English, Mike, "Mike teaches better stuff", 75, nil)

	return nil
}
