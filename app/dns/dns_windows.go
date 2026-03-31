package dns

import (
	"fmt"
	"log"
	"os/exec"
	"regexp"
	"strings"
)

func GetGlobalDNS() ([]string, error) {
	out, err := exec.Command("powershell", "-NoProfile", "-Command",
		`Get-DnsClientServerAddress -AddressFamily IPv4 | Where-Object { $_.ServerAddresses } | ForEach-Object { $_.ServerAddresses }`).Output()
	if err != nil {
		return nil, fmt.Errorf("failed to get DNS servers: %v", err)
	}

	ipRe := regexp.MustCompile(`(\d{1,3}\.\d{1,3}\.\d{1,3}\.\d{1,3})`)
	var servers []string
	for _, line := range strings.Split(string(out), "\n") {
		matches := ipRe.FindStringSubmatch(strings.TrimSpace(line))
		if len(matches) > 1 {
			servers = append(servers, matches[1])
		}
	}

	if len(servers) == 0 {
		return nil, fmt.Errorf("no DNS servers found")
	}
	return servers, nil
}

func SetGlobalDNS(servers []string) error {
	if len(servers) == 0 {
		return fmt.Errorf("no DNS servers provided")
	}

	interfaces, err := exec.Command("powershell", "-NoProfile", "-Command",
		`Get-NetAdapter | Where-Object { $_.Status -eq "Up" } | Select-Object -First 1 -ExpandProperty InterfaceIndex`).Output()
	if err != nil {
		return fmt.Errorf("failed to get network interface: %v", err)
	}

	interfaceIndex := strings.TrimSpace(string(interfaces))
	if interfaceIndex == "" {
		return fmt.Errorf("no active network interface found")
	}

	dnsServers := strings.Join(servers, ",")
	cmd := exec.Command("powershell", "-NoProfile", "-Command",
		fmt.Sprintf(`Set-DnsClientServerAddress -InterfaceIndex %s -ServerAddresses %s`,
			interfaceIndex, dnsServers))

	if err := cmd.Run(); err != nil {
		return fmt.Errorf("failed to set DNS servers: %v", err)
	}

	log.Printf("successfully set DNS servers: %v", servers)
	return nil
}
