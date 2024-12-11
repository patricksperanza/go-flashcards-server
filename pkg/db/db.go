package db

import (
	"database/sql"
	"fmt"
	"go-flashcards-server/pkg/config"
	"log"

	_ "github.com/go-sql-driver/mysql"
)

var DB *sql.DB

type Deck struct {
	ID        int    `json:"id"`
	UserID    int    `json:"user_id"`
	Name      string `json:"name"`
	CreatedAt string `json:"created_at"`
}

func Init() {
	dataSourceName := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s", config.DBUsername, config.DBPassword, config.DBHost, config.DBPort, config.DBName)
	var err error
    DB, err = sql.Open("mysql", dataSourceName)
    if err != nil {
        log.Fatalf("Error connecting to database: %v", err)
    }
    err = DB.Ping()
    if err != nil {
        log.Fatalf("Error pinging database: %v", err)
    }
    fmt.Println("Database connected")
}

func CreateDeck(userID int, name string) (int64, error) {
	query := "INSERT INTO deck (user_id, name) VALUES (?, ?)"
	result, err := DB.Exec(query, userID, name)
	if err != nil {
		log.Printf("Error creating deck: %v", err)
		return 0, err
	}
	return result.LastInsertId()
}

func GetDecksByUser(userID int) ([]Deck, error) {
	query := "SELECT id, user_id, name, created_at FROM deck WHERE user_id = ?"
	rows, err := DB.Query(query, userID)
	if err != nil {
		log.Printf("Error retrieving decks: %v", err)
		return nil, err
	}
	defer rows.Close()

	var decks []Deck
	for rows.Next() {
		var deck Deck
		if err := rows.Scan(&deck.ID, &deck.UserID, &deck.Name, &deck.CreatedAt); err != nil {
			log.Printf("Error scanning deck row: %v", err)
			return nil, err
		}
		decks = append(decks, deck)
	}
	return decks, nil
}

func UpdateDeck(deckID int, name string) (int64, error) {
	query := "UPDATE deck SET name = ? WHERE id = ?"
	result, err := DB.Exec(query, name, deckID)
	if err != nil {
		log.Printf("Error update deck with ID %d: %v\n", deckID, err)
		return 0, err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		log.Printf("Error retrieving rows affected: %v", err)
		return 0, err
	}
	return rowsAffected, nil
}

func DeleteDeck(deckID int) (int64, error) {
	query := "DELETE FROM deck WHERE id = ?"
	result, err := DB.Exec(query, deckID)
	if err != nil {
		log.Printf("Error deleting deck with id %d: %v\n", deckID, err)
		return 0, err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		log.Printf("Error retrieving rows: %v", err)
		return 0, err
	}
	return rowsAffected, nil
}

func UpdateCard(cardID int, question, answer string) (int64, error) {
	query := "UPDATE card SET question = ?, answer = ? WHERE id = ?"
	result, err := DB.Exec(query, question, answer, cardID)
	if err != nil {
		log.Printf("Error update card with ID %d: %v\n", cardID, err)
		return 0, err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		log.Printf("Error retrieving rows affected: %v", err)
		return 0, err
	}
	return rowsAffected, nil
}

func DeleteCard(cardID int) (int64, error) {
	query := "DELETE FROM card WHERE id = ?"
	result, err := DB.Exec(query, cardID)
	if err != nil {
		log.Printf("Error deleting card with id %d: %v\n", cardID, err)
		return 0, err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		log.Printf("Error retrieving rows: %v", err)
		return 0, err
	}
	return rowsAffected, nil
}
