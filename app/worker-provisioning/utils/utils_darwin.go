package provisoning

import (
	"bytes"
	"fmt"
	"log"
	"os/exec"

	"github.com/shirou/gopsutil/v3/mem"
)


func RunCommand(command Command) error {
	cmd := exec.Command(command.Header, command.Commands...)

	//buffers for debugging
	var stdoutBuff bytes.Buffer
	var stderrBuff bytes.Buffer

	cmd.Stdout = &stdoutBuff
	cmd.Stderr = &stderrBuff

	err := cmd.Run()

	if err != nil {
		stderrString := stderrBuff.String()
		log.Printf("error encountered when %s: %s",command.ID, stderrString)
		return fmt.Errorf("error encountered when %s: %s ", command.ID, stderrString)
	}

	log.Printf("%s command successfully executed: %s", command.ID, stdoutBuff.String())
	return nil
}

func GetAvailableMem() (float64, error) {
	v, err := mem.VirtualMemory()
	if err != nil {
		return 0.0, fmt.Errorf("could not query available memory", err) 
	}
	memValue := float64(v.Available / 1e9)
	fmt.Printf("Available: %.2f GB\n", memValue)
	
	return memValue, err
}
