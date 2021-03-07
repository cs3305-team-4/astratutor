package services

import (
	"bufio"
	"bytes"
	"encoding/base64"
	"fmt"
	"image/jpeg"
	"math/rand"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/cs3305-team-4/api/pkg/database"
	"github.com/google/uuid"
	"github.com/nfnt/resize"
	log "github.com/sirupsen/logrus"
	stripe "github.com/stripe/stripe-go/v72"
	stripeAccount "github.com/stripe/stripe-go/v72/account"
	stripeCustomer "github.com/stripe/stripe-go/v72/customer"
)

var seedFirstNames []string
var seedLastNames []string
var seedWords []string
var seedDegrees = []string{"Bachelors", "Masters", "Postdoctorate"}
var seedSchools = []string{}
var seedSubjects []string
var seedLocations []string

var subjects []Subject
var students []*Account
var tutors []*Account
var lessons []Lesson

func CreateDesc(minLen int) string {
	var words []string
	charLen := 0
	for charLen < minLen {
		word := seedWords[rand.Intn(len(seedWords))]
		words = append(words, word)
		charLen += len(word)
	}
	return strings.Join(words, " ")
}

func ReadFileToLines(path string) ([]string, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}

	defer f.Close()

	scanner := bufio.NewScanner(f)

	var lines []string

	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}

	return lines, nil
}

func RandomizeProfile(account *Account) {
	firstName := seedFirstNames[rand.Intn(len(seedFirstNames))]
	lastName := seedLastNames[rand.Intn(len(seedLastNames))]
	slug := strings.ToLower(firstName + "-" + lastName + "-" + strconv.Itoa(time.Now().UTC().Nanosecond()))
	location := strings.Split(seedLocations[rand.Intn(len(seedLocations))], ",")
	city, country := location[0], location[1]

	var subjectsTaught []SubjectTaught
	var qualifications []Qualification
	var workExperience []WorkExperience
	availability := make(Availability, 168)

	profileId, err := uuid.NewRandom()
	if err != nil {
		panic(err)
	}

	if account.Type == Tutor {
		for s := 1; s < rand.Intn(4)+2; s++ {
		createNew:
			taughtId, err := uuid.NewRandom()
			if err != nil {
				panic(err)
			}
			subject := SubjectTaught{
				Model: database.Model{
					ID: taughtId,
				},
				Subject:        subjects[rand.Intn(len(subjects))],
				Description:    CreateDesc(300),
				Price:          int64(rand.Intn(10000) + 10000),
				TutorID:        account.ID,
				TutorProfileID: profileId,
			}
			// Ensure tutor isnt already teaching this subject
			for _, subjectTaught := range subjectsTaught {
				if subjectTaught.Subject == subject.Subject {
					goto createNew
				}
			}
			subjectsTaught = append(subjectsTaught, subject)
		}

		// Generate random qualifications
		for q := 1; q < rand.Intn(6)+2; q++ {
			qualifications = append(qualifications, Qualification{
				Field:    seedSubjects[rand.Intn(len(seedSubjects))],
				Degree:   seedDegrees[rand.Intn(len(seedDegrees))],
				School:   seedSchools[rand.Intn(len(seedSchools))],
				Verified: true,
			})
		}

		// Generate random work experience
		if len(qualifications) > 0 {
			for w := 0; w < rand.Intn(3); w++ {
				workExperience = append(workExperience, WorkExperience{
					// Not much of a role but best I could do without having nicer
					// lists of subjects that include sub lists of roles
					Role:        qualifications[rand.Intn(len(qualifications))].Field,
					YearsExp:    2,
					Description: CreateDesc(30),
					Verified:    true,
				})
			}
		}

		// Generate random availability
		for i := 0; i < 168; i++ {
			if rand.Intn(10) == 0 {
				availability[i] = true
			}
		}
	}

	// Generate random hex color
	r := rand.Int63n(256)
	b := rand.Int63n(256)
	g := rand.Int63n(256)
	color := "#" +
		strconv.FormatInt(r, 16) +
		strconv.FormatInt(g, 16) +
		strconv.FormatInt(b, 16)

	// Get random avatar
	response, err := http.Get("https://thispersondoesnotexist.com/image")
	var avatarB64 string
	if err != nil {
		log.Warn("Failed to retreive avatar image from https://thispersondoesnotexist.com/image")
	} else {
		defer response.Body.Close()
		img, err := jpeg.Decode(response.Body)
		if err != nil {
			log.Warn("Failed to decode image from https://thispersondoesnotexist.com/image")
		} else {
			resizedImg := resize.Resize(400, 400, img, resize.Lanczos3)
			buf := new(bytes.Buffer)
			err = jpeg.Encode(buf, resizedImg, &jpeg.Options{Quality: 35})
			if err != nil {
				log.Error(err)
			}
			avatarB64 = "data:image/jpeg;base64," + base64.StdEncoding.EncodeToString([]byte(buf.Bytes()))
		}
	}

	account.Profile = &Profile{
		Model: database.Model{
			ID: profileId,
		},
		AccountID:      account.ID,
		Avatar:         avatarB64,
		Slug:           slug,
		FirstName:      firstName,
		LastName:       lastName,
		City:           city,
		Country:        country,
		Description:    CreateDesc(300),
		Color:          color,
		Subjects:       subjectsTaught,
		Qualifications: qualifications,
		WorkExperience: workExperience,
		Availability:   &availability,
	}
}

