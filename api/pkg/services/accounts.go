package services

import (
	"fmt"

	"github.com/cs3305-team-4/api/pkg/database"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

// AccountType is the type of account.
type AccountType string

const (
	Tutor   AccountType = "tutor"
	Student AccountType = "student"
)

// ToAccountType will case to AccounType if it exists.
func ToAccountType(s string) (AccountType, error) {
	switch AccountType(s) {
	case Tutor:
		return Tutor, nil
	case Student:
		return Student, nil
	default:
		return "", fmt.Errorf("Couldn't find account type %s", s)
	}
}

// Account model.
type Account struct {
	database.Model
	Email         string `gorm:"unique;not null;"`
	EmailVerified bool
	Type          AccountType
	Suspended     bool
	PasswordHash  PasswordHash `gorm:"foreignKey:AccountID"`
	Profile       Profile      `gorm:"foreignKey:AccountID"`
}

// CreateAccount will create an account entry in the DB.
func CreateAccount(a *Account) error {
	return conn.Create(a).Error
}

type PasswordHash struct {
	database.Model
	AccountID uuid.UUID `gorm:"type:uuid"`
	Hash      []byte    `gorm:"type:text"`
}

// NewPasswordHash will generate a password hash object. Storage should be done via CreateAccount.
func NewPasswordHash(password string) (*PasswordHash, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost) // Salt embedded in hash
	if err != nil {
		return nil, err
	}
	return &PasswordHash{Hash: hash}, nil
}

type Profile struct {
	database.Model
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
	Availability []bool `gorm:"type:text"`
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
