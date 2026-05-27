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

	_, err = db.Exec("DELETE FROM disposable_domains;")
	if err != nil {
		fmt.Printf("Error deleting table: %v\n", err)
		os.Exit(1)
	}

	fmt.Println("✅ Successfully cleared disposable_domains table!")
}
