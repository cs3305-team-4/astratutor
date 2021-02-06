package routes

import (
	"github.com/gorilla/mux"
)

func InjectStudentsRoutes(subrouter *mux.Router) {
	// Profile routes
	subrouter.HandleFunc("/{uuid}/profile", handleProfileGet).Methods("GET")

	accountResource := subrouter.PathPrefix("/{uuid}").Subrouter()
	accountResource.Use(authAccount())
	accountResource.HandleFunc("/profile", handleProfilePost).Methods("POST")

	// Profile update routes
	accountResource.HandleFunc("/profile/avatar", handleProfileUpdateAvatar).Methods("POST")
	accountResource.HandleFunc("/profile/first_name", handleProfileUpdateFirstName).Methods("POST")
	accountResource.HandleFunc("/profile/last_name", handleProfileUpdateLastName).Methods("POST")
	accountResource.HandleFunc("/profile/city", handleProfileUpdateCity).Methods("POST")
	accountResource.HandleFunc("/profile/country", handleProfileUpdateCountry).Methods("POST")
	accountResource.HandleFunc("/profile/description", handleProfileUpdateDescription).Methods("POST")

	// Lessons routes
	// subrouter.HandleFunc("/{uuid}/lessons", handleStudentsLessonsGet).Methods("GET")
	// subrouter.HandleFunc("/{uuid}/lessons", handleStudentsLessonsPost).Methods("POST")
}

// func handleStudentsLessonsGet(w http.ResponseWriter, r *http.Request) {
// 	id, err := getUUID(r, "uuid")
// 	if err != nil {
// 		restError(w, r, err, http.StatusBadRequest)
// 		return
// 	}
// 	serviceProfile, err := services.ReadProfileByAccountID(id, nil)
// 	if err != nil {
// 		restError(w, r, err, http.StatusBadRequest)
// 		return
// 	}
// 	outProfile := dtoFromProfile(serviceProfile, services.Student)
// 	if err = json.NewEncoder(w).Encode(outProfile); err != nil {
// 		restError(w, r, err, http.StatusInternalServerError)
// 		return
// 	}
// }

// func handleStudentsLessonsPost(w http.ResponseWriter, r *http.Request) {
// 	id, err := getUUID(r, "uuid")
// 	if err != nil {
// 		restError(w, r, err, http.StatusBadRequest)
// 		return
// 	}

// 	lessonRequest := &LessonRequestDTO{}
// 	if !ParseBody(w, r, lessonRequest) {
// 		return
// 	}

// 	tutor, err := services.ReadTutorByID(lessonRequest.RequesterID, nil)
// 	if err != nil {
// 		restError(w, r, err, http.StatusBadRequest)
// 		return
// 	}

// 	student, err := services.ReadStudentByID(id, nil)
// 	if err != nil {
// 		restError(w, r, err, http.StatusBadRequest)
// 		return
// 	}

// 	err = services.CreateLesson(student, tutor, lessonRequest.TimeStarts, lessonRequest.LessonDetail)
// 	if err != nil {
// 		restError(w, r, err, http.StatusBadRequest)
// 		return
// 	}
// }
