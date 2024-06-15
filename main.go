package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

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
	router.HandleFunc("/users/{id}", getUserController).Methods("GET")
	router.HandleFunc("/users/create", createUserController).Methods("POST")
	router.HandleFunc("/users/update", updateUserController).Methods("PUT")
	router.HandleFunc("/users/delete", deleteUserController).Methods("DELETE")
	
	http.ListenAndServe(":8080", router)
	
}


// Router controllers

func getUsersController(w http.ResponseWriter, r *http.Request) {

	users := getUsers()
	fmt.Println(users)
	json.NewEncoder(w).Encode(users)
}

func getUserController(w http.ResponseWriter, r *http.Request) {
	
	params := mux.Vars(r)
	id := params["id"]
	userId, _ := strconv.Atoi(id)

	var user User = getUser(userId)

	fmt.Println(user)
	json.NewEncoder(w).Encode(user)
}

func createUserController(w http.ResponseWriter, r *http.Request) {

	var userName, userEmail string = r.FormValue("name"), r.FormValue("email")


	user := createUser(userName, userEmail)
	fmt.Println(user)
	json.NewEncoder(w).Encode(user)
}

func updateUserController(w http.ResponseWriter, r *http.Request) {

	var userId, _ = strconv.Atoi(r.FormValue("id"))
	var userName , userEmail string = r.FormValue("name"), r.FormValue("email")
	var user User = updateUser(userId, userName, userEmail)

	fmt.Println(user)
	json.NewEncoder(w).Encode(user)
}

func deleteUserController(w http.ResponseWriter, r *http.Request) {

	var userId, _ = strconv.Atoi(r.FormValue("id"))
	if (deleteUser(userId)) {
		json.NewEncoder(w).Encode(`{
			message: "Row deleted successfully"
		}`)
	} else {
		json.NewEncoder(w).Encode(`{
			message: "Could'nt perform deletion"
		}`)
	}
}


// Database operations

func connectDatabase() (*sql.DB) {
	databaseURI := "user=postgres password=test@123 dbname=GoDB sslmode=disable"
	db, err := sql.Open("postgres", databaseURI)
	if err != nil {
		panic(err)
	}
	return db
}

func getUsers() ([]User) {
	db := connectDatabase()

	rows, err := db.Query("SELECT * FROM users")

	if err != nil {
		panic(err)
		// var users []User
		// return users
	}
	
	var users []User
	for rows.Next() {
		var userId int
		var userName, userEmail string
		rows.Scan(&userId, &userName, &userEmail)
		var user User = User{
			Id: userId,
			Name: userName,
			Email: userEmail,
		}

		users = append(users, user)
	}

	return users
}

func getUser(id int) (User) {
	db := connectDatabase()

	rows := db.QueryRow("SELECT * FROM users WHERE id = $1", id)

	var user User
	rows.Scan(&user.Id, &user.Name, &user.Email)

	return user
}

func createUser(userName, userEmail string) (User) {
	db := connectDatabase()

	rowCountQuery, _ := db.Query("SELECT COUNT(*) FROM users")
	var rowCount int
	rowCountQuery.Next()
	rowCountQuery.Scan(&rowCount)

	_, err := db.Exec("INSERT INTO users (id, name, email) VALUES ($1, $2, $3)", rowCount + 1, userName, userEmail)
	if err != nil {
		fmt.Println(rowCount)
		panic(err)
	}

	return getUser(rowCount + 1)
}

func updateUser(userId int, newName, newEmail string) (User) {
	db := connectDatabase()
	
	var user User = getUser(userId)
	if (newName == "") {
		newName = user.Name
	}

	if (newEmail == "") {
		newEmail = user.Email
	}

	_, err := db.Exec("UPDATE users SET name = $1, email = $2 WHERE id = $3", newName, newEmail, userId)
	if err != nil {
		panic(err)
	}

	return getUser(userId)
}

func deleteUser(userId int) (bool) {
	db := connectDatabase()

	_, err := db.Exec("DELETE FROM users WHERE id = $1", userId)

	return err == nil
}