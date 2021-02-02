package services

import (
	"errors"
	"fmt"
	"strings"

	"github.com/cs3305-team-4/api/pkg/database"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

// AccountType is the type of account.
type AccountType string

const (
	Tutor   AccountType = "tutor"
	Student AccountType = "student"
)

// ToAccountType will cast to AccounType if it exists.
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
	Profile       *Profile     `gorm:"foreignKey:AccountID"`
}

// CreateAccount will create an account entry in the DB.
func CreateAccount(a *Account) error {
	conn, err := database.Open()
	if err != nil {
		return err
	}
	return conn.Create(a).Error
}

// GetAccountByID queries the DB by account ID.
func GetAccounteByID(id uuid.UUID, preloads ...string) (*Account, error) {
	conn, err := database.Open()
	if err != nil {
		return nil, err
	}
	for _, preload := range preloads {
		conn = conn.Preload(preload)
	}
	account := &Account{}
	return account, conn.First(account, id).Error
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

// Profile model.
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

// CreateProfile will create a profile entry in the DB relating to the Account from AccountID.
func CreateProfile(p *Profile) error {
	account, err := GetAccounteByID(p.AccountID, "Profile")
	if err != nil {
		return err
	}
	if account.Profile != nil {
		return errors.New("profile already exists")
	}

	// Generate slug
	name := fmt.Sprintf("%s-%s", strings.ToLower(p.FirstName), strings.ToLower(p.LastName))
	_, slugErr := GetProfileBySlug(name)
	i := 1
	slug := name
	for !errors.Is(slugErr, gorm.ErrRecordNotFound) {
		slug = fmt.Sprintf("%s-%d", name, i)
		_, slugErr = GetProfileBySlug(slug)
		i++
	}
	p.Slug = slug

	account.Profile = p
	conn, err := database.Open()
	if err != nil {
		return err
	}
	return conn.Save(account).Error
}

// GetProfileByAccountSlug queries the DB by slug.
func GetProfileBySlug(slug string) (*Profile, error) {
	conn, err := database.Open()
	if err != nil {
		return nil, err
	}
	profile := &Profile{}
	return profile, conn.First(profile, "slug = ?", slug).Error
}

// GetProfileByAccountID queries the DB by account ID.
func GetProfileByAccountID(id uuid.UUID) (*Profile, error) {
	return nil, nil
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
