package handler

import (
	"encoding/json"
	"go-flashcards-server/pkg/db"
	"net/http"

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

func CreateCardHandler(w http.ResponseWriter, r *http.Request) {
    var card Card
    decoder := json.NewDecoder(r.Body)
    if err := decoder.Decode(&card); err != nil {
        http.Error(w, "Invalid request", http.StatusBadRequest)
        return
    }

    if card.DeckID == 0 || card.Question == "" || card.Answer == "" {
        http.Error(w, "Deck, Question, and Answer are required", http.StatusBadRequest)
        return
    }

    statement, err := db.DB.Prepare("INSERT INTO card (deck_id, question, answer) VALUES (?, ?, ?)")
    if err != nil {
        http.Error(w, "Error creating statement", http.StatusInternalServerError)
        return
    }
    defer statement.Close()

    res, err := statement.Exec(card.DeckID, card.Question, card.Answer)
    if err != nil {
        http.Error(w, "Error creating card", http.StatusInternalServerError)
        return
    }

    lastID, err := res.LastInsertId()
    if err != nil {
        http.Error(w, "Error getting card ID", http.StatusInternalServerError)
        return
    }
    card.ID = int(lastID)
    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(card)
}


func GetCardsByDeckHandler(w http.ResponseWriter, r *http.Request) {
    params := mux.Vars(r)
    deckID := params["deck_id"]

    if deckID == "" {
        http.Error(w, "Deck is required", http.StatusBadRequest)
        return
    }

    rows, err := db.DB.Query("SELECT id, deck_id, question, answer, created_at FROM card WHERE deck_id = ?", deckID)
    if err != nil {
        http.Error(w, "Error retrieving cards", http.StatusInternalServerError)
        return
    }
    defer rows.Close()

    var cards []Card
    for rows.Next() {
        var card Card
        if err := rows.Scan(&card.ID, &card.DeckID, &card.Question, &card.Answer, &card.CreatedAt); err != nil {
            http.Error(w, "Error scanning flashcard", http.StatusInternalServerError)
            return
        }
        cards = append(cards, card)
    }

    if len(cards) == 0 {
        http.Error(w, "No flashcards found for this deck", http.StatusNotFound)
        return
    }

    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(cards)
}



