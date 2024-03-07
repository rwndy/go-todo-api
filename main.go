package main

import (
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"log"
	"net/http"
)

type User struct {
	Email    string `json:"Email"`
	Password string `json:"password"` // In a real application, passwords should be hashed
}

type ToDo struct {
	ID     int    `json:"id"`
	UserID int    `json:"userId"`
	Task   string `json:"task"`
	Status string `json:"status"` // "todo", "in_progress", "done"
}

type Users struct {
	NAME     string `json:"name"`
	EMAIL    string `json:"email"`
	PASSWORD string `json:"password"`
}

type ResponseRegister struct {
	STATUS string `json:"status"`
}

type ResponseLogin struct {
	STATUS string `json:"status"`
}

// var users []User
// var todos []ToDo
var payloadRegister Users
var payloadLogin User

func main() {
	router := mux.NewRouter()

	// Setup routes
	router.HandleFunc("/register", RegisterHandler).Methods("POST")
	router.HandleFunc("/login", LoginHandler).Methods("POST")
	router.HandleFunc("/todos", CreateToDoHandler).Methods("POST")
	router.HandleFunc("/todos/{id}", UpdateToDoHandler).Methods("PUT")
	router.HandleFunc("/todos/{id}", DeleteToDoHandler).Methods("DELETE")

	// Start server
	log.Fatal(http.ListenAndServe(":8080", router))
}

func connecetDB() *gorm.DB {
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local", "root", "", "localhost", "3306", "todo-db")
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		panic("failed to connect database")
	}

	fmt.Println(db)
	return db
}

func RegisterHandler(w http.ResponseWriter, r *http.Request) {
	// Implement user registration
	db := connecetDB()

	var response ResponseRegister

	err := json.NewDecoder(r.Body).Decode(&payloadRegister)
	if err != nil {
		response.STATUS = "failed to register!"

		w.Header().Set("Content-Type", "application/json")

		err := json.NewEncoder(w).Encode(response)
		if err != nil {
			return
		}

		return
	}

	err = db.Exec(`INSERT INTO users (NAME,EMAIL,PASSWORD) VALUES (?,?,?)`, payloadRegister.NAME, payloadRegister.EMAIL, payloadRegister.PASSWORD).Error

	if err != nil {
		fmt.Println(err)
		return
	}

	response.STATUS = "succes register!"

	w.Header().Set("Content-Type", "application/json")

	_ = json.NewEncoder(w).Encode(response)
	return
}

func LoginHandler(w http.ResponseWriter, r *http.Request) {
	// Implement user login
	db := connecetDB()
	var response ResponseLogin

	err := json.NewDecoder(r.Body).Decode(&payloadLogin)

	if err != nil {
		response.STATUS = "failed to login!"

		w.Header().Set("Content-Type", "application/json")

		err := json.NewEncoder(w).Encode(response)
		if err != nil {
			return
		}

		return
	}

	var userFromDB User

	err = db.Raw(`SELECT PASSWORD, EMAIL FROM users WHERE PASSWORD = ? AND EMAIL = ?`, payloadLogin.Password, payloadLogin.Email).Scan(&userFromDB).Error

	if err != nil {
		fmt.Println(err)
		response.STATUS = "password or email not valid!"
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(400)
		_ = json.NewEncoder(w).Encode(response)

		return
	}

	if payloadLogin.Password != userFromDB.Password {
		fmt.Println(err)
		response.STATUS = "password not valid!"
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(400)
		_ = json.NewEncoder(w).Encode(response)

		return
	}

	if payloadLogin.Email != userFromDB.Email {
		fmt.Println(err)
		response.STATUS = "Email not register!"
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(400)
		_ = json.NewEncoder(w).Encode(response)

		return
	}

	response.STATUS = "succes login!"

	w.Header().Set("Content-Type", "application/json")

	_ = json.NewEncoder(w).Encode(response)
	return
}

func CreateToDoHandler(w http.ResponseWriter, r *http.Request) {
	// Implement ToDo creation
}

func UpdateToDoHandler(w http.ResponseWriter, r *http.Request) {
	// Implement ToDo update
}

func DeleteToDoHandler(w http.ResponseWriter, r *http.Request) {
	// Implement ToDo deletion
}
