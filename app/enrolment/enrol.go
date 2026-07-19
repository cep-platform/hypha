package main

import (
	"bufio"
	"fmt"
	"hypha/app/pkg"
	"io"
	"log"
	"os"
	"strings"
)

//NOTE: This is mainly used for development currently
//Main entry point will be via cmd/widget
func Enrol() io.ReadCloser {
	
	nebulaExists := pkg.IfNebulaExists()
	
	if !nebulaExists{
		
		err := pkg.InstallNebula()

		if err != nil {
			log.Fatal("Failed to install Nebula")
			os.Exit(1)
		}
	 }

	err := pkg.ValidateDir(pkg.DIRS)
		
	if err != nil {
		log.Fatal(err)
		os.Exit(1)
	}

	err = pkg.Unzip(pkg.HOST_PATH,pkg.DESTINATION_FOLDER)
		
	if err != nil {
		log.Fatal("Failed to unzip: %w", err)
		os.Exit(1)
	}
	
	fmt.Print("Enter sudo password: ")
	passwordBytes, err := bufio.NewReader(os.Stdin).ReadString('\n')
	if err != nil {
		log.Fatal("Failed to read password: ", err)
	}
	sudoPassword := strings.TrimRight(passwordBytes, "\r\n")

	pipe, err := pkg.NebulaStart(pkg.NEBULA_PATH, pkg.DESTINATION_CERTS, sudoPassword)
	
	if err != nil {
		log.Fatal("Failed to start Nebula: %w", err)
		os.Exit(1) 
	}
	
	return pipe
}

func main() {
	Enrol()
}
