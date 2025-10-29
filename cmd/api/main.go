package main

import (
	"bookstore/configs"
	"bookstore/internal/db"
	"bookstore/internal/router"
	"encoding/json"
	"fmt"
	"log"
	"os"

	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

func main() {
	var err error

	if err := godotenv.Load(".env"); err != nil {
		log.Fatalf("Error loading .env file: %s", err)
	}

	// Read configuration from file
	configFile, err := os.Open("configs/config.json")
	if err != nil {
		log.Fatal("Error opening config file:", err)
	}
	defer configFile.Close()

	var config configs.Config
	jsonParser := json.NewDecoder(configFile)
	err = jsonParser.Decode(&config)
	if err != nil {
		log.Fatal("Error decoding config JSON:", err)
	}

	// Construct connection string
	connStr := fmt.Sprintf(`host=%s
							port=%d 
							user=%s 
							password=%s 
							dbname=%s 
							sslmode=%s`,
		config.Database.Host,
		config.Database.Port,
		config.Database.User,
		os.Getenv("DB_PASSWORD"),
		config.Database.DBName,
		config.Database.SSLMode)

	// Connect to PostgreSQL
	db.InitDB(connStr)
	defer db.DB.Close()

	router.RunRouters()
}
