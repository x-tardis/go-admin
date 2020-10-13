package main

import (
	"log"

	"github.com/x-tardis/go-admin/pkg/infra"
)

func main() {
	log.Println(infra.ActiveNetwork())
	log.Println(infra.GetNetInformation())
}