func CreateRandomAccountWithID(id uuid.UUID, index int, passwordHash *PasswordHash) *Account {
	// Generate random account type
	accountType := Student
	if rand.Intn(2) == 1 {
		accountType = Tutor
	}

	account := Account{
		Model: database.Model{
			ID: id,
		},
		Email:         strconv.Itoa(index) + "@grindsapp.localhost",
		EmailVerified: true,
		Type:          accountType,
		Suspended:     false,
		PasswordHash:  *passwordHash,
	}
	RandomizeProfile(&account)
	return &account
}

func SeedDatabase() error {
	db, err := database.Open()
	if err != nil {
		return err
	}

	log.Info("Seeding database. This may take serveral minutes...")

	//Read in random data for seeding
	seedFirstNames, err = ReadFileToLines("./seed/first_names")
	if err != nil {
		return err
	}

	seedLastNames, err = ReadFileToLines("./seed/last_names")
	if err != nil {
		return err
	}

	seedWords, err = ReadFileToLines("./seed/words")
	if err != nil {
		return err
	}

	seedSubjects, err = ReadFileToLines("./seed/subjects")
	if err != nil {
		return err
	}

	seedLocations, err = ReadFileToLines("./seed/locations")
	if err != nil {
		return err
	}

	seedSchools, err = ReadFileToLines("./seed/schools")
	if err != nil {
		return err
	}

	hash, err := NewPasswordHash("grindsapp")
	if err != nil {
		return err
	}

	accList := stripeAccount.List(nil)
	var emailToStripeConnectAccID map[string]string
	emailToStripeConnectAccID = map[string]string{}
	for accList.Next() {
		a := accList.Account()
		emailToStripeConnectAccID[a.Email] = a.ID
	}

	// Setting a fixed seed so that seeding is derterministic
	rand.Seed(2131287698123)

	// Create Subjects
	log.Info("Seeding Subjects...")
	for s := 0; s < len(seedSubjects); s++ {
		id := uuid.MustParse("11111111-1111-1111-1111-" + fmt.Sprintf("%012d", s))
		var subject Subject
		if db.Find(&subject, id).RowsAffected > 0 {
			continue
		}

		subject = Subject{
			Model: database.Model{
				ID: id,
			},
			Name: seedSubjects[s],
			Slug: strings.ToLower(seedSubjects[s]),
		}
		subjects = append(subjects, subject)
	}
	db.CreateInBatches(subjects, 20)

	fmt.Println(subjects)
	// Seed main test accounts
	log.Info("Seeding Tutor and Student account")
	var search Account
	tid := uuid.MustParse("11111111-1111-1111-1111-111111111111")
	sid := uuid.MustParse("22222222-2222-2222-2222-222222222222")
	if db.Find(&search, tid, sid).RowsAffected == 0 {
		tutor := Account{
			Model: database.Model{
				ID: tid,
			},
			Email:         "tutor@grindsapp.localhost",
			EmailVerified: true,
			Type:          Tutor,
			Suspended:     false,
			PasswordHash:  *hash,
		}

		RandomizeProfile(&tutor)
		if stripeConnectAccID, ok := emailToStripeConnectAccID[tutor.Email]; ok {
			tutor.StripeID = stripeConnectAccID
		} else {
			tutor.SetupBilling()
		}

		tutors = append(tutors, &tutor)

		student := Account{
			Model: database.Model{
				ID: sid,
			},
			Email:         "student@grindsapp.localhost",
			EmailVerified: true,
			Type:          Student,
			Suspended:     false,
			PasswordHash:  *hash,
		}
		custList := stripeCustomer.List(&stripe.CustomerListParams{
			Email: &student.Email,
		})

		for custList.Next() {
			c := custList.Customer()
			if c.Email == student.Email {
				student.StripeID = c.ID
			}
		}
		RandomizeProfile(&student)
		if student.StripeID == "" {
			student.SetupBilling()
		}
		students = append(students, &student)

		// Add lesson between tutor and student
		id := uuid.MustParse("11111111-1111-1111-1111-111111111111")
		startTime := time.Now().Add(time.Duration(rand.Intn(7000)) * time.Hour)
		endTime := startTime.Add(1 * time.Hour)

		lesson := Lesson{
			Model: database.Model{
				ID: id,
			},
			TutorID:               tutor.ID,
			StudentID:             student.ID,
			Tutor:                 tutor,
			Student:               student,
			StartTime:             startTime,
			EndTime:               endTime,
			LessonDetail:          CreateDesc(20),
			RequestStage:          Requested,
			RequestStageDetail:    CreateDesc(10),
			Requester:             student,
			RequesterID:           student.ID,
			RequestStageChangerID: student.ID,
			SubjectTaught:         tutor.Profile.Subjects[0],
			SubjectTaughtID:       tutor.Profile.Subjects[0].ID,
		}
		err = lesson.SetupPaymentIntent()
		if err != nil {
			panic(err)
		}
		lessons = append(lessons, lesson)
	}

	// Create random accounts
	log.Info("Seeding Accounts...")
	for i := 0; i < 20; i++ {
		id := uuid.MustParse("11111111-1111-1111-1111-" + fmt.Sprintf("%012d", i))
		// Check if already in database
		// This is to avoid generating a random account
		var user Account
		if db.Find(&user, id).RowsAffected > 0 {
			continue
		}

		account := CreateRandomAccountWithID(id, i, hash)
		if account.Type == Student {
			custList := stripeCustomer.List(&stripe.CustomerListParams{
				Email: &account.Email,
			})

			for custList.Next() {
				c := custList.Customer()
				if c.Email == account.Email {
					account.StripeID = c.ID
				}
			}

			if account.StripeID == "" {
				account.SetupBilling()
			}

			students = append(students, account)
		} else {
			if stripeConnectAccID, ok := emailToStripeConnectAccID[account.Email]; ok {
				account.StripeID = stripeConnectAccID
			} else {
				account.SetupBilling()
			}

			tutors = append(tutors, account)
		}
		log.Info("Seeding database with Account ID: ", id)

		// Pause to try avoid cached avatars being returned
		time.Sleep(10 * time.Millisecond)
	}
	db.CreateInBatches(students, 20)
	db.CreateInBatches(tutors, 20)

	// Create lessons between students and tutors
	log.Info("Seeding lessons...")
	// Disabled due to payment integration issues
	// for i, tutor := range tutors {
	// 	for l := 5; l < rand.Intn(6)+6; l++ {
	// 		id := uuid.MustParse("11111111-1111-1111-1111-" + fmt.Sprintf("%012d", l+((i+1)*10)))
	// 		student := students[rand.Intn(len(students))]

	// 		startTime := time.Now().Add(time.Duration(rand.Intn(7000)) * time.Hour)
	// 		endTime := startTime.Add(1 * time.Hour)

	// 		stage := []LessonRequestStage{Scheduled, Denied, Cancelled, Requested}[rand.Intn(4)]
	// 		changer := tutor.ID
	// 		if stage == Requested {
	// 			changer = student.ID
	// 		}

	// 		lesson := Lesson{
	// 			Model: database.Model{
	// 				ID: id,
	// 			},
	// 			TutorID:               tutor.ID,
	// 			StudentID:             student.ID,
	// 			StartTime:             startTime,
	// 			EndTime:               endTime,
	// 			LessonDetail:          CreateDesc(20),
	// 			RequestStage:          stage,
	// 			RequestStageDetail:    CreateDesc(10),
	// 			RequesterID:           student.ID,
	// 			RequestStageChangerID: changer,
	// 		}
	// 		lessons = append(lessons, lesson)
	// 	}
	// }

	db.CreateInBatches(lessons, 20)

	// Setting seed to be current time so that rand is no longer deterministic
	rand.Seed(time.Now().UnixNano())
	return nil
}
