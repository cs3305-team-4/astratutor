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
	Slug string `gorm:"unique;not null;"`
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

	return db.Create(&Subject{Name: name, Slug: slug}).Error
}

//Contains information on
type SubjectTaught struct {
	database.Model

	Subject   Subject `gorm:"foreignKey:SubjectID"`
	SubjectID uuid.UUID

	TutorProfile   Profile `gorm:"foreignKey:TutorProfileID"`
	TutorProfileID uuid.UUID

	Description string  `gorm:"not null;"`
	Price       float32 `gorn:"not null;"`
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

// returns a subjects when given a list of subject slugs
func GetSubjectsBySlugs(slugs []string, db *gorm.DB) (*[]Subject, error) {
	if db == nil {
		var err error
		db, err = database.Open()
		if err != nil {
			return nil, err
		}
	}
	var subjects []Subject
	var err error
	err = db.Where("slug IN (?)", slugs).Find(&subjects).Error

	if err != nil {
		return nil, err
	}

	return &subjects, nil
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
func GetTutorsBySubjectIDs(sid uuid.UUID, db *gorm.DB, preloads ...string) ([]SubjectTaught, error) {

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

	for _, preload := range preloads {
		db = db.Preload(preload)
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
func GetAllTutors(db *gorm.DB, preloads ...string) ([]SubjectTaught, error) {
	if db == nil {
		var err error
		db, err = database.Open()
		if err != nil {
			return nil, err
		}
	}

	for _, preload := range preloads {
		db = db.Preload(preload)
	}

	var subjectTaught []SubjectTaught
	return subjectTaught, db.Find(&subjectTaught).Error
}

//Returns subjectTaught for specific tutors using their ID
func GetSubjectsTaughtByTutorID(tid uuid.UUID, db *gorm.DB) ([]SubjectTaught, error) {
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
	return subjectTaught, db.Where(&SubjectTaught{TutorProfileID: tid}).Find(&subjectTaught).Error
}

//creats a StudentTaught based on the subject and tutor with a set price description.
func TeachSubject(subject *Subject, tutor *Account, description string, price float32, db *gorm.DB) error {
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
			Subject:      *subject,
			TutorProfile: *tutor.Profile,
			Price:        price,
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
	db.Create(&Subject{Name: "French", Slug: "french"})
	db.Create(&Subject{Name: "Irish", Slug: "irish"})
	hash, err := NewPasswordHash("grindshub")
	if err != nil {
		return err
	}

	John := &Account{Model: database.Model{ID: uuid.MustParse("11111111-1111-1111-1111-111111111111")},
		Email:         "tutor4@grindshub.localhost",
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
	}

	irish, err := GetSubjectBySlug("irish", nil)
	TeachSubject(irish, John, "I teach irish", 67, nil)

	if err != nil {
		return err
	}

	french, err := GetSubjectBySlug("french", nil)
	TeachSubject(french, John, "I teach French", 59, nil)
	if err != nil {
		return err
	}

	return nil
}
