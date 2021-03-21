package services

import (
	"time"

	"github.com/cs3305-team-4/api/pkg/database"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

//contains the information on a single review and what accounts it is connected to
type Review struct {
	database.Model
	Rating  int    `gorm:"not null;check:rating >= 0; check:rating <= 5;"`
	Comment string `gorm:""`

	TutorProfileID uuid.UUID `gorm:"type:uuid;index"`

	StudentProfileID uuid.UUID `gorm:"type:uuid;index"`
	Student          Profile   `gorm:"foreignKey:StudentProfileID;references:AccountID"`
}

type ReviewError string

func (e ReviewError) Error() string {
	return string(e)
}

const (
	ReviewErrorTutorNotFound     ReviewError = "Tutor account does not exist"
	ReviewErrorNoCompletedLesson ReviewError = "Student has not completed a lesson with this tutor"
	ReviewErrorStudentsOnly      ReviewError = "Only students can review tutors"
	ReviewErrorOnce              ReviewError = "You can onyl review a tutor once"
)

//Represents an incoming review to be added to the database
type ReviewCreateDTO struct {
	Rating  int    `json:"rating" validate:"required,gte=1,lte=5"`
	Comment string `json:"comment"`
}

//Represents a single review and the account connected to it
type ReviewDTO struct {
	ID                    uuid.UUID `json:"id" gorm:"type:uuid;column:id"`
	CreatedAt             time.Time `json:"created_at"`
	Rating                int       `json:"rating"`
	Comment               string    `json:"comment"`
	ProfileResponseDTOMin `json:"student" gorm:""`
}

//represents a profile connected to a review
type ProfileResponseDTOMin struct {
	AccountID string `json:"account_id"`
	Avatar    string `json:"avatar" validate:"omitempty"`
	Slug      string `json:"slug"`
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
}

//represents a  reviews update
type ReviewUpdateDTO struct {
	Rating  int    `json:"rating" validate:"gte=0,lte=5"`
	Comment string `json:"comment"`
}

//represnts a tutors average review score
type ReviewAverageDTO struct {
	Average float32 `json:"average"`
}

//adds a review to the database
func CreateReview(review *Review) error {
	conn, err := database.Open()
	if err != nil {
		return err
	}
	query := conn.Debug().Model(&Review{}).Where(review)
	if query.RowsAffected != 0 {
		return ReviewErrorOnce
	}

	return conn.Create(review).Error
}

//connects the pofile of a student who wrote the review to the review
func joinReviewProfile(db *gorm.DB) *gorm.DB {
	return db.Joins("LEFT JOIN profiles ON profiles.account_id = reviews.student_profile_id").
		Select([]string{"reviews.*", "profiles.account_id", "profiles.avatar", "profiles.slug", "profiles.first_name", "profiles.last_name"})
}

//returns all reviews for a given tutor
func TutorAllReviews(id uuid.UUID) ([]ReviewDTO, error) {
	conn, err := database.Open()
	if err != nil {
		return nil, err
	}

	var reviews []ReviewDTO
	err = conn.Table("reviews").
		Debug().
		Scopes(joinReviewProfile).
		Where("reviews.deleted_at IS NULL").
		Order("reviews.created_at desc").Where(&Review{
		TutorProfileID: id,
	}).Find(&reviews).Error
	return reviews, err
}

//returns a single review when given the associated tutor and review ids
func TutorSingleReview(tid uuid.UUID, rid uuid.UUID) (ReviewDTO, error) {
	conn, err := database.Open()
	if err != nil {
		return ReviewDTO{}, err
	}

	var review ReviewDTO
	query := Review{}
	query.ID = rid
	query.TutorProfileID = tid
	err = conn.Table("reviews").
		Where("reviews.deleted_at IS NULL").
		Scopes(joinReviewProfile).
		Where(&query).Limit(1).Find(&review).Error
	return review, err
}

//retuens a tutors average review score when given their ID
func TutorReviewsAverage(tid uuid.UUID) (ReviewAverageDTO, error) {
	conn, err := database.Open()
	if err != nil {
		return ReviewAverageDTO{}, err
	}

	var average ReviewAverageDTO
	err = conn.Model(&Review{}).Where(&Review{
		TutorProfileID: tid,
	}).Select("avg(rating) as average").Row().Scan(&average.Average)

	return average, err
}

//Returns a spefic review when given the associated tutor and review id
func TutorReviewByStudent(tid uuid.UUID, sid uuid.UUID) (ReviewDTO, error) {
	conn, err := database.Open()
	if err != nil {
		return ReviewDTO{}, err
	}

	var review ReviewDTO
	query := Review{}
	query.TutorProfileID = tid
	query.StudentProfileID = sid
	err = conn.Table("reviews").
		Where("reviews.deleted_at IS NULL").
		Scopes(joinReviewProfile).
		Where(&query).Limit(1).Find(&review).Error
	return review, err
}

//updates a reviews rating by a given review id
func UpdateReviewRating(rid uuid.UUID, rating int, sid uuid.UUID) error {
	conn, err := database.Open()
	if err != nil {
		return err
	}

	review := Review{}
	review.ID = rid
	review.StudentProfileID = sid
	return conn.Model(&Review{}).Where(&review).Update("rating", rating).Error
}

//updates a reviews comment by a given review id
func UpdateReviewComment(rid uuid.UUID, comment string, sid uuid.UUID) error {
	conn, err := database.Open()
	if err != nil {
		return err
	}

	review := Review{}
	review.ID = rid
	review.StudentProfileID = sid
	return conn.Model(&Review{}).Where(&review).Update("comment", comment).Error
}

//deletes a speific review by the given tutor, review and student id
func TutorDeleteReview(tid uuid.UUID, rid uuid.UUID, sid uuid.UUID) error {
	conn, err := database.Open()
	if err != nil {
		return err
	}

	review := &Review{}
	review.TutorProfileID = tid
	review.ID = rid
	review.StudentProfileID = sid
	return conn.Model(&Review{}).Delete(review).Error
}
