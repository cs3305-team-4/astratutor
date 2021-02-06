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

}
