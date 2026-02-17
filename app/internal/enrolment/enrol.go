package main

import (
	"hypha/app/pkg"
	"log"
	"os"
	// "runtime"
)

func Enrol() {
	
	nebulaExists := pkg.IfNebulaExists()
	
	if nebulaExists{
		//TODO: get enrolment payload from unzipped contents
		// pkg.NebulaStart()
		
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

		payload, err := pkg.ParseCertFolder(pkg.DESTINATION_FOLDER)	
			
		log.Printf("Payload parsed: %s", payload)
		
		pkg.NebulaStart()
		return 
	}

	mok := pkg.InstallNebula()	
	
	if mok != nil{
		// runtime.Breakpoint()
	}
}

func main() {
	Enrol()
}
