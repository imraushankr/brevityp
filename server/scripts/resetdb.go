package main

import (
	"fmt"
	"os"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Usage: go run resetdb.go <db-file>")
		os.Exit(1)
	}

	dbFile := os.Args[1]
	if _, err := os.Stat(dbFile); err == nil {
		err := os.Remove(dbFile)
		if err != nil {
			fmt.Printf("Error removing database file: %v\n", err)
			os.Exit(1)
		}
		fmt.Println("Database file removed successfully")
	} else {
		fmt.Println("Database file does not exist, nothing to remove")
	}
}
