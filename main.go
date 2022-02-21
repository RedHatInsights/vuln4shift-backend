package main

import (
	"app/database_admin"
	"log"
	"os"
)

func main() {
	if len(os.Args) > 1 {
		switch os.Args[1] {
		case "database_admin":
			database_admin.MigrateUp()
			return
		}
	}
	log.Fatal("Please specify service name as the first argument.\n")
}
