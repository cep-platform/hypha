package dns

import (
	"os/exec"
	"regexp"
	"fmt"
	"strings"
)

//TODO: from claude needs to be tested on mac @Sven
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



