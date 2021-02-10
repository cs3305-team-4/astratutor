package services

import (
	"github.com/cs3305-team-4/api/pkg/database"
	"github.com/google/uuid"
)

func CreateDebugData() error {
	db, err := database.Open()
	if err != nil {
		return err
	}

	hash, err := NewPasswordHash("grindsapp")
	if err != nil {
		return err
	}

	err = db.FirstOrCreate(&Account{
		Model: database.Model{
			ID: uuid.MustParse("22222222-2222-2222-2222-222222222222"),
		},
		Email:         "student@grindsapp.localhost",
		EmailVerified: true,
		Type:          Student,
		Suspended:     false,
		PasswordHash:  *hash,
		Profile: &Profile{
			Avatar:         "",
			Slug:           "john-student",
			FirstName:      "John",
			LastName:       "Student",
			City:           "Cork",
			Country:        "Ireland",
			Description:    "A student",
			Qualifications: []Qualification{},
			WorkExperience: []WorkExperience{},
			Availability:   nil,
			Color:          "#56847a",
		},
	}).Error
	if err != nil {
		return err
	}

	err = db.FirstOrCreate(&Account{
		Model: database.Model{
			ID: uuid.MustParse("11111111-1111-1111-1111-111111111111"),
		},
		Email:         "tutor@grindsapp.localhost",
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
			Availability:   nil,
			Color:          "#56847a",
		},
	}).Error
	if err != nil {
		return err
	}

	return nil
}
