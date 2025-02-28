package handler

import (
	"encoding/json"
	"errors"
	"fmt"
	"go-flashcards-server/pkg/config"
	"go-flashcards-server/pkg/db"
	"go-flashcards-server/pkg/types"
	"go-flashcards-server/pkg/utils"

	"log"
	"math/rand"
	"net/http"
	"strconv"
	"time"

	"github.com/dgrijalva/jwt-go"

	"golang.org/x/crypto/bcrypt"
)

type User struct {
	ID        int    `json:"id"`
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	Email     string `json:"email"`
	Password  string `json:"password"`
	Salt      string `json:"salt"`
}

type UserPayload struct {
	ID        int    `json:"id"`
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	Email     string `json:"email"`
}

type Credentials struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

func generateSalt() string {
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
		utils.HandleErrorResponse(w, "Invalid input", http.StatusBadRequest)
		return
	}

	salt := generateSalt()

	passwordWithSalt := user.Password + salt
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(passwordWithSalt), bcrypt.DefaultCost)
	if err != nil {
		utils.HandleErrorResponse(w, "Error signing up", http.StatusInternalServerError)
		return
	}

	_, err = db.DB.Exec("INSERT INTO user (first_name, last_name, email, password_hash, salt) VALUES (?, ?, ?, ?, ?)",
		user.FirstName, user.LastName, user.Email, hashedPassword, salt)
	if err != nil {
		log.Printf("Error creating user: %v\n", err)
		utils.HandleErrorResponse(w, "Error creating user", http.StatusInternalServerError)
		return
	}
	creds := Credentials{
		Email:    user.Email,
		Password: user.Password,
	}
	userPayload, err := loginUser(&creds, w)
	if err != nil {
		utils.HandleErrorResponse(w, err.Error(), http.StatusBadRequest)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(types.GCResponse[UserPayload]{
		IsOK:    true,
		Message: "Signed Up Successfully",
		Payload: userPayload,
	})

}

func Login(w http.ResponseWriter, r *http.Request) {
	var creds Credentials
	err := json.NewDecoder(r.Body).Decode(&creds)
	if err != nil {
		utils.HandleErrorResponse(w, "Invalid input", http.StatusBadRequest)
		return
	}
	userPayload, err := loginUser(&creds, w)
	if err != nil {
		utils.HandleErrorResponse(w, err.Error(), http.StatusBadRequest)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(types.GCResponse[UserPayload]{
		IsOK:    true,
		Message: "Logged In",
		Payload: userPayload,
	})
}

func loginUser(creds *Credentials, w http.ResponseWriter) (*UserPayload, error) {
	var user User
	var storedHashedPassword string
	err := db.DB.QueryRow(
		"SELECT id, first_name, last_name, email, password_hash, salt FROM user WHERE email = ?",
		creds.Email,
	).Scan(&user.ID, &user.FirstName, &user.LastName, &user.Email, &storedHashedPassword, &user.Salt)
	if err != nil {
		fmt.Printf("Error retrieving user: %v\n", err.Error())
		return nil, errors.New("Error retrieving user")
	}

	passwordWithSalt := creds.Password + user.Salt
	err = bcrypt.CompareHashAndPassword([]byte(storedHashedPassword), []byte(passwordWithSalt))
	if err != nil {
		log.Printf("Failed login attempt for user: %s", creds.Email)
		return nil, errors.New("Failed login attempt")
	}
	sub := strconv.Itoa(user.ID)
	expirationTime := time.Now().Add(60 * time.Minute)
	claims := &jwt.MapClaims{
		"sub": sub,
		"exp": expirationTime.Unix(),
		"iat": time.Now().Unix(),
	}
	jwtKey := config.JWTSecretKey
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString(jwtKey)
	if err != nil {
		fmt.Printf("Error generating token:  %v\n", err)
		return nil, errors.New("Server error")
	}

	http.SetCookie(w, &http.Cookie{
		Name:     "token",
		Value:    tokenString,
		Expires:  expirationTime,
		HttpOnly: false,
		SameSite: http.SameSiteNoneMode,
		Path:     "/",
		Secure:   true,
	})
	userPayload := UserPayload{
		ID:        user.ID,
		FirstName: user.FirstName,
		LastName:  user.LastName,
		Email:     user.Email,
	}
	return &userPayload, nil
}
