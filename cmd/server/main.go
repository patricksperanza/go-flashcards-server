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
    r.HandleFunc("/signup", handler.SignUp).Methods("POST")
    r.HandleFunc("/login", handler.Login).Methods("POST")

    deckRouter := r.PathPrefix("/decks").Subrouter()
    deckRouter.Use(middleware.AuthMiddleware)
    deckRouter.HandleFunc("", handler.CreateDeckHandler).Methods("POST")
    deckRouter.HandleFunc("", handler.GetDecksHandler).Methods("GET")

    cardRouter := r.PathPrefix("/cards").Subrouter()
    cardRouter.Use(middleware.AuthMiddleware)
    cardRouter.HandleFunc("/create", handler.CreateCardHandler).Methods("POST")
    cardRouter.HandleFunc("/{deck_id}", handler.GetCardsByDeckHandler).Methods("GET")
    cardRouter.HandleFunc("/update/{card_id}", handler.UpdateCardHandler).Methods("PUT")
    cardRouter.HandleFunc("/delete/{card_id}", handler.DeleteCardHandler).Methods("DELETE")
    
    log.Println("Server started on :8000")
    log.Fatal(http.ListenAndServe(":8000", r))
}
