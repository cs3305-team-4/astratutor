package services

import (
	"time"

	"github.com/cs3305-team-4/api/pkg/database"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

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

type ReviewCreateDTO struct {
	Rating  int    `json:"rating" validate:"required,gte=1,lte=5"`
	Comment string `json:"comment"`
}

type ReviewDTO struct {
	ID                    uuid.UUID `json:"id" gorm:"type:uuid;column:id"`
	CreatedAt             time.Time `json:"created_at"`
	Rating                int       `json:"rating"`
	Comment               string    `json:"comment"`
	ProfileResponseDTOMin `json:"student" gorm:""`
}

type ProfileResponseDTOMin struct {
	AccountID string `json:"account_id"`
	Avatar    string `json:"avatar" validate:"omitempty"`
	Slug      string `json:"slug"`
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
}

type ReviewUpdateDTO struct {
	Rating  int    `json:"rating" validate:"gte=0,lte=5"`
	Comment string `json:"comment"`
}

type ReviewAverageDTO struct {
	Average float32 `json:"average"`
}

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

func joinReviewProfile(db *gorm.DB) *gorm.DB {
	return db.Joins("LEFT JOIN profiles ON profiles.account_id = reviews.student_profile_id").
		Select([]string{"reviews.*", "profiles.account_id", "profiles.avatar", "profiles.slug", "profiles.first_name", "profiles.last_name"})
}

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
