package dns

import (
	"bytes"
	"fmt"
	"log"
	"os/exec"
	"regexp"
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

func getActiveInterface() (string, error) {
	out, err := exec.Command("networksetup", "-listallhardwareports").Output()
	if err != nil {
		return "", fmt.Errorf("failed to list hardware ports: %v", err)
	}

	lines := strings.Split(string(out), "\n")
	for _, line := range lines {
		if strings.HasPrefix(line, "Device:") {
			device := strings.TrimSpace(strings.TrimPrefix(line, "Device:"))
			out, err := exec.Command("ifconfig", device).Output()
			if err != nil {
				continue
			}
			if strings.Contains(string(out), "inet ") && !strings.Contains(string(out), "inet6 ") {
				return device, nil
			}
		}
	}
	return "", fmt.Errorf("no active network interface found")
}

func SetGlobalDNS(servers []string) error {
	if len(servers) == 0 {
		return fmt.Errorf("no DNS servers provided")
	}

	iface, err := getActiveInterface()
	if err != nil {
		return fmt.Errorf("failed to get active interface: %v", err)
	}

	args := append([]string{"-setdnsservers", iface}, servers...)
	cmd := exec.Command("networksetup", args...)

	var stderrBuff bytes.Buffer
	cmd.Stderr = &stderrBuff

	if err := cmd.Run(); err != nil {
		return fmt.Errorf("error setting DNS servers: %s", stderrBuff.String())
	}

	log.Printf("successfully set DNS servers on %s: %v", iface, servers)
	return nil
}
