package dns

import (
	"os/exec"
	"regexp"
	"fmt"
	"bytes"
	"log"
	"strings"
)

func GetGlobalDNS() ([]string, error) {
    out, err := exec.Command("scutil", "--dns").Output()
    if err != nil {
        return nil, err
    }

    ipRe := regexp.MustCompile(`nameserver\[.*?\]\s*:\s*(\S+)`)
    
    for _, line := range strings.Split(string(out), "\n") {
        matches := ipRe.FindStringSubmatch(line)
        if len(matches) > 1 {
            return []string{matches[1]}, nil // first nameserver in global config
        }
    }
    return nil, fmt.Errorf("no DNS found in scutil output")
}

func SetGlobalDNS(servers []string) error {
	
	for _, server := range servers {
		
		cmd := exec.Command("networksetup", "-setdnsservers", "Wi-Fi", server)

		//buffers for collecting output
		var stdoutBuff bytes.Buffer
		var stderrBuff bytes.Buffer

		cmd.Stdout = &stdoutBuff
		cmd.Stderr = &stderrBuff

	
		if err:= cmd.Run(); err != nil {
			stderrString := stderrBuff.String()
			log.Printf("Error setting DNS address %s", stderrString)
			return fmt.Errorf("error setting DNS address %s", stderrString)
		}
		
		log.Printf("successfully set DNS server %s", server)
	}
	return nil
}


