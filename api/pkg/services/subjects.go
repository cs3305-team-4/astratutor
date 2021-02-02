package services

import "github.com/cs3305-team-4/api/pkg/database"

type Subject struct {
	database.Model
	Name string
	Slug string
}

type SubjectTaught struct {
	database.Model
	Subject     Subject
	Tutor       Account
	Price       uint
	Description string
}

func ReadSubjects() ([]Subject, error) {
	return nil, nil
}
