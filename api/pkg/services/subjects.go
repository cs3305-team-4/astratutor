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
	Slug  string `gorm:"unique;not null;"`
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

// returns a subject when given a subjects slug
func GetSubjectBySlug(slug string, db *gorm.DB) (*Subject, error) {
	if db == nil {
		var err error
		db, err = database.Open()
		if err != nil {
			return nil, err
		}
	}
	var subject Subject
	var err error
	err = db.Where(&Subject{Slug: slug}).Find(&subject).Error

	if err != nil {
		return nil, err
	}

	return &subject, nil
}

//returns a subject when given its ID
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

	if sid == uuid.MustParse("00000000-0000-0000-0000-000000000000") {
		var err error
		return nil, err
	}

	if db == nil {
		var err error
		db, err = database.Open()
		if err != nil {
			return nil, err
		}
	}
	var subjectTaught []SubjectTaught
	var err error
	err = db.Where(&SubjectTaught{SubjectID: sid}).Find(&subjectTaught).Error
	if err != nil {
		return nil, err
	}
	return subjectTaught, nil
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

//Returns subjectTaught for specific tutors using their ID
func GetSubjectsByTutorID(tid uuid.UUID, db *gorm.DB) ([]SubjectTaught, error) {
	if tid == uuid.MustParse("00000000-0000-0000-0000-000000000000") {
		var err error
		return nil, err
	}

	if db == nil {
		var err error
		db, err = database.Open()
		if err != nil {
			return nil, err
		}
	}
	var subjectTaught []SubjectTaught
	return subjectTaught, db.Where(&SubjectTaught{TutorID: tid}).Find(&subjectTaught).Error
}

//creats a StudentTaught based on the subject and tutor with a set price description.
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

//updates the price of a subjecttaught by the sid
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

//updates the description of a subjecttaught by the sid
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

//function used to create data for tests
func CreateSubjectTestAccounts() error {
	db, err := database.Open()
	if err != nil {
		return err
	}
	db.Create(&Subject{Name: "English", Slug: "English"})
	db.Create(&Subject{Name: "Math", Slug: "Math"})
	hash, err := NewPasswordHash("grindshub")
	if err != nil {
		return err
	}

	John := &Account{Model: database.Model{ID: uuid.MustParse("22222222-2222-2222-2222-222222222222")},
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
	}

	french, err := GetSubjectBySlug("English", nil)
	teachSubject(french, John, "I teach Emglish", 67, nil)

	if err != nil {
		return err
	}

	english, err := GetSubjectBySlug("Math", nil)
	teachSubject(english, John, "I teach maths", 59, nil)
	if err != nil {
		return err
	}

	return nil
}
