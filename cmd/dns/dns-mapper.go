package main

import (
	"hypha/app/dns"
	"log"
)

func main() {
	servers, err := dns.GetGlobalDNS()
	if err != nil {
		log.Fatalf("Fatal! %s", err)
	}
	log.Printf("servers: %v", servers)

	err = dns.SetGlobalDNS([]string{
		"8.8.8.8",  //google
	})
	
	if err != nil {
		log.Fatalf("error when setting %s",err)
	}


	servers, err = dns.GetGlobalDNS()
	if err != nil {
		log.Fatalf("Fatal! %s", err)
	}

	log.Printf("servers after setting: %v", servers)

}
