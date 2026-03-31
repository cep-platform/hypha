package dns

import (
	"bytes"
	"fmt"
	dns "hypha/app/dns/utils"
	"log"
	"net"
	"os/exec"
	"strings"
	"runtime"
	"github.com/godbus/dbus/v5"
)

func RunCommand(command dns.Command) error {
	cmd := exec.Command(command.Header, command.Commands...)

	//buffers for debugging
	var stdoutBuff bytes.Buffer
	var stderrBuff bytes.Buffer

	cmd.Stdout = &stdoutBuff
	cmd.Stderr = &stderrBuff

	err := cmd.Run()

	if err != nil {
		stderrString := stderrBuff.String()
		log.Printf("error encountered when %s: %s", command.ID, stderrString)
		return fmt.Errorf("error encountered when %s: %s ", command.ID, stderrString)
	}

	log.Printf("%s command successfully executed: %s", command.ID, stdoutBuff.String())
	return nil
}

func GetGlobalDNS() ([]string, error) {

	//NOTE: this will be made consistent just wrapping in arr for now
	dnsArrTemp := make([]string, 1)
	var result dbus.Variant
	conn, err := dbus.SystemBus()
	if err != nil {
		return dnsArrTemp, err
	}

	defer conn.Close()

	obj := conn.Object("org.freedesktop.resolve1", "/org/freedesktop/resolve1")

	err = obj.Call("org.freedesktop.DBus.Properties.Get", 0,
		"org.freedesktop.resolve1.Manager", "CurrentDNSServer").Store(&result)

	if err != nil {
		return dnsArrTemp, err
	}

	//hacky to get indexable arr
	dnsArr := result.Value().([]interface{})
	dnsArrIdx := len(dnsArr) - 1
	address := dnsArr[dnsArrIdx].([]byte)

	ip := net.IP(address).String()

	return append(dnsArrTemp, ip), nil
}

// TODO: Tomorrow: check if NMCLI is installed
// Do greps:
//   - get which one is active
//   - nmcli connection show --active | grep ethernet/Wi-fi
//   - sudo nmcli connection modify "ethernet" ipv4.dns "8.8.8.8 8.8.4.4"
//   - That should be it
func SetGlobalDNS(servers []string) error {
	if len(servers) == 0 {
		return fmt.Errorf("no DNS servers provided")
	}

	commands := dns.CommandsLinux()

	err := RunCommand(commands.CheckNetworkUtil)
	if err != nil {
		return fmt.Errorf("network util check failed: %s", err)
	}

	out, err := exec.Command("bash", "-c", commands.RetrieveActiveConnection.Commands[0]).Output()
	if err != nil {
		return fmt.Errorf("failed to get active connection: %s", err)
	}

	//fuck it we usnafe indexing, because big pickle is a little bitch
	activeInterface := strings.Split(string(out), "\n")[0]
	if activeInterface == "" {
		return fmt.Errorf("no active network connection found")
	}

	dnsArg := strings.Join(servers, " ")
	runtime.Breakpoint()
	modifyCmd := exec.Command("nmcli", "connection", "modify", activeInterface, "ipv4.dns", dnsArg)
	if err := modifyCmd.Run(); err != nil {
		return fmt.Errorf("failed to modify DNS: %s", err)
	}

	upCmd := exec.Command("nmcli", "connection", "up", activeInterface)
	
	if err := upCmd.Run(); err != nil {
		return fmt.Errorf("failed to apply DNS settings: %s", err)
	}

	log.Printf("successfully set DNS servers: %s", dnsArg)
	return nil
}
