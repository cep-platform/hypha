package main

import (
	"hypha/app/pkg"
	// "runtime"
)

func Enrol() {
	
	nebulaExists := pkg.IfNebulaExists()
	
	if nebulaExists{
		//TODO: get enrolment payload from unzipped contents
		// pkg.NebulaStart()
	}
	mok := pkg.InstallNebula()	
	
	if mok != nil{
		// runtime.Breakpoint()
	}
}

func main() {
	Enrol()
}
