package routes

import (
	"encoding/json"
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

// Profile DTO.
type ProfileDTO struct {
	AccountID   string `json:"account_id" validate:"len=0"`
	ID          string `json:"id" validate:"len=0"`
	Avatar      string `json:"avatar" validate:"omitempty,base64"`
	Slug        string `json:"slug" validate:"len=0"`
	FirstName   string `json:"first_name" validate:"required"`
	LastName    string `json:"last_name" validate:"required"`
	City        string `json:"city" validate:"required"`
	Country     string `json:"country" validate:"required"`
	Description string `json:"description"`
}

// TutorDTO DTO.
type TutorDTO struct {
	ProfileDTO
	Qualifications []QualificationsDTO `json:"qualifications" validate:"len=0"`
	WorkExperience []WorkExperienceDTO `json:"work_experience" validate:"len=0"`
	Availability   []bool              `json:"availability" validate:"omitempty,len=336"`
}

// QualificationsDTO DTO.
type QualificationsDTO struct {
	ID       string `json:"id" validate:"len=0"`
	Field    string `json:"field" validate:"required"`
	Degree   string `json:"degree" validate:"required"`
	School   string `json:"school" validate:"required"`
	Verified bool   `json:"verified" validate:"eq=false"`
}

// WorkExperienceDTO DTO.
type WorkExperienceDTO struct {
	ID          string `json:"id" validate:"len=0"`
	Role        string `json:"role" validate:"required"`
	YearsExp    int    `json:"years_exp" validate:"required"`
	Description string `json:"description"`
	Verified    bool   `json:"verified" validate:"eq=false"`
}

func dtoFromProfile(p *services.Profile, accountType services.AccountType) interface{} {
	switch accountType {
	case services.Student:
		return &ProfileDTO{
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
		qualifications := make([]QualificationsDTO, 0)
		for _, val := range p.Qualifications {
			qualifications = append(qualifications, QualificationsDTO{
				ID:       val.ID.String(),
				Field:    val.Field,
				Degree:   val.Degree,
				School:   val.School,
				Verified: val.Verified,
			})
		}
		workExperience := make([]WorkExperienceDTO, 0)
		for _, val := range p.WorkExperience {
			workExperience = append(workExperience, WorkExperienceDTO{
				ID:          val.ID.String(),
				Role:        val.Role,
				YearsExp:    val.YearsExp,
				Description: val.Description,
				Verified:    val.Verified,
			})
		}
		return &TutorDTO{
			ProfileDTO: ProfileDTO{
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
	// TODO(ericm): Filter verified qualifications and work experience for other users.
	serviceProfile, err := services.ReadProfileByAccountID(id, nil)
	if err != nil {
		restError(w, r, err, http.StatusBadRequest)
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
	if err = json.NewEncoder(w).Encode(outProfile); err != nil {
		restError(w, r, err, http.StatusInternalServerError)
		return
	}
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
	profile := &ProfileDTO{}
	if !ParseBody(w, r, profile) {
		return
	}

	serviceProfile := &services.Profile{
		AccountID:   id,
		Avatar:      profile.Avatar,
		Slug:        profile.Slug,
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
	if err = json.NewEncoder(w).Encode(outProfile); err != nil {
		restError(w, r, err, http.StatusInternalServerError)
		return
	}
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
	if err := validateUpdate("Avatar", value, &ProfileDTO{}); err != nil {
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
	dto := dtoFromProfile(profile, t)
	if err = json.NewEncoder(w).Encode(dto); err != nil {
		restError(w, r, err, http.StatusInternalServerError)
		return
	}
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
	if err := validateUpdate("FirstName", value, &ProfileDTO{}); err != nil {
		restError(w, r, err, http.StatusBadRequest)
		return
	}
	var profile *services.Profile
	if profile, err = services.UpdateProfileField(id, "first_name", value); err != nil {
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
	dto := dtoFromProfile(profile, t)
	if err = json.NewEncoder(w).Encode(dto); err != nil {
		restError(w, r, err, http.StatusInternalServerError)
		return
	}
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
	if err := validateUpdate("LastName", value, &ProfileDTO{}); err != nil {
		restError(w, r, err, http.StatusBadRequest)
		return
	}
	var profile *services.Profile
	if profile, err = services.UpdateProfileField(id, "last_name", value); err != nil {
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
	dto := dtoFromProfile(profile, t)
	if err = json.NewEncoder(w).Encode(dto); err != nil {
		restError(w, r, err, http.StatusInternalServerError)
		return
	}
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
	if err := validateUpdate("City", value, &ProfileDTO{}); err != nil {
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
	dto := dtoFromProfile(profile, t)
	if err = json.NewEncoder(w).Encode(dto); err != nil {
		restError(w, r, err, http.StatusInternalServerError)
		return
	}
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
	if err := validateUpdate("Country", value, &ProfileDTO{}); err != nil {
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
	dto := dtoFromProfile(profile, t)
	if err = json.NewEncoder(w).Encode(dto); err != nil {
		restError(w, r, err, http.StatusInternalServerError)
		return
	}
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
	if err := validateUpdate("Description", value, &ProfileDTO{}); err != nil {
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
	dto := dtoFromProfile(profile, t)
	if err = json.NewEncoder(w).Encode(dto); err != nil {
		restError(w, r, err, http.StatusInternalServerError)
		return
	}
}
