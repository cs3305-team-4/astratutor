package services

import "github.com/cs3305-team-4/api/pkg/db"

type Subject struct {
	db.Model
	Name string
	Slug string
}

type SubjectTaught struct {
	db.Model
	Subject     Subject
	Tutor       Account
	Price       uint
	Description string
}

func ReadSubjects() ([]Subject, error) {
	return nil, nil
}
