package main

import (
	"log"

	"atonixcorp/cloud/devtools/portal/backend"
)

func main() {
	if err := backend.Start(":8081"); err != nil {
		log.Fatal(err)
	}
}
