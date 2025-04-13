package main

import (
	"database/sql"
	"flag"
	"fmt"
	"log"
	"net/http"

	"github.com/CaiqueRibeiro/client-api-ex/server/src/gateways"
	"github.com/CaiqueRibeiro/client-api-ex/server/src/handlers"
	"github.com/CaiqueRibeiro/client-api-ex/server/src/repositories"
	_ "github.com/mattn/go-sqlite3"
)

func main() {
	// Analisa os flags da linha de comando
	port := flag.String("port", "8080", "HTTP server port")
	dbPath := flag.String("db", "./quotations.db", "Path to SQLite database file")
	flag.Parse()

	db, err := sql.Open("sqlite3", *dbPath)
	if err != nil {
		log.Fatalf("Failed to connect to SQLite database: %v", err)
	}
	defer db.Close()

	err = createDatabaseTable(db)
	if err != nil {
		log.Fatalf("Failed to create database table: %v", err)
	}

	quotationsRepository := repositories.NewQuotationsRepository(db)
	quotationGateway := gateways.NewQuotationGateway()
	quotationHandler := handlers.NewQuotationHandler(quotationGateway, quotationsRepository)

	mux := http.NewServeMux()
	mux.HandleFunc("GET /cotacao", quotationHandler.HandleGetQuotation)

	serverAddr := fmt.Sprintf(":%s", *port)

	log.Printf("Starting server on %s", serverAddr)
	log.Fatal(http.ListenAndServe(serverAddr, mux))
}

func createDatabaseTable(conn *sql.DB) error {
	_, err := conn.Exec(`CREATE TABLE IF NOT EXISTS quotations (
		id TEXT PRIMARY KEY,
		code TEXT,
		codein TEXT,
		name TEXT,
		high TEXT,
		low TEXT,
		varBid TEXT,
		pctChange TEXT,
		bid TEXT,
		ask TEXT,
		timestamp TEXT,
		create_date TEXT
	)`)
	return err
}
