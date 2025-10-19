// main package of the app
package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"

	"github.com/gin-gonic/gin"
	_ "github.com/lib/pq" 
)

type TextEntry struct {
	ID   int    `json:"id"`
	Text string `json:"text" binding:"required"`
}

var db *sql.DB

func main() {
	err := setupDatabase()
	if err != nil {
		log.Fatalf("Failed to setup database: %v", err)
	}
	defer db.Close()

	router := gin.Default()

	router.POST("/text", postText(db))
	router.GET("/text/:id", getText(db))

	log.Println("Server running on http://localhost:8080")
	if err := router.Run(":8080"); err != nil {
		log.Fatal(err)
	}
}

func setupDatabase() error {
	connStr := os.Getenv("DATABASE_URL")
	if connStr == "" {
		log.Println("DATABASE_URL not set, using default connection string.")
		connStr = "user=postgres password=mysecretpassword dbname=simpleapp sslmode=disable port=1200"
	}

	var err error
	db, err = sql.Open("postgres", connStr)
	if err != nil {
		return fmt.Errorf("error opening database: %w", err)
	}

	if err = db.Ping(); err != nil {
		return fmt.Errorf("error pinging database: %w", err)
	}

	log.Println("Successfully connected to PostgreSQL!")

	createTableSQL := `
	CREATE TABLE IF NOT EXISTS texts (
		id SERIAL PRIMARY KEY,
		content TEXT NOT NULL
	);`

	_, err = db.Exec(createTableSQL)
	if err != nil {
		return fmt.Errorf("error creating table: %w", err)
	}

	log.Println("Database table 'texts' ensured.")
	return nil
}

func postText(db *sql.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		var newEntry TextEntry
		if err := c.ShouldBindJSON(&newEntry); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body or missing 'text' field"})
			return
		}

		sqlStatement := `INSERT INTO texts (content) VALUES ($1) RETURNING id`
		var id int
		err := db.QueryRow(sqlStatement, newEntry.Text).Scan(&id)

		if err != nil {
			log.Printf("DB INSERT error: %v", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not store text"})
			return
		}

		c.JSON(http.StatusCreated, gin.H{"id": id, "message": "Text stored successfully"})
	}
}

func getText(db *sql.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		idStr := c.Param("id")
		id, err := strconv.Atoi(idStr)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID format"})
			return
		}

		sqlStatement := `SELECT id, content FROM texts WHERE id = $1`
		row := db.QueryRow(sqlStatement, id)

		var entry TextEntry
		err = row.Scan(&entry.ID, &entry.Text)

		if err != nil {
			if err == sql.ErrNoRows {
				c.JSON(http.StatusNotFound, gin.H{"error": fmt.Sprintf("Text with id %d not found", id)})
				return
			}
			log.Printf("DB SELECT error: %v", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not retrieve text"})
			return
		}

		c.JSON(http.StatusOK, entry)
	}
}
