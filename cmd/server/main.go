package main

import (
	"go-flashcards-server/pkg/config"
	"go-flashcards-server/pkg/db"
	"go-flashcards-server/pkg/handler"
	"go-flashcards-server/pkg/middleware"
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

func main() {
    config.Init()
    db.Init()

    r := mux.NewRouter()
    r.Use(middleware.CorsMiddleware)
    r.Use(middleware.LoggingMiddleware)
    r.HandleFunc("/signup", handler.SignUp).Methods("POST", "OPTIONS")
    r.HandleFunc("/login", handler.Login).Methods("POST", "OPTIONS")

    deckRouter := r.PathPrefix("/deck").Subrouter()
    deckRouter.Use(middleware.AuthMiddleware)
    deckRouter.HandleFunc("/create", handler.CreateDeck).Methods("POST")
    deckRouter.HandleFunc("", handler.GetDecks).Methods("GET", "OPTIONS")
    deckRouter.HandleFunc("/update/{deck_id}", handler.UpdateDeck).Methods("PUT")
    deckRouter.HandleFunc("/delete/{deck_id}", handler.DeleteDeck).Methods("DELETE")

    cardRouter := r.PathPrefix("/card").Subrouter()
    cardRouter.Use(middleware.AuthMiddleware)
    cardRouter.HandleFunc("/create", handler.CreateCard).Methods("POST")
    cardRouter.HandleFunc("/{deck_id}", handler.GetCardsByDeck).Methods("GET")
    cardRouter.HandleFunc("/update/{card_id}", handler.UpdateCard).Methods("PUT")
    cardRouter.HandleFunc("/delete/{card_id}", handler.DeleteCard).Methods("DELETE")
    
    log.Println("Server started on :8000")
    log.Fatal(http.ListenAndServe(":8000", r))
}
