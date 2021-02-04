package services

import (
	"database/sql/driver"
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
	AccountErrorProfileExists              AccountError = "A profile already exists for this account."
	AccountErrorAccountDoesNotExist        AccountError = "This account does not exist."
	AccountErrorProfileDoesNotExists       AccountError = "A profile does not exist for this account."
	AccountErrorQualificationDoesNotExists AccountError = "This qualification does not exist."
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

func CreateTestAccounts() error {
	db, err := database.Open()
	if err != nil {
		return err
	}

	hash, err := NewPasswordHash("grindshub")
	if err != nil {
		return err
	}

	err = db.FirstOrCreate(&Account{
		Model: database.Model{
			ID: uuid.MustParse("deadbeef-cafe-badd-c0de-facadebadbad"),
		},
		Email:         "student@grindshub.localhost",
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
			Availability:   &Availability{},
		},
	}).Error
	if err != nil {
		return err
	}

	err = db.FirstOrCreate(&Account{
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
			Slug:           "john-tutor",
			FirstName:      "John",
			LastName:       "Tutor",
			City:           "Cork",
			Country:        "Ireland",
			Description:    "A tutor",
			Qualifications: []Qualification{},
			WorkExperience: []WorkExperience{},
			Availability:   &Availability{},
		},
	}).Error
	if err != nil {
		return err
	}

	return nil
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

// AvailabilityLength length of slice.
const AvailabilityLength = 336

// Availability for tutors.
type Availability []bool

// Scan scan value into availability, implements sql.Scanner interface.
func (a *Availability) Scan(value interface{}) error {
	if value == nil {
		return nil
	}
	text, ok := value.(string)
	if !ok {
		return errors.New("Invalid value for availability.")
	}
	if len(text) < (AvailabilityLength*2)-1 {
		return errors.New("Invalid availability length.")
	}
	out := make(Availability, 0)
	text = text[1 : len(text)-1]
	for i := 0; i < len(text); i += 2 {
		fmt.Println(text[i])
		switch text[i] {
		case '0':
			out = append(out, false)
		case '1':
			out = append(out, true)
		}
	}
	*a = out
	return nil
}

// Value return availability value, implement driver.Valuer interface.
func (a *Availability) Value() (driver.Value, error) {
	if a != nil {
		out := []int{}
		for _, val := range *a {
			switch val {
			case false:
				out = append(out, 0)
			case true:
				out = append(out, 1)
			}
		}
		return out, nil
	}
	return nil, nil
}

func (a *Availability) Get() []bool {
	if a != nil {
		return *a
	}
	return make([]bool, AvailabilityLength)
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
	Availability *Availability `gorm:"type:int[]"`
}

// IsAccountType checks the account type of the given profile.
func (p *Profile) IsAccountType(accountType AccountType) (bool, error) {
	conn, err := database.Open()
	if err != nil {
		return false, err
	}
	account := &Account{}
	if err := conn.First(account, p.AccountID).Error; err != nil {
		return false, err
	}
	return account.Type == accountType, nil
}

// RemoveQualificationByID removes qualification inplace and in the DB.
func (p *Profile) RemoveQualificationByID(qualificationID uuid.UUID) error {
	conn, err := database.Open()
	if err != nil {
		return err
	}
	return conn.Transaction(func(tx *gorm.DB) error {
		for i, val := range p.Qualifications {
			if val.ID == qualificationID {
				p.Qualifications = append(p.Qualifications[:i], p.Qualifications[i+1:]...)
				if err = tx.Delete(&val).Error; err != nil {
					return err
				}
				return tx.Save(p).Error
			}
		}
		return AccountErrorQualificationDoesNotExists
	})
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
		return tx.Save(account).Error
	})
}

// UpdateProfileField will update a single profile field belonging to the provided account ID.
func UpdateProfileField(id uuid.UUID, key string, value interface{}) (*Profile, error) {
	conn, err := database.Open()
	if err != nil {
		return nil, err
	}
	var profile *Profile
	return profile, conn.Transaction(func(tx *gorm.DB) error {
		account, err := ReadAccountByID(id, tx, "Profile")
		if err != nil {
			return err
		}
		if profile = account.Profile; profile == nil {
			return AccountErrorProfileDoesNotExists
		}
		return tx.Model(profile).Update(key, value).Error
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
func ReadProfileByAccountID(id uuid.UUID, conn *gorm.DB, preloads ...string) (*Profile, error) {
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
	profile := &Profile{}
	return profile, conn.First(profile, "account_id = ?", id).Error
}

type Qualification struct {
	database.Model
	ProfileID uuid.UUID `gorm:"type:uuid"`
	Field     string
	Degree    string
	School    string
	Verified  bool
	// SupportingDocuments
}

// SetOnProfileByAccountID will set the qualification on the profile matching the
// provided account ID.
func (q *Qualification) SetOnProfileByAccountID(id uuid.UUID) (*Profile, error) {
	var err error
	conn, err := database.Open()
	if err != nil {
		return nil, err
	}
	var profile *Profile
	return profile, conn.Transaction(func(tx *gorm.DB) error {
		profile, err = ReadProfileByAccountID(id, tx, "Qualifications")
		if err != nil {
			return err
		}
		profile.Qualifications = append(profile.Qualifications, *q)
		return tx.Save(profile).Error
	})
}

type WorkExperience struct {
	database.Model
	ProfileID   uuid.UUID `gorm:"type:uuid"`
	Role        string
	YearsExp    string
	Description string
	Verified    bool
	// Supporting Documents
}
