package main

import (
	"hypha/app/dns"
	"log"
)

func main() {
	servers, err := dns.GetGlobalDNS()
	if err != nil {
		log.Fatalf("Fatal! %w", err)
	}
	log.Printf("servers: %v", servers)
}
