package services

import (
	"time"

	"github.com/cs3305-team-4/api/pkg/database"
	"github.com/google/uuid"
	log "github.com/sirupsen/logrus"
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

	// Add Student Account
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
			Color:          "#46247a",
		},
	}).Error
	if err != nil {
		return err
	}

	log.Info("Added student account: student@grindsapp.localhost grindsapp")

	// Add Tutor Acccounts
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

	log.Info("Added tutor account: tutor@grindsapp.localhost grindsapp")

	err = db.FirstOrCreate(&Account{
		Model: database.Model{
			ID: uuid.MustParse("33333333-3333-3333-3333-333333333333"),
		},
		Email:         "jane@grindsapp.localhost",
		EmailVerified: true,
		Type:          Tutor,
		Suspended:     false,
		PasswordHash:  *hash,
		Profile: &Profile{
			Avatar:         "",
			Slug:           "jane-smith",
			FirstName:      "Jane",
			LastName:       "Smith",
			City:           "Cork",
			Country:        "Ireland",
			Description:    "A tutor",
			Qualifications: []Qualification{},
			WorkExperience: []WorkExperience{},
			Availability:   nil,
			Color:          "#27c97a",
		},
	}).Error
	if err != nil {
		return err
	}

	log.Info("Added tutor account: jane@grindsapp.localhost grindsapp")

	// Add Subjects
	err = db.FirstOrCreate(&Subject{
		Model: database.Model{ID: uuid.MustParse("11111111-1111-1111-1111-111111111111")},
		Name:  "Leaving Certificate - English", Slug: "lc-english"}).Error
	if err != nil {
		return err
	}
	log.Info("Added subject Leaving Certificate - English")
	err = db.FirstOrCreate(&Subject{
		Model: database.Model{ID: uuid.MustParse("22222222-2222-2222-2222-222222222222")},
		Name:  "Leaving Certificate - Irish", Slug: "lc-irish"}).Error
	if err != nil {
		return err
	}
	log.Info("Added subject Leaving Certificate - Irish")
	err = db.FirstOrCreate(&Subject{
		Model: database.Model{ID: uuid.MustParse("33333333-3333-3333-3333-333333333333")},
		Name:  "Leaving Certificate - Maths", Slug: "lc-maths"}).Error
	if err != nil {
		return err
	}
	log.Info("Added subject Leaving Certificate - Maths")
	err = db.FirstOrCreate(&Subject{
		Model: database.Model{ID: uuid.MustParse("44444444-4444-4444-4444-444444444444")},
		Name:  "Leaving Certificate - Physics", Slug: "lc-physics"}).Error
	if err != nil {
		return err
	}
	log.Info("Added subject Leaving Certificate - Physics")
	err = db.FirstOrCreate(&Subject{
		Model: database.Model{ID: uuid.MustParse("55555555-5555-5555-5555-555555555555")},
		Name:  "Leaving Certificate - Chemistry", Slug: "lc-chemistry"}).Error
	if err != nil {
		return err
	}
	log.Info("Added subject Leaving Certificate - Chemistry")
	err = db.FirstOrCreate(&Subject{
		Model: database.Model{ID: uuid.MustParse("66666666-6666-6666-6666-666666666666")},
		Name:  "Leaving Certificate - Biology", Slug: "lc-biology"}).Error
	if err != nil {
		return err
	}
	log.Info("Added subject Leaving Certificate - Biology")
	err = db.FirstOrCreate(&Subject{
		Model: database.Model{ID: uuid.MustParse("77777777-7777-7777-7777-777777777777")},
		Name:  "Leaving Certificate - Engineering", Slug: "lc-engineering"}).Error
	if err != nil {
		return err
	}
	log.Info("Added subject Leaving Certificate - Engineering")
	err = db.FirstOrCreate(&Subject{
		Model: database.Model{ID: uuid.MustParse("88888888-8888-8888-8888-888888888888")},
		Name:  "Leaving Certificate - Construction Studies", Slug: "lc-construction-studies"}).Error
	if err != nil {
		return err
	}
	log.Info("Added subject Leaving Certificate - Construction Studies")
	err = db.FirstOrCreate(&Subject{
		Model: database.Model{ID: uuid.MustParse("99999999-9999-9999-9999-999999999999")},
		Name:  "Leaving Certificate - Technical Graphics", Slug: "lc-technical-graphics"}).Error
	if err != nil {
		return err
	}
	log.Info("Added subject Leaving Certificate - Technical Graphics")
	err = db.FirstOrCreate(&Subject{
		Model: database.Model{ID: uuid.MustParse("21233124-2222-2222-2222-234235656756")},
		Name:  "Leaving Certificate - Religion", Slug: "lc-religion"}).Error
	if err != nil {
		return err
	}
	log.Info("Added subject Leaving Certificate - Religion")

	// Add Tutors to subjects
	john, err := ReadAccountByID(uuid.MustParse("11111111-1111-1111-1111-111111111111"), nil, "Profile")
	if err != nil {
		return err
	}
	jane, err := ReadAccountByID(uuid.MustParse("33333333-3333-3333-3333-333333333333"), nil, "Profile")
	if err != nil {
		return err
	}
	english, err := GetSubjectBySlug("lc-english", nil)
	if err != nil {
		return err
	}
	irish, err := GetSubjectBySlug("lc-irish", nil)
	if err != nil {
		return err
	}
	physics, err := GetSubjectBySlug("lc-physics", nil)
	if err != nil {
		return err
	}

	log.Info("x")
	err = db.FirstOrCreate(&SubjectTaught{
		Model:        database.Model{ID: uuid.MustParse("11111111-1111-1111-1111-111111111111")},
		Subject:      *english,
		TutorProfile: *john.Profile,
		Description:  "Lorem ipsum dolor sit amet, consectetur adipiscing elit. Morbi interdum dignissim ipsum, sit amet scelerisque quam auctor id. Suspendisse laoreet commodo libero vitae volutpat. Integer hendrerit congue posuere. Pellentesque vestibulum leo at nunc interdum, gravida consequat dui egestas. Donec vel lobortis lorem. Donec suscipit, arcu vel dignissim ultricies, mi nibh tincidunt velit, eu dapibus justo metus id metus. Pellentesque porttitor nec augue eu molestie. Morbi eget lacinia arcu. Aliquam ornare risus mi, aliquam eleifend dolor consequat at.",
		Price:        45,
	}).Error
	if err != nil {
		return err
	}
	log.Info("John Tutor now teaching Leaving Certificate English")
	err = db.FirstOrCreate(&SubjectTaught{
		Model:        database.Model{ID: uuid.MustParse("22222222-2222-2222-2222-222222222222")},
		Subject:      *irish,
		TutorProfile: *john.Profile,
		Description:  "Lorem ipsum dolor sit amet, consectetur adipiscing elit. Morbi interdum dignissim ipsum, sit amet scelerisque quam auctor id. Suspendisse laoreet commodo libero vitae volutpat. Integer hendrerit congue posuere. Pellentesque vestibulum leo at nunc interdum, gravida consequat dui egestas. Donec vel lobortis lorem. Donec suscipit, arcu vel dignissim ultricies, mi nibh tincidunt velit, eu dapibus justo metus id metus. Pellentesque porttitor nec augue eu molestie. Morbi eget lacinia arcu. Aliquam ornare risus mi, aliquam eleifend dolor consequat at.",
		Price:        50,
	}).Error
	if err != nil {
		return err
	}
	log.Info("John Tutor now teaching Leaving Certificate Irish")
	err = db.FirstOrCreate(&SubjectTaught{
		Model:        database.Model{ID: uuid.MustParse("33333333-3333-3333-3333-333333333333")},
		Subject:      *physics,
		TutorProfile: *jane.Profile,
		Description:  "Lorem ipsum dolor sit amet, consectetur adipiscing elit. Morbi interdum dignissim ipsum, sit amet scelerisque quam auctor id. Suspendisse laoreet commodo libero vitae volutpat. Integer hendrerit congue posuere. Pellentesque vestibulum leo at nunc interdum, gravida consequat dui egestas. Donec vel lobortis lorem. Donec suscipit, arcu vel dignissim ultricies, mi nibh tincidunt velit, eu dapibus justo metus id metus. Pellentesque porttitor nec augue eu molestie. Morbi eget lacinia arcu. Aliquam ornare risus mi, aliquam eleifend dolor consequat at.",
		Price:        35,
	}).Error
	if err != nil {
		return err
	}
	log.Info("Jane Smith now teaching Leaving Certificate Physics")

	// Add lesson between student and tutor
	student, err := ReadAccountByID(uuid.MustParse("22222222-2222-2222-2222-222222222222"), nil)
	if err != nil {
		return err
	}

	startTime, err := time.Parse("02 Jan 2006 15:04:05", "24 Feb 2021 17:00:00")
	if err != nil {
		return err
	}
	endTime, err := time.Parse("02 Jan 2006 15:04:05", "24 Feb 2021 17:59:00")
	if err != nil {
		return err
	}
	err = db.FirstOrCreate(&Lesson{
		//Model:               database.Model{ID: uuid.MustParse("11111111-1111-1111-1111-111111111111")},
		StartTime:           startTime,
		EndTime:             endTime,
		Student:             *student,
		Tutor:               *john,
		LessonDetail:        "Quick question",
		RequestStage:        Accepted,
		RequestStageDetail:  "",
		Requester:           *student,
		RequestStageChanger: *john,
	}).Error
	if err != nil {
		return err
	}
	log.Info("John Tutor has lesson with John Student")

	return nil
}
