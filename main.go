package main

import (
	"app/dbadmin"
	"log"
	"os"
)

func main() {
	if len(os.Args) > 1 {
		switch os.Args[1] {
		case "dbadmin":
			dbadmin.MigrateUp()
			return
		}
	}
	log.Fatal("Please specify service name as the first argument.\n")
}
