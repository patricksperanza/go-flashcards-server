package handler

import (
	"encoding/json"
	"fmt"
	"go-flashcards-server/pkg/db"
	"go-flashcards-server/pkg/types"
	"go-flashcards-server/pkg/utils"
	"log"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
)

type DeckPayload struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

func CreateDeck(w http.ResponseWriter, r *http.Request) {
	userID, err := getUserIdFromContext(r)
	if err != nil {
		fmt.Printf("Error getting user: %v\n", err.Error())
		utils.HandleErrorResponse(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	var deck types.Deck
	if err := json.NewDecoder(r.Body).Decode(&deck); err != nil {
		utils.HandleErrorResponse(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	if deck.Name == "" {
		utils.HandleErrorResponse(w, "Deck name is required", http.StatusBadRequest)
		return
	}

	deckID, err := db.CreateDeck(userID, deck.Name)
	if err != nil {
		utils.HandleErrorResponse(w, "Failed to create deck", http.StatusInternalServerError)
		return
	}

	response := types.GCResponse[DeckPayload]{
		IsOK:    true,
		Message: "Deck Created",
		Payload: &DeckPayload{
			ID:   int(deckID),
			Name: deck.Name,
		},
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)
}

func GetDecks(w http.ResponseWriter, r *http.Request) {
	userID, err := getUserIdFromContext(r)
	if err != nil {
		log.Println("Error getting user: ", err.Error())
		utils.HandleErrorResponse(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	decks, err := db.GetDecksByUser(userID)
	if err != nil {
		utils.HandleErrorResponse(w, "Failed to retrieve decks", http.StatusInternalServerError)
		return
	}
	payload := mapToDeckPayload(decks)
	response := types.GCResponse[[]DeckPayload]{
		IsOK:    true,
		Message: "Decks retrieved",
		Payload: &payload,
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)
}

func mapToDeckPayload(decks []types.Deck) []DeckPayload {
	payload := make([]DeckPayload, len(decks))
	for i, d := range decks {
		payload[i] = DeckPayload{
			ID:   d.ID,
			Name: d.Name,
		}
	}
	return payload
}

func UpdateDeck(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	deckID, err := strconv.Atoi(params["deck_id"])
	if err != nil {
		utils.HandleErrorResponse(w, "Invalid deck id", http.StatusBadRequest)
		return
	}

	var deck DeckPayload
	if err = json.NewDecoder(r.Body).Decode(&deck); err != nil {
		utils.HandleErrorResponse(w, "Invalid request", http.StatusBadRequest)
		return
	}

	if deck.Name == "" {
		utils.HandleErrorResponse(w, "Name is required", http.StatusBadRequest)
		return
	}

	rowsAffected, err := db.UpdateDeck(deckID, deck.Name)
	if err != nil {
		utils.HandleErrorResponse(w, "Failed to update deck", http.StatusInternalServerError)
		return
	}

	if rowsAffected == 0 {
		utils.HandleErrorResponse(w, "Deck not found", http.StatusNotFound)
		return
	}

	response := types.GCResponse[DeckPayload]{
		IsOK:    true,
		Message: "Deck Updated Succesfully",
		Payload: &DeckPayload{
			ID:   deckID,
			Name: deck.Name,
		},
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)
}

func DeleteDeck(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	deckID, err := strconv.Atoi(params["deck_id"])
	if err != nil {
		utils.HandleErrorResponse(w, "Invalid deck ID", http.StatusBadRequest)
		return
	}

	rowsAffected, err := db.DeleteDeck(deckID)
	if err != nil {
		log.Printf("Error deleting deck %v\n", err)
		utils.HandleErrorResponse(w, "Error deleting deck", http.StatusInternalServerError)
		return
	}

	if rowsAffected == 0 {
		utils.HandleErrorResponse(w, "Deck not found", http.StatusNotFound)
		return
	}
	message := fmt.Sprintf("Deck %d deleted successfully", deckID)
	response := types.GCResponse[string]{
		IsOK:    true,
		Message: message,
		Payload: nil,
	}
	w.Header().Set("Content-Type", "application/json")
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
