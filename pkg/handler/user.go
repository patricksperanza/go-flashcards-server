package handler

import (
	"encoding/json"
	"fmt"
	"go-flashcards-server/pkg/config"
	"go-flashcards-server/pkg/db"
	"log"
	"math/rand"
	"net/http"
	"strconv"
	"time"

	"github.com/dgrijalva/jwt-go"

	"golang.org/x/crypto/bcrypt"
)

type User struct {
	ID 				int	`json:"id"`
    FirstName  		string `json:"first_name"`
    LastName   		string `json:"last_name"`
    Email      		string `json:"email"`
    Password       	string `json:"password"`
	Salt			string `json:"salt"`
}

type Credentials struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

func generateSalt() string {
	rand.Seed(time.Now().UnixNano())
	const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	var salt []byte
	for i := 0; i < 16; i++ {
		salt = append(salt, charset[rand.Intn(len(charset))])
	}
	return string(salt)
}

func SignUp(w http.ResponseWriter, r *http.Request) {
	var user User
	err := json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		http.Error(w, "Invalid input", http.StatusBadRequest)
		return
	}

	salt := generateSalt()

	passwordWithSalt := user.Password + salt
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(passwordWithSalt), bcrypt.DefaultCost)
	if err != nil {
		http.Error(w, "Error signing up", http.StatusInternalServerError)
		return
	}

	_, err = db.DB.Exec("INSERT INTO user (first_name, last_name, email, password_hash, salt) VALUES (?, ?, ?, ?, ?)",
		user.FirstName, user.LastName, user.Email, hashedPassword, salt)
	if err != nil {
		log.Printf("Error creating user: %v\n", err)
		http.Error(w, "Error creating user", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode("User created")
}


func Login(w http.ResponseWriter, r *http.Request) {
	jwtKey := config.JWTSecretKey
    var creds Credentials
    err := json.NewDecoder(r.Body).Decode(&creds)
    if err != nil {
        http.Error(w, "Invalid input", http.StatusBadRequest)
        return
    }

	var userID int
    var storedSalt, storedHashedPassword string
    err = db.DB.QueryRow("SELECT id, password_hash, salt FROM user WHERE email = ?", creds.Email).Scan(&userID, &storedHashedPassword, &storedSalt)
    if err != nil {
        http.Error(w, "User not found", http.StatusUnauthorized)
        return
    }

    passwordWithSalt := creds.Password + storedSalt
    err = bcrypt.CompareHashAndPassword([]byte(storedHashedPassword), []byte(passwordWithSalt))
    if err != nil {
		log.Printf("Failed login attempt for user: %s", creds.Email)
        http.Error(w, "Invalid credentials", http.StatusUnauthorized)
        return
    }
	sub := strconv.Itoa(userID)
	expirationTime := time.Now().Add(60 * time.Minute) 
    claims := &jwt.MapClaims{
		"sub": sub,
		"exp": expirationTime.Unix(),
		"iat": time.Now().Unix(),
    }
    token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
    tokenString, err := token.SignedString(jwtKey)
    if err != nil {
		fmt.Printf("Error generating token:  %v\n", err)
        http.Error(w, "Server error", http.StatusInternalServerError)
        return
    }

    json.NewEncoder(w).Encode(map[string]string{"token": tokenString})
}




