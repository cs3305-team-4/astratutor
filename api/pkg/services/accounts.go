package services

import (
	"github.com/cs3305-team-4/api/pkg/db"
	"github.com/google/uuid"
)

// Type is the type of account
type AccountType string

const (
	Tutor   AccountType = "tutor"
	Student             = "student"
)

type Account struct {
	db.Model
	Email         string
	EmailVerified bool
	Type          AccountType
	Suspended     bool
	PasswordHash  PasswordHash `gorm:"foreignKey:AccountID"`
	Profile       Profile      `gorm:"foreignKey:AccountID"`
}

type PasswordHash struct {
	db.Model
	AccountID uuid.UUID `gorm:"type:uuid"`
	Hash      []byte
	Salt      string
}

type Profile struct {
	db.Model
	AccountID      uuid.UUID `gorm:"type:uuid"`
	Avatar         string
	Slug           string
	FirstName      string
	LastName       string
	City           string
	Country        string
	Description    string
	Qualifications []Qualification  `gorm:"foreignKey:ProfileID"`
	WorkExperience []WorkExperience `gorm:"foreignKey:ProfileID"`

	// Contains the next 14x24 hrs of availbility modulus to 2 weeks
	Availability []bool
}

type Qualification struct {
	ProfileID uuid.UUID `gorm:"type:uuid"`
	Field     string
	Degree    string
	School    string
	Verified  bool
	// SupportingDocuments
}

type WorkExperience struct {
	ProfileID   uuid.UUID `gorm:"type:uuid"`
	Role        string
	YearsExp    string
	Description string
	Verified    bool
	// Supporting Documents
}
