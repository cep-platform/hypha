package main

import (
	"hypha/app/dns"
	"fmt"
	"log"
)

func main() {
	dnsAddr, err := dns.GetGlobalDNS()
	if err!= nil {
		fmt.Println("addr: %w, %w", err)
	}
	log.Printf("Dns Address is: %s", dnsAddr)
}
