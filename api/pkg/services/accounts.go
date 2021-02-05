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

// AccountError types.
type AccountError string

func (e AccountError) Error() string {
	return string(e)
}

const (
	AccountErrorProfileExists AccountError = "a profile already exists for this account"
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

func (a *Account) IsStudent() bool {
	return a.Type == Student
}

func (a *Account) IsTutor() bool {
	return a.Type == Tutor
}

// CreateAccount will create an account entry in the DB.
func CreateAccount(a *Account) error {
	conn, err := database.Open()
	if err != nil {
		return err
	}

	return conn.Create(a).Error
}

// ReadAccountByID queries the DB by account ID.
// conn is optional.
func ReadAccountByID(id uuid.UUID, conn *gorm.DB, preloads ...string) (*Account, error) {
	if conn == nil {
		var err error
		conn, err = database.Open()
		if err != nil {
			return nil, err
		}
	}
	for _, preload := range preloads {
		conn = conn.Preload(preload)
	}
	account := &Account{}
	return account, conn.First(account, id).Error
}

func ReadAccountByEmail(email string, conn *gorm.DB, preloads ...string) (*Account, error) {
	if conn == nil {
		var err error
		conn, err = database.Open()
		if err != nil {
			return nil, err
		}
	}
	for _, preload := range preloads {
		conn = conn.Preload(preload)
	}

	account := &Account{}
	return account, conn.Where(
		&Account{Email: email},
	).First(&account).Error
}

func ReadTutorByID(id uuid.UUID, conn *gorm.DB, preloads ...string) (*Account, error) {
	account, err := ReadAccountByID(id, conn, preloads...)
	if err != nil {
		return nil, err
	}

	if !account.IsTutor() {
		return nil, errors.New("the specified account is not a tutor")
	}

	return account, nil
}

func ReadStudentByID(id uuid.UUID, conn *gorm.DB, preloads ...string) (*Account, error) {
	account, err := ReadAccountByID(id, conn, preloads...)
	if err != nil {
		return nil, err
	}

	if !account.IsStudent() {
		return nil, errors.New("the specified account is not a student")
	}

	return account, nil
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

func (p *PasswordHash) ValidMatch(match string) bool {
	res := bcrypt.CompareHashAndPassword(p.Hash, []byte(match))

	if res == nil {
		return true
	}

	return false
}

func ReadPasswordHashByAccountID(id uuid.UUID) (*PasswordHash, error) {
	db, err := database.Open()
	if err != nil {
		return nil, err
	}

	var hash PasswordHash
	err = db.Where(&PasswordHash{
		AccountID: id,
	}).Find(&hash).Error

	if err != nil {
		return nil, err
	}

	return &hash, nil
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
	conn, err := database.Open()
	if err != nil {
		return err
	}
	return conn.Transaction(func(tx *gorm.DB) error {
		account, err := ReadAccountByID(p.AccountID, tx, "Profile")
		if err != nil {
			return err
		}
		if account.Profile != nil {
			return AccountErrorProfileExists
		}

		// Generate slug
		name := fmt.Sprintf("%s-%s", strings.ToLower(p.FirstName), strings.ToLower(p.LastName))
		_, slugErr := ReadProfileBySlug(name, tx)
		i := 1
		slug := name
		for !errors.Is(slugErr, gorm.ErrRecordNotFound) {
			slug = fmt.Sprintf("%s-%d", name, i)
			_, slugErr = ReadProfileBySlug(slug, tx)
			i++
		}
		p.Slug = slug

		account.Profile = p
		return conn.Save(account).Error
	})
}

// ReadProfileBySlug queries the DB by slug.
// conn is optional.
func ReadProfileBySlug(slug string, conn *gorm.DB) (*Profile, error) {
	if conn == nil {
		var err error
		conn, err = database.Open()
		if err != nil {
			return nil, err
		}
	}
	profile := &Profile{}
	return profile, conn.First(profile, "slug = ?", slug).Error
}

// ReadProfileByAccountID queries the DB by account ID.
// conn is optional.
func ReadProfileByAccountID(id uuid.UUID, conn *gorm.DB) (*Profile, error) {
	if conn == nil {
		var err error
		conn, err = database.Open()
		if err != nil {
			return nil, err
		}
	}
	profile := &Profile{}
	return profile, conn.First(profile, "account_id = ?", id).Error
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
