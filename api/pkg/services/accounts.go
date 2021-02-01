package services

// Type is the type of account
type Type string

const (
	Tutor   = "tutor"
	Student = "student"
)

type Account struct {
	ID            uint64
	Slug          string
	Email         string
	EmailVerified bool
	Type          Type
	Suspended     bool
}

type PasswordHash struct {
	ID   string
	Hash []byte
	Salt string
}

type TutorProfile struct {
	FirstName     string
	LastName      string
	City          string
	Country       string
	Description   string
}

type StudentProfile struct {
	FirstName	  string
	LastName      string
	City		  string
	Country		  string
}

