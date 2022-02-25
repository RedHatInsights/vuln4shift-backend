package main

import (
	"app/dbadmin"
	"app/manager"
	"log"
	"os"
)

func main() {
	if len(os.Args) > 1 {
		switch os.Args[1] {
		case "dbadmin":
			dbadmin.MigrateUp()
		case "manager":
			manager.Start()
		}
		return
	}
	log.Fatal("Please specify service name as the first argument.\n")
}
