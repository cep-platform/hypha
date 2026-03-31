package provisoning

import (
	"fmt"
	"hypha/app/worker-provisioning/utils"
	"log"
	"math"
)

func SpinUp(mem float64) error {
	commands := provisoning.CommandsDarwin()	
	
	err := provisoning.RunCommand(commands.CheckPackageManager)
	if err != nil {
		return fmt.Errorf("spin-up failed: %s", err)
	}
	
	err = provisoning.RunCommand(commands.CheckDocker)
	if err != nil {
		return fmt.Errorf("spin-up failed: %s", err)
	}


	err = provisoning.RunCommand(commands.InstallDocker)
	if err != nil {
		return fmt.Errorf("spin-up failed: %s", err)
	}
	
	memAvailable, err := provisoning.GetAvailableMem()
	
	if err != nil {
		return fmt.Errorf("spin-up failed: %s", err)
	}
	
	memToReserve := fmt.Sprintf("%vG", math.Round(memAvailable * mem))
	diskToReserve := fmt.Sprintf("%vG", 4)

	log.Printf("All utilities correctly installed, reserving %v of memory now", memToReserve)	
	
	err = provisoning.RunCommand(commands.BuildImage)
	if err != nil {
		return fmt.Errorf("spin-up failed: %s", err)
	}

	commands = provisoning.SetDockerCommand(memToReserve, diskToReserve)
	
	

	err = provisoning.RunCommand(commands.SpinUpContainer)
	
	if err != nil {
		return fmt.Errorf("spin-up failed: %s", err)
	}

	return nil
}
