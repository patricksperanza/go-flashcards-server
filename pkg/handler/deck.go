package handler

import (
	"encoding/json"
	"fmt"
	"go-flashcards-server/pkg/db"
	"log"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
)

func CreateDeck(w http.ResponseWriter, r *http.Request) {
	userID, err := getUserIdFromContext(r)
	if err != nil {
		fmt.Printf("Error getting user: %v\n", err.Error())
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	var payload struct {
		Name string `json:"name"`
	}
	if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	if payload.Name == "" {
		http.Error(w, "Deck name is required", http.StatusBadRequest)
		return
	}

	deckID, err := db.CreateDeck(userID, payload.Name)
	if err != nil {
		http.Error(w, "Failed to create deck", http.StatusInternalServerError)
		return
	}

	response := map[string]interface{}{
		"message": "Deck created successfully",
		"deck_id": deckID,
	}
	json.NewEncoder(w).Encode(response)
}

func GetDecks(w http.ResponseWriter, r *http.Request) {
	log.Println("GetDecks()")
	userID, err := getUserIdFromContext(r)
	if err != nil {
		log.Println("Error getting user: ", err.Error())
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	decks, err := db.GetDecksByUser(userID)
	if err != nil {
		http.Error(w, "Failed to retrieve decks", http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(decks)
}

func UpdateDeck(w http.ResponseWriter, r *http.Request) {
    params := mux.Vars(r)
    deckID, err := strconv.Atoi(params["deck_id"])
    if err != nil {
        http.Error(w, "Invalid deck id", http.StatusBadRequest)
        return
    }
    var payload struct {
        Name string `json: "name"`
    }

    if err = json.NewDecoder(r.Body).Decode(&payload); err != nil {
        http.Error(w, "Invalid request", http.StatusBadRequest)
		return
    }

    if payload.Name == "" {
		http.Error(w, "Name is required", http.StatusBadRequest)
		return
	}

    rowsAffected, err := db.UpdateDeck(deckID, payload.Name)
	if err != nil {
		http.Error(w, "Failed to update deck", http.StatusInternalServerError)
		return
	}

    if rowsAffected == 0 {
		http.Error(w, "Deck not found", http.StatusNotFound)
		return
	}

    response := map[string]string{"message": "Deck updated successfully"}
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)
}

func DeleteDeck(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	deckID, err := strconv.Atoi(params["deck_id"])
    if err != nil {
        http.Error(w, "Invalid deck ID", http.StatusBadRequest)
        return
    }

	rowsAffected, err := db.DeleteDeck(deckID)
    if err != nil {
        log.Printf("Error deleting deck %v\n", err)
        http.Error(w, "Error deleting deck", http.StatusInternalServerError)
        return
    }

	if rowsAffected == 0 {
        http.Error(w, "Deck not found", http.StatusNotFound)
        return
    }

	
    response := map[string]string{"message": "Deck deleted successfully"}
    w.WriteHeader(http.StatusOK)
    json.NewEncoder(w).Encode(response)
}

func getUserIdFromContext(r *http.Request) (int, error) {
	ctx := r.Context()
	userID, ok := ctx.Value("userID").(int)
	if !ok {
		return 0, http.ErrNoCookie
	}
	return userID, nil
}