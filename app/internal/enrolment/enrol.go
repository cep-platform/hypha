package main

import (
	"hypha/app/pkg"
	"log"
	"os"
	"runtime"
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

		pkg.Unzip(pkg.HOST_NAME,pkg.DESTINATION_CERT_PATH)
		runtime.Breakpoint()	
	
	}
	mok := pkg.InstallNebula()	
	if mok != nil{
		// runtime.Breakpoint()
	}
}

func main() {
	Enrol()
}
