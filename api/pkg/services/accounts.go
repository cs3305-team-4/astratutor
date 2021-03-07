package services

import (
	"database/sql/driver"
	"errors"
	"fmt"
	"math/rand"
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
	AccountErrorProfileExists        AccountError = "A profile already exists for this account."
	AccountErrorAccountDoesNotExist  AccountError = "This account does not exist."
	AccountErrorProfileDoesNotExists AccountError = "A profile does not exist for this account."
	AccountErrorEntryDoesNotExists   AccountError = "This Entry does not exist."
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

	// StripeID corresponds to a customer ID if the account type is a Student or a Stripe Connect account ID if the account type is a Tutor
	StripeID     string
	PasswordHash PasswordHash `gorm:"foreignKey:AccountID"`
	Profile      *Profile     `gorm:"foreignKey:AccountID"`
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

// UpdateAccountField will update a single account field belonging to the provided account ID.
func UpdateAccountField(id uuid.UUID, key string, value interface{}) (*Account, error) {
	conn, err := database.Open()
	if err != nil {
		return nil, err
	}
	var account *Account
	return account, conn.Transaction(func(tx *gorm.DB) error {
		account, err = ReadAccountByID(id, tx)
		if err != nil {
			return err
		}
		return tx.Model(account).Update(key, value).Error
	})
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

// SetOnAccountID will delete the previous password hash and set it to the new one.
func (p PasswordHash) SetOnAccountByID(id uuid.UUID) (*Account, error) {
	conn, err := database.Open()
	if err != nil {
		return nil, err
	}
	var account *Account
	return account, conn.Transaction(func(tx *gorm.DB) error {
		account, err = ReadAccountByID(id, tx, "PasswordHash")
		if err != nil {
			return err
		}
		if err = tx.Delete(&account.PasswordHash).Error; err != nil {
			return err
		}
		account.PasswordHash = p
		return tx.Save(account).Error
	})
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

// NewPasswordHash will generate a password hash object. Storage should be done via CreateAccount.
func NewPasswordHash(password string) (*PasswordHash, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost) // Salt embedded in hash
	if err != nil {
		return nil, err
	}
	return &PasswordHash{Hash: hash}, nil
}

// AvailabilityLength length of slice.
const AvailabilityLength = 168

// Availability for tutors.
type Availability []bool

// Scan scan value into availability, implements sql.Scanner interface.
func (a *Availability) Scan(value interface{}) error {
	if value == nil {
		return nil
	}
	text, ok := value.(string)
	if !ok {
		return errors.New("invalid value for availability.")
	}
	if len(text) < (AvailabilityLength*2)-1 {
		return errors.New("invalid availability length.")
	}
	out := make(Availability, 0)
	text = text[1 : len(text)-1]
	for i := 0; i < len(text); i += 2 {
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
	AccountID      uuid.UUID `gorm:"type:uuid;unique"`
	Avatar         string
	Slug           string
	FirstName      string
	LastName       string
	City           string
	Country        string
	Subtitle       string
	Description    string
	Color          string
	Subjects       []SubjectTaught  `gorm:"foreignKey:TutorProfileID"`
	Qualifications []Qualification  `gorm:"foreignKey:ProfileID"`
	WorkExperience []WorkExperience `gorm:"foreignKey:ProfileID"`

	// Contains the next 14x24 hrs of availbility modulus to 1 week
	Availability *Availability `gorm:"type:int[]"`
}

// FilterVerifiedFields will filter qualifications and work experience for verified in-place.
func (p *Profile) FilterVerifiedFields() {
	qualifications := make([]Qualification, 0)
	workExperience := make([]WorkExperience, 0)
	for _, val := range p.Qualifications {
		if val.Verified {
			qualifications = append(qualifications, val)
		}
	}
	for _, val := range p.WorkExperience {
		if val.Verified {
			workExperience = append(workExperience, val)
		}
	}
	p.Qualifications = qualifications
	p.WorkExperience = workExperience
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
		return AccountErrorEntryDoesNotExists
	})
}

// RemoveWorkExperienceByID removes work experience inplace and in the DB.
func (p *Profile) RemoveWorkExperienceByID(expID uuid.UUID) error {
	conn, err := database.Open()
	if err != nil {
		return err
	}
	return conn.Transaction(func(tx *gorm.DB) error {
		for i, val := range p.WorkExperience {
			if val.ID == expID {
				p.WorkExperience = append(p.WorkExperience[:i], p.WorkExperience[i+1:]...)
				if err = tx.Delete(&val).Error; err != nil {
					return err
				}
				return tx.Save(p).Error
			}
		}
		return AccountErrorEntryDoesNotExists
	})
}

// Save the current profile in the DB.
func (p *Profile) Save() error {
	conn, err := database.Open()
	if err != nil {
		return err
	}
	return conn.Save(p).Error
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
		if err = p.GenerateNewSlug(tx); err != nil {
			return err
		}
		p.GenerateNewColor()

		// Setup billing for the account (as the profile has the required fields we can pre-input as a customer)
		account.Profile = p
		err = account.SetupBilling()
		if err != nil {
			tx.Rollback()
			return err
		}

		return tx.Save(account).Error
	})
}

// GenerateNewColor for sample avatar.
func (p *Profile) GenerateNewColor() {
	r := rand.Intn(255)
	g := rand.Intn(255)
	b := rand.Intn(255)
	p.Color = fmt.Sprintf("#%x%x%x", r, g, b)
}

// GenerateNewSlug for account
func (p *Profile) GenerateNewSlug(conn *gorm.DB) error {
	if conn == nil {
		var err error
		conn, err = database.Open()
		if err != nil {
			return err
		}
	}
	return conn.Transaction(func(tx *gorm.DB) error {
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
		return nil
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
func ReadProfileByAccountID(id uuid.UUID, conn *gorm.DB) (*Profile, error) {
	return readProfileByAccountID(id, conn, "Qualifications", "WorkExperience")
}

func readProfileByAccountID(id uuid.UUID, conn *gorm.DB, preloads ...string) (*Profile, error) {
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

// Settable on Profile.
type Settable interface {
	SetOnProfileByAccountID(id uuid.UUID) (*Profile, error)
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
		profile, err = ReadProfileByAccountID(id, tx)
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
	YearsExp    int
	Description string
	Verified    bool
	// Supporting Documents
}

// SetOnProfileByAccountID will set the work experience on the profile matching the
// provided account ID.
func (w *WorkExperience) SetOnProfileByAccountID(id uuid.UUID) (*Profile, error) {
	var err error
	conn, err := database.Open()
	if err != nil {
		return nil, err
	}
	var profile *Profile
	return profile, conn.Transaction(func(tx *gorm.DB) error {
		profile, err = ReadProfileByAccountID(id, tx)
		if err != nil {
			return err
		}
		profile.WorkExperience = append(profile.WorkExperience, *w)
		return tx.Save(profile).Error
	})
}
