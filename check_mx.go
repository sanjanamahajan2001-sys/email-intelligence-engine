package main

import (
	"database/sql"
	"fmt"
	"os"

	_ "github.com/mattn/go-sqlite3"
)

func main() {
	db, err := sql.Open("sqlite3", "emails.db")
	if err != nil {
		fmt.Printf("Error opening DB: %v\n", err)
		os.Exit(1)
	}
	defer db.Close()

	rows, err := db.Query("SELECT signature, provider_name FROM disposable_mx_signatures")
	if err != nil {
		fmt.Printf("Error querying table: %v\n", err)
		os.Exit(1)
	}
	defer rows.Close()

	fmt.Println("MX Signatures in DB:")
	for rows.Next() {
		var sig, prov string
		rows.Scan(&sig, &prov)
		fmt.Printf(" - %s: %s\n", sig, prov)
	}
}
