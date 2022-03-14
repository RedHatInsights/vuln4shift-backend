package main

import (
	"app/dbadmin"
	"app/manager"
	"app/pyxis"
	"log"
	"os"
)

func main() {
	if len(os.Args) > 1 {
		switch os.Args[1] {
		case "dbadmin":
			dbadmin.Start()
		case "manager":
			manager.Start()
		case "pyxis":
			pyxis.Start()
		default:
			log.Fatalf("Unknown service name: %s\n", os.Args[1])
		}
	} else {
		log.Fatal("Please specify service name as the first argument.\n")
	}
}
