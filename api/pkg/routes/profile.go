package routes

import (
	"errors"
	"net/http"
	"strings"

	"github.com/cs3305-team-4/api/pkg/services"
)

func getAccountType(r *http.Request) (services.AccountType, error) {
	base := strings.Split(r.URL.Path, "/")[1]
	if len(base) > 1 {
		return services.AccountType(strings.TrimSuffix(base, "s")), nil
	}
	return "", errors.New("invalid account type")
}

// ProfileRequestDTO create DTO.
type ProfileRequestDTO struct {
	Avatar      string `json:"avatar" validate:"omitempty"`
	FirstName   string `json:"first_name" validate:"required"`
	LastName    string `json:"last_name" validate:"required"`
	City        string `json:"city" validate:"required"`
	Country     string `json:"country" validate:"required"`
	Description string `json:"description" validate:"omitempty,lte=1000"`
}

// ProfileResponseDTO return DTO.
type ProfileResponseDTO struct {
	AccountID   string `json:"account_id" validate:"required,uuid"`
	ID          string `json:"id" validate:"required,uuid"`
	Avatar      string `json:"avatar" validate:"omitempty"`
	Slug        string `json:"slug" validate:"required"`
	FirstName   string `json:"first_name" validate:"required"`
	LastName    string `json:"last_name" validate:"required"`
	City        string `json:"city" validate:"required"`
	Country     string `json:"country" validate:"required"`
	Description string `json:"description" validate:"omitempty,lte=1000"`
}

// TutorResponseDTO return DTO.
type TutorResponseDTO struct {
	ProfileResponseDTO
	Qualifications []QualificationsResponseDTO `json:"qualifications"`
	WorkExperience []WorkExperienceResponseDTO `json:"work_experience"`
	Availability   []bool                      `json:"availability" validate:""`
}

// QualificationsRequestDTO create DTO.
type QualificationsRequestDTO struct {
	Field  string `json:"field" validate:"required"`
	Degree string `json:"degree" validate:"required"`
	School string `json:"school" validate:"required"`
}

// QualificationsResponseDTO return DTO.
type QualificationsResponseDTO struct {
	ID       string `json:"id" validate:"required,uuid"`
	Field    string `json:"field" validate:"required"`
	Degree   string `json:"degree" validate:"required"`
	School   string `json:"school" validate:"required"`
	Verified bool   `json:"verified" validate:"required"`
}

// WorkExperienceRequestDTO create DTO.
type WorkExperienceRequestDTO struct {
	Role        string `json:"role" validate:"required"`
	YearsExp    int    `json:"years_exp" validate:"required"`
	Description string `json:"description" validate:"required,lte=1000"`
}

// WorkExperienceResponseDTO return DTO.
type WorkExperienceResponseDTO struct {
	ID          string `json:"id" validate:"required,uuid"`
	Role        string `json:"role" validate:"required"`
	YearsExp    int    `json:"years_exp" validate:"required"`
	Description string `json:"description" validate:"required,lte=1000"`
	Verified    bool   `json:"verified" validate:"required"`
}

func dtoFromProfile(p *services.Profile, accountType services.AccountType) interface{} {
	switch accountType {
	case services.Student:
		return &ProfileResponseDTO{
			AccountID:   p.AccountID.String(),
			ID:          p.ID.String(),
			Avatar:      p.Avatar,
			Slug:        p.Slug,
			FirstName:   p.FirstName,
			LastName:    p.LastName,
			City:        p.City,
			Country:     p.Country,
			Description: p.Description,
		}
	case services.Tutor:
		qualifications := make([]QualificationsResponseDTO, 0)
		for _, val := range p.Qualifications {
			qualifications = append(qualifications, QualificationsResponseDTO{
				ID:       val.ID.String(),
				Field:    val.Field,
				Degree:   val.Degree,
				School:   val.School,
				Verified: val.Verified,
			})
		}
		workExperience := make([]WorkExperienceResponseDTO, 0)
		for _, val := range p.WorkExperience {
			workExperience = append(workExperience, WorkExperienceResponseDTO{
				ID:          val.ID.String(),
				Role:        val.Role,
				YearsExp:    val.YearsExp,
				Description: val.Description,
				Verified:    val.Verified,
			})
		}
		return &TutorResponseDTO{
			ProfileResponseDTO: ProfileResponseDTO{
				AccountID:   p.AccountID.String(),
				ID:          p.ID.String(),
				Avatar:      p.Avatar,
				Slug:        p.Slug,
				FirstName:   p.FirstName,
				LastName:    p.LastName,
				City:        p.City,
				Country:     p.Country,
				Description: p.Description,
			},
			Availability:   p.Availability.Get(),
			Qualifications: qualifications,
			WorkExperience: workExperience,
		}
	}
	return nil
}

