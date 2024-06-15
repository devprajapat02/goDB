package main

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
	_ "github.com/lib/pq"
)

type User struct {
	Id int `json:"id"`
	Name string `json:"name"`
	Email string `json:"email"`
}


func main() {

	fmt.Println("Server started...")

	router := mux.NewRouter()

	router.HandleFunc("/users", getUsersController).Methods("GET")
	
	http.ListenAndServe(":8080", router)
	
}


// Router controllers

func getUsersController(w http.ResponseWriter, r *http.Request) {

	json.NewEncoder(w).Encode(`{
		message: "Hello User!!"
	}`)
}
