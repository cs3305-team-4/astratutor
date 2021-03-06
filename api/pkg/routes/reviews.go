package routes

import (
	"encoding/json"
	"net/http"

	"github.com/cs3305-team-4/api/pkg/services"
	"github.com/gorilla/mux"
)

func InjectReviewsRoutes(subrouter *mux.Router) {
	// Reviews/Ratings
	subrouter.HandleFunc("/{tid}", handleReviewGetAll).Methods("GET")
	subrouter.HandleFunc("/{tid}/average", handleReviewAverage).Methods("GET")
	subrouter.HandleFunc("/{tid}/{rid}", handleReviewGetSingle).Methods("GET")
	subrouter.HandleFunc("/{tid}/author/{sid}", handleReviewGetByStudent).Methods("GET")

	authRoutes := subrouter.NewRoute().Subrouter()
	authRoutes.Use(authRequired)

	authRoutes.HandleFunc("/{tid}", handleReviewCreate).Methods("POST")
	authRoutes.HandleFunc("/{tid}/{rid}/rating", handleReviewUpdateRating).Methods("POST")
	authRoutes.HandleFunc("/{tid}/{rid}/comment", handleReviewUpdateComment).Methods("POST")
	authRoutes.HandleFunc("/{tid}/{rid}", handleReviewDelete).Methods("DELETE")
}

//TODO(james): Pagination
func handleReviewGetAll(w http.ResponseWriter, r *http.Request) {
	tid, err := getUUID(r, "tid")
	if err != nil {
		restError(w, r, err, http.StatusBadRequest)
		return
	}

	reviews, err := services.TutorAllReviews(tid)
	if err != nil {
		restError(w, r, err, http.StatusInternalServerError)
		return
	}
	if err = json.NewEncoder(w).Encode(reviews); err != nil {
		restError(w, r, err, http.StatusInternalServerError)
		return
	}
}

func handleReviewAverage(w http.ResponseWriter, r *http.Request) {
	tid, err := getUUID(r, "tid")
	if err != nil {
		restError(w, r, err, http.StatusBadRequest)
		return
	}

	avg, _ := services.TutorReviewsAverage(tid)
	WriteBody(w, r, avg)
}

func handleReviewGetSingle(w http.ResponseWriter, r *http.Request) {
	tid, err := getUUID(r, "tid")
	if err != nil {
		restError(w, r, err, http.StatusBadRequest)
		return
	}
	rid, err := getUUID(r, "rid")
	if err != nil {
		restError(w, r, err, http.StatusBadRequest)
		return
	}
	review, err := services.TutorSingleReview(tid, rid)
	if err != nil {
		restError(w, r, err, http.StatusBadRequest)
		return
	}
	WriteBody(w, r, review)
}

func handleReviewGetByStudent(w http.ResponseWriter, r *http.Request) {
	tid, err := getUUID(r, "tid")
	if err != nil {
		restError(w, r, err, http.StatusBadRequest)
		return
	}
	sid, err := getUUID(r, "sid")
	if err != nil {
		restError(w, r, err, http.StatusBadRequest)
		return
	}
	review, err := services.TutorReviewByStudent(tid, sid)
	if err != nil {
		restError(w, r, err, http.StatusBadRequest)
		return
	}
	WriteBody(w, r, review)
}

func handleReviewCreate(w http.ResponseWriter, r *http.Request) {
	tid, err := getUUID(r, "tid")
	if err != nil {
		restError(w, r, err, http.StatusBadRequest)
		return
	}

	authContext, err := ReadRequestAuthContext(r)
	if authContext.Account.Type == services.Tutor {
		restError(w, r, services.ReviewErrorStudentsOnly, http.StatusForbidden)
		return
	}

	res, err := services.HaveCompletedLesson(authContext.Account.ID, tid)
	if !res {
		restError(w, r, services.ReviewErrorNoCompletedLesson, http.StatusForbidden)
		return
	}

	create := &services.ReviewCreateDTO{}
	if !ParseBody(w, r, create) {
		return
	}

	err = services.CreateReview(&services.Review{
		Rating:           create.Rating,
		Comment:          create.Comment,
		TutorProfileID:   tid,
		StudentProfileID: authContext.Account.ID,
	})
	if err != nil {
		restError(w, r, err, http.StatusInternalServerError)
	}
}

func handleReviewUpdateRating(w http.ResponseWriter, r *http.Request) {
	rid, err := getUUID(r, "rid")
	if err != nil {
		restError(w, r, err, http.StatusBadRequest)
		return
	}

	authContext, err := ReadRequestAuthContext(r)
	if authContext.Account.Type == services.Tutor {
		restError(w, r, services.ReviewErrorStudentsOnly, http.StatusForbidden)
		return
	}

	update := &services.ReviewUpdateDTO{}
	if !ParseBody(w, r, update) {
		return
	}

	err = services.UpdateReviewRating(rid, update.Rating, authContext.Account.ID)
	if err != nil {
		restError(w, r, err, http.StatusInternalServerError)
		return
	}
}

func handleReviewUpdateComment(w http.ResponseWriter, r *http.Request) {
	rid, err := getUUID(r, "rid")
	if err != nil {
		restError(w, r, err, http.StatusBadRequest)
		return
	}

	authContext, err := ReadRequestAuthContext(r)
	if authContext.Account.Type == services.Tutor {
		restError(w, r, services.ReviewErrorStudentsOnly, http.StatusForbidden)
		return
	}

	update := &services.ReviewUpdateDTO{}
	if !ParseBody(w, r, update) {
		return
	}

	err = services.UpdateReviewComment(rid, update.Comment, authContext.Account.ID)
	if err != nil {
		restError(w, r, err, http.StatusInternalServerError)
		return
	}
}

func handleReviewDelete(w http.ResponseWriter, r *http.Request) {
	tid, err := getUUID(r, "tid")
	if err != nil {
		restError(w, r, err, http.StatusBadRequest)
		return
	}

	rid, err := getUUID(r, "rid")
	if err != nil {
		restError(w, r, err, http.StatusBadRequest)
		return
	}

	authContext, err := ReadRequestAuthContext(r)
	if authContext.Account.Type == services.Tutor {
		restError(w, r, services.ReviewErrorStudentsOnly, http.StatusForbidden)
		return
	}

	err = services.TutorDeleteReview(tid, rid, authContext.Account.ID)
}