func handleProfileGet(w http.ResponseWriter, r *http.Request) {
	t, err := getAccountType(r)
	if err != nil {
		restError(w, r, err, http.StatusBadRequest)
		return
	}
	id, err := getUUID(r, "uuid")
	if err != nil {
		restError(w, r, err, http.StatusBadRequest)
		return
	}
	serviceProfile, err := services.ReadProfileByAccountID(id, nil)
	if err != nil {
		restError(w, r, err, http.StatusNotFound)
		return
	}
	if ok, err := serviceProfile.IsAccountType(t); err != nil {
		restError(w, r, err, http.StatusInternalServerError)
		return
	} else if !ok {
		restError(w, r, errors.New("Account type does not match endpoint."), http.StatusBadRequest)
		return
	}
	// Filter verified if same user isn't authenticated.
	auth, err := ParseRequestAuth(r)
	if err != nil {
		restError(w, r, err, http.StatusInternalServerError)
		return
	}
	if !auth.Authenticated() || auth.Account.ID != id {
		serviceProfile.FilterVerifiedFields()
	}
	outProfile := dtoFromProfile(serviceProfile, t)
	WriteBody(w, r, outProfile)
}

func handleProfilePost(w http.ResponseWriter, r *http.Request) {
	t, err := getAccountType(r)
	if err != nil {
		restError(w, r, err, http.StatusBadRequest)
		return
	}
	id, err := getUUID(r, "uuid")
	if err != nil {
		restError(w, r, err, http.StatusBadRequest)
		return
	}
	profile := &ProfileRequestDTO{}
	if !ParseBody(w, r, profile) {
		return
	}

	serviceProfile := &services.Profile{
		AccountID:   id,
		Avatar:      profile.Avatar,
		FirstName:   profile.FirstName,
		LastName:    profile.LastName,
		City:        profile.City,
		Country:     profile.Country,
		Description: profile.Description,
	}
	if err := services.CreateProfile(serviceProfile); err != nil {
		restError(w, r, err, http.StatusInternalServerError)
		return
	}
	if ok, err := serviceProfile.IsAccountType(t); err != nil {
		restError(w, r, err, http.StatusInternalServerError)
		return
	} else if !ok {
		restError(w, r, errors.New("Account type does not match endpoint."), http.StatusBadRequest)
		return
	}
	outProfile := dtoFromProfile(serviceProfile, t)
	WriteBody(w, r, outProfile)
}

func handleProfileUpdateAvatar(w http.ResponseWriter, r *http.Request) {
	t, err := getAccountType(r)
	if err != nil {
		restError(w, r, err, http.StatusBadRequest)
		return
	}
	id, err := getUUID(r, "uuid")
	if err != nil {
		restError(w, r, err, http.StatusBadRequest)
		return
	}
	value := ParseUpdateString(w, r)
	if err := validateUpdate("Avatar", value, &ProfileRequestDTO{}); err != nil {
		restError(w, r, err, http.StatusBadRequest)
		return
	}
	var profile *services.Profile
	if profile, err = services.UpdateProfileField(id, "avatar", value); err != nil {
		restError(w, r, err, http.StatusBadRequest)
		return
	}
	if ok, err := profile.IsAccountType(t); err != nil {
		restError(w, r, err, http.StatusInternalServerError)
		return
	} else if !ok {
		restError(w, r, errors.New("Account type does not match endpoint."), http.StatusBadRequest)
		return
	}
	outProfile := dtoFromProfile(profile, t)
	WriteBody(w, r, outProfile)
}

func handleProfileUpdateFirstName(w http.ResponseWriter, r *http.Request) {
	t, err := getAccountType(r)
	if err != nil {
		restError(w, r, err, http.StatusBadRequest)
		return
	}
	id, err := getUUID(r, "uuid")
	if err != nil {
		restError(w, r, err, http.StatusBadRequest)
		return
	}
	value := ParseUpdateString(w, r)
	if err := validateUpdate("FirstName", value, &ProfileRequestDTO{}); err != nil {
		restError(w, r, err, http.StatusBadRequest)
		return
	}
	var profile *services.Profile
	if profile, err = services.UpdateProfileField(id, "first_name", value); err != nil {
		restError(w, r, err, http.StatusBadRequest)
		return
	}
	if err = profile.GenerateNewSlug(nil); err != nil {
		restError(w, r, err, http.StatusInternalServerError)
		return
	}
	if err = profile.Save(); err != nil {
		restError(w, r, err, http.StatusInternalServerError)
		return
	}
	if ok, err := profile.IsAccountType(t); err != nil {
		restError(w, r, err, http.StatusInternalServerError)
		return
	} else if !ok {
		restError(w, r, errors.New("Account type does not match endpoint."), http.StatusBadRequest)
		return
	}
	outProfile := dtoFromProfile(profile, t)
	WriteBody(w, r, outProfile)
}

