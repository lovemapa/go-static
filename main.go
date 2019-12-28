package main

import (
	"github.com/gorilla/mux"
	"log"
	"net/http"
)

func main() {

	route := mux.NewRouter()

	s := route.PathPrefix("/api").Subrouter() //Base Path
	route.PathPrefix("/").Handler(http.FileServer(http.Dir("./static")))

	//ROUTES

	s.HandleFunc("/createProfile", createProfile).Methods("POST") //upload File handler
	s.HandleFunc("/getAllUsers", getAllUsers).Methods("GET")      //display All list of users

	log.Print("Server Connected 🚀 ")
	log.Fatal(http.ListenAndServe(":8000", route)) // Run Server
}
