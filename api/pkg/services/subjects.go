package services

import (
	"fmt"
	"math"
	"strings"

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

//Subject Request contains information about a requested subject and who has requested it
type SubjectRequest struct {
	database.Model
	RequesterID uuid.UUID
	Requester   Account `gorm:"foreignKey:RequesterID"`
	Name        string  `gorm:"not null;"`
	Status      SubjectRequestStatus
	Reason      string
}

type SubjectRequestStatus string

const (
	SubjectRequestApproved SubjectRequestStatus = "approved"
	SubjectRequestDenied   SubjectRequestStatus = "denied"
	SubjectRequestPending  SubjectRequestStatus = "pending"
)

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

//Contains information on a single tutor subject relation
type SubjectTaught struct {
	database.Model

	Subject   Subject `gorm:"foreignKey:SubjectID"`
	SubjectID uuid.UUID

	Tutor   Account `gorm:"foreignKey:TutorID"`
	TutorID uuid.UUID

	TutorProfile   Profile   `gorm:"foreignKey:TutorProfileID"`
	TutorProfileID uuid.UUID // Foreign key for the Profile table

	Description string `gorm:"not null;"`
	Price       int64  `gorn:"not null;"`
}

//gets all subjects in the DB
func GetSubjects(query string, db *gorm.DB) ([]Subject, error) {
	if db == nil {
		var err error
		db, err = database.Open()
		if err != nil {
			return nil, err
		}
	}
	var subject []Subject
	err := db.
		Scopes(
			Search(SearchQuery{"name", query}),
		).
		Find(&subject).
		Error
	if err != nil {
		return nil, err
	}
	return subject, nil
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
func GetTutorsBySubjectsPaginated(subjects *[]Subject, pageSize int, page int, query string, sort string, db *gorm.DB) ([]Profile, int, error) {
	if db == nil {
		var err error
		db, err = database.Open()
		if err != nil {
			return nil, 0, err
		}
	}

	var subject_ids []string
	for _, subject := range *subjects {
		subject_ids = append(subject_ids, subject.ID.String())
	}

	scopes := []func(*gorm.DB) *gorm.DB{}
	for _, q := range strings.Split(query, " ") {
		scopes = append(scopes, Search(SearchQuery{"profiles.first_name", q}, SearchQuery{"profiles.last_name", q}, SearchQuery{"profiles.country", q}, SearchQuery{"profiles.city", q}, SearchQuery{"profiles.description", q}))
	}

	var totalTutors int64
	db.Model(&SubjectTaught{}).
		Joins("JOIN profiles ON profiles.id = subject_taughts.tutor_profile_id").
		Where("subject_taughts.subject_id IN (?)", subject_ids).
		Scopes(
			scopes...,
		).
		Distinct("tutor_profile_id").Count(&totalTutors)

	scopes = append(scopes, Paginate(pageSize, page))
	asc := ""
	switch sort {
	case "low":
		asc = "asc"
	case "high":
		asc = "desc"
	}
	scopes = append(scopes, Sort("subject_taughts.price", asc, Join{
		new:     Table{"subject_taughts", "tutor_profile_id"},
		current: Table{"profiles", "id"},
	}))
	// Get tutors who are teaching subjects paginated,
	var profiles []Profile
	err := db.
		Where("profiles.id IN (?)",
			db.Model(&SubjectTaught{}).
				Where("subject_id IN (?)", subject_ids).
				Select("tutor_profile_id")).
		Preload("Subjects").Preload("Subjects.Subject").
		Scopes(
			scopes...,
		).
		Group("subject_taughts.price").
		Where("subject_taughts.subject_id IN ?", subject_ids).
		Find(&profiles).Error
	if err != nil {
		return nil, 0, err
	}

	return profiles, int(math.Ceil(float64(totalTutors) / float64(pageSize))), nil
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
	return subjectTaught, db.Preload("TutorProfile").Where(&SubjectTaught{Model: database.Model{ID: stid}}).Find(&subjectTaught).Error
}

//Returns all subjectTaught
func GetAllTutorsPaginated(db *gorm.DB, pageSize int, query string, sort string, page int) ([]Profile, int, error) {
	if db == nil {
		var err error
		db, err = database.Open()
		if err != nil {
			return nil, 0, err
		}
	}

	scopes := []func(*gorm.DB) *gorm.DB{}
	for _, q := range strings.Split(query, " ") {
		scopes = append(scopes, Search(SearchQuery{"profiles.first_name", q}, SearchQuery{"profiles.last_name", q}, SearchQuery{"profiles.country", q}, SearchQuery{"profiles.city", q}, SearchQuery{"profiles.description", q}, SearchQuery{"subjects.name", q}))
	}

	var totalTutors int64
	db.Model(&SubjectTaught{}).
		Joins("JOIN profiles ON profiles.id = subject_taughts.tutor_profile_id").
		Joins("JOIN subjects ON subject_taughts.subject_id = subjects.id").
		Scopes(
			scopes...,
		).
		Distinct("tutor_profile_id").Count(&totalTutors)

	scopes = append(scopes, Paginate(pageSize, page))
	asc := ""
	switch sort {
	case "low":
		asc = "asc"
	case "high":
		asc = "desc"
	}
	scopes = append(scopes, Sort("AVG( subject_taughts.price )", asc,
		Join{
			new:     Table{"subject_taughts", "tutor_profile_id"},
			current: Table{"profiles", "id"},
		},
		Join{
			new:     Table{"subjects", "id"},
			current: Table{"subject_taughts", "subject_id"},
		},
	))
	// Get tutors who are teaching subjects paginated
	var profiles []Profile
	err := db.
		Where("profiles.id IN (?)", db.Model(&SubjectTaught{}).Select("tutor_profile_id")).
		Preload("Subjects").Preload("Subjects.Subject").
		Scopes(
			scopes...,
		).
		Find(&profiles).
		Error
	if err != nil {
		return nil, 0, err
	}

	return profiles, int(math.Ceil(float64(totalTutors) / float64(pageSize))), nil
}

//Returns subjectTaught for specific tutors using their ID
func GetSubjectsTaughtByTutorID(tpid uuid.UUID, db *gorm.DB, preloads ...string) ([]SubjectTaught, error) {
	if tpid == uuid.MustParse("00000000-0000-0000-0000-000000000000") {
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
	return subjectTaught, db.Where(&SubjectTaught{TutorID: tpid}).Find(&subjectTaught).Error
}

//creats a StudentTaught based on the subject and tutor with a set price description.
func TeachSubject(subject *Subject, tutor *Account, description string, price int64, db *gorm.DB) error {
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
			Subject:        *subject,
			Tutor:          *tutor,
			TutorProfileID: tutor.Profile.ID,
			Description:    description,
			Price:          price,
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
func UpdateCost(stid uuid.UUID, price float32, db *gorm.DB) (*SubjectTaught, error) {
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

// Request a subject to be added
func RequestSubject(tutor *Account, name string) error {
	db, err := database.Open()
	if err != nil {
		return err
	}

	// Check if a subject matching this name is already in the database
	var subject Subject
	res := db.Where(&Subject{Name: name}).Find(&subject)
	if res.Error != nil {
		return err
	}
	if res.RowsAffected > 0 {
		return fmt.Errorf("A subject matching that name already exists")
	}

	return db.Create(&SubjectRequest{
		RequesterID: tutor.ID,
		Name:        name,
		Status:      SubjectRequestPending,
	}).Error
}