func handleProfileUpdateLastName(w http.ResponseWriter, r *http.Request) {
	t, err := getAccountType(r)
	if err != nil {
		restError(w, r, err, http.StatusBadRequest)
		return
	}
	id, err := getUUID(r, "uuid")
	if err != nil {
		restError(w, r, err, http.StatusBadRequest)
		return
	}
	value := ParseUpdateString(w, r)
	if err := validateUpdate("LastName", value, &ProfileRequestDTO{}); err != nil {
		restError(w, r, err, http.StatusBadRequest)
		return
	}
	var profile *services.Profile
	if profile, err = services.UpdateProfileField(id, "last_name", value); err != nil {
		restError(w, r, err, http.StatusBadRequest)
		return
	}
	if err = profile.GenerateNewSlug(nil); err != nil {
		restError(w, r, err, http.StatusInternalServerError)
		return
	}
	if err = profile.Save(); err != nil {
		restError(w, r, err, http.StatusInternalServerError)
		return
	}
	if ok, err := profile.IsAccountType(t); err != nil {
		restError(w, r, err, http.StatusInternalServerError)
		return
	} else if !ok {
		restError(w, r, errors.New("Account type does not match endpoint."), http.StatusBadRequest)
		return
	}
	outProfile := dtoFromProfile(profile, t)
	WriteBody(w, r, outProfile)
}

func handleProfileUpdateCity(w http.ResponseWriter, r *http.Request) {
	t, err := getAccountType(r)
	if err != nil {
		restError(w, r, err, http.StatusBadRequest)
		return
	}
	id, err := getUUID(r, "uuid")
	if err != nil {
		restError(w, r, err, http.StatusBadRequest)
		return
	}
	value := ParseUpdateString(w, r)
	if err := validateUpdate("City", value, &ProfileRequestDTO{}); err != nil {
		restError(w, r, err, http.StatusBadRequest)
		return
	}
	var profile *services.Profile
	if profile, err = services.UpdateProfileField(id, "city", value); err != nil {
		restError(w, r, err, http.StatusBadRequest)
		return
	}
	if ok, err := profile.IsAccountType(t); err != nil {
		restError(w, r, err, http.StatusInternalServerError)
		return
	} else if !ok {
		restError(w, r, errors.New("Account type does not match endpoint."), http.StatusBadRequest)
		return
	}
	outProfile := dtoFromProfile(profile, t)
	WriteBody(w, r, outProfile)
}

func handleProfileUpdateCountry(w http.ResponseWriter, r *http.Request) {
	t, err := getAccountType(r)
	if err != nil {
		restError(w, r, err, http.StatusBadRequest)
		return
	}
	id, err := getUUID(r, "uuid")
	if err != nil {
		restError(w, r, err, http.StatusBadRequest)
		return
	}
	value := ParseUpdateString(w, r)
	if err := validateUpdate("Country", value, &ProfileRequestDTO{}); err != nil {
		restError(w, r, err, http.StatusBadRequest)
		return
	}
	var profile *services.Profile
	if profile, err = services.UpdateProfileField(id, "country", value); err != nil {
		restError(w, r, err, http.StatusBadRequest)
		return
	}
	if ok, err := profile.IsAccountType(t); err != nil {
		restError(w, r, err, http.StatusInternalServerError)
		return
	} else if !ok {
		restError(w, r, errors.New("Account type does not match endpoint."), http.StatusBadRequest)
		return
	}
	outProfile := dtoFromProfile(profile, t)
	WriteBody(w, r, outProfile)
}

func handleProfileUpdateDescription(w http.ResponseWriter, r *http.Request) {
	t, err := getAccountType(r)
	if err != nil {
		restError(w, r, err, http.StatusBadRequest)
		return
	}
	id, err := getUUID(r, "uuid")
	if err != nil {
		restError(w, r, err, http.StatusBadRequest)
		return
	}
	value := ParseUpdateString(w, r)
	if err := validateUpdate("Description", value, &ProfileRequestDTO{}); err != nil {
		restError(w, r, err, http.StatusBadRequest)
		return
	}
	var profile *services.Profile
	if profile, err = services.UpdateProfileField(id, "description", value); err != nil {
		restError(w, r, err, http.StatusBadRequest)
		return
	}
	if ok, err := profile.IsAccountType(t); err != nil {
		restError(w, r, err, http.StatusInternalServerError)
		return
	} else if !ok {
		restError(w, r, errors.New("Account type does not match endpoint."), http.StatusBadRequest)
		return
	}
	outProfile := dtoFromProfile(profile, t)
	WriteBody(w, r, outProfile)
}
