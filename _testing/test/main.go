package main

import (
	"log"

	"github.com/x-tardis/go-admin/pkg/deployed"
)

func main() {
	log.Println(deployed.IPLocation("117.136.75.9"))
}
