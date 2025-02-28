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

	_ "github.com/go-sql-driver/mysql"
	"github.com/gorilla/mux"
)

type Card struct {
	ID        int    `json:"id"`
	DeckID    int    `json:"deck_id"`
	Question  string `json:"question"`
	Answer    string `json:"answer"`
	CreatedAt string `json:"created_at"`
}

func CreateCard(w http.ResponseWriter, r *http.Request) {
	var card Card
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&card); err != nil {
		utils.HandleErrorResponse(w, "Invalid request", http.StatusBadRequest)
		return
	}

	if card.DeckID == 0 || card.Question == "" || card.Answer == "" {
		utils.HandleErrorResponse(w, "Deck, Question, and Answer are required", http.StatusBadRequest)
		return
	}

	statement, err := db.DB.Prepare("INSERT INTO card (deck_id, question, answer) VALUES (?, ?, ?)")
	if err != nil {
		utils.HandleErrorResponse(w, "Error creating statement", http.StatusInternalServerError)
		return
	}
	defer statement.Close()

	res, err := statement.Exec(card.DeckID, card.Question, card.Answer)
	if err != nil {
		utils.HandleErrorResponse(w, "Error creating card", http.StatusInternalServerError)
		return
	}

	lastID, err := res.LastInsertId()
	if err != nil {
		utils.HandleErrorResponse(w, "Error getting card ID", http.StatusInternalServerError)
		return
	}
	card.ID = int(lastID)
	response := types.GCResponse[Card]{
		IsOK:    true,
		Message: "Card Created",
		Payload: &card,
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)
}

func GetCardsByDeck(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	deckID := params["deck_id"]

	if deckID == "" {
		utils.HandleErrorResponse(w, "Deck is required", http.StatusBadRequest)
		return
	}

	rows, err := db.DB.Query("SELECT id, deck_id, question, answer, created_at FROM card WHERE deck_id = ?", deckID)
	if err != nil {
		utils.HandleErrorResponse(w, "Error retrieving cards", http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	var cards []Card
	for rows.Next() {
		var card Card
		if err := rows.Scan(&card.ID, &card.DeckID, &card.Question, &card.Answer, &card.CreatedAt); err != nil {
			log.Printf("Error scanning flashcard: %v", err.Error())
			utils.HandleErrorResponse(w, "Error scanning flashcard", http.StatusInternalServerError)
			return
		}
		cards = append(cards, card)
	}

	if len(cards) == 0 {
		utils.HandleErrorResponse(w, "No flashcards found for this deck", http.StatusNotFound)
		return
	}
	response := types.GCResponse[[]Card]{
		IsOK:    true,
		Message: "Cards Retrieved",
		Payload: &cards,
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)
}

func UpdateCard(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	cardID, err := strconv.Atoi(params["card_id"])
	if err != nil {
		utils.HandleErrorResponse(w, "Invalid card id", http.StatusBadRequest)
		return
	}
	var payload struct {
		Question string
		Answer   string
	}

	if err = json.NewDecoder(r.Body).Decode(&payload); err != nil {
		utils.HandleErrorResponse(w, "Invalid request", http.StatusBadRequest)
		return
	}

	if payload.Question == "" || payload.Answer == "" {
		utils.HandleErrorResponse(w, "Question and answer are required", http.StatusBadRequest)
		return
	}

	rowsAffected, err := db.UpdateCard(cardID, payload.Question, payload.Answer)
	if err != nil {
		utils.HandleErrorResponse(w, "Failed to update card", http.StatusInternalServerError)
		return
	}

	if rowsAffected == 0 {
		utils.HandleErrorResponse(w, "Card not found", http.StatusNotFound)
		return
	}
	message := fmt.Sprintf("Card %d updated", cardID)
	response := types.GCResponse[Card]{
		IsOK:    true,
		Message: message,
		Payload: &Card{
			ID:       cardID,
			Question: payload.Question,
			Answer:   payload.Answer,
		},
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)
}

func DeleteCard(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	cardID, err := strconv.Atoi(params["card_id"])
	if err != nil {
		utils.HandleErrorResponse(w, "Invalid card ID", http.StatusBadRequest)
		return
	}
	rowsAffected, err := db.DeleteCard(cardID)
	if err != nil {
		log.Printf("Error deleting card %v\n", err)
		utils.HandleErrorResponse(w, "Error deleting card", http.StatusInternalServerError)
		return
	}

	if rowsAffected == 0 {
		utils.HandleErrorResponse(w, "Card not found", http.StatusNotFound)
		return
	}

	message := fmt.Sprintf("Card %d deleted successfully", cardID)
	response := types.GCResponse[string]{
		IsOK:    true,
		Message: message,
		Payload: nil,
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)
}
