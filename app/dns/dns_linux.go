package dns

import (
	"bytes"
	"fmt"
	dns "hypha/app/dns/utils"
	"log"
	"net"
	"os/exec"
	"runtime"
	"strings"

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
		log.Printf("error encountered when %s: %s",command.ID, stderrString)
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
  	dnsArr:= result.Value().([]interface{})
		dnsArrIdx := len(dnsArr) - 1
		address := dnsArr[dnsArrIdx].([]byte)
		
		ip := net.IP(address).String()
			
		return append(dnsArrTemp, ip), nil
}

//TODO: Tomorrow: check if NMCLI is installed
//Do greps:
//	- get which one is active
//	- nmcli connection show --active | grep ethernet/Wi-fi
//	- sudo nmcli connection modify "ethernet" ipv4.dns "8.8.8.8 8.8.4.4"
//  - That should be it
func SetGlobalDNS(servers []string) error {
	
	commands := dns.CommandsLinux()	
	
	err := RunCommand(commands.CheckPackageManager)
	
	if err != nil {
		return fmt.Errorf("Package manager chekc failed: %s", err)
	}
	
	err = RunCommand(commands.CheckNetworkUtil)
	if err != nil {
		return fmt.Errorf("network util check failed: %s", err)
	}


	err = RunCommand(commands.StartNetworkUtil)

	if err != nil {
		return fmt.Errorf("start of network util failed: %s", err)
	}
	
	//hacky af	
	out, err := exec.Command(commands.RetrieveActiveConnection.Header, "-c", commands.RetrieveActiveConnection.Commands[0]).Output()

	if err != nil {
		return fmt.Errorf("spin-up failed: %s", err)
	}
	
	//still hacky, we assume the first one to be the active dns
	activeInterface := strings.Split(fmt.Sprintf("%s", out), "\n")[0]
	
	for _, server := range servers {
		//fuck it
		runtime.Breakpoint()
		out, err = exec.Command("nmcli", "connection", "modify", activeInterface, "ipv4.dns", server).Output()
		
		if err != nil {
			return fmt.Errorf("failed to modify DNS: %s", err)
		}

		out, err = exec.Command("nmcli", "connection", "up", activeInterface).Output()

		if err != nil {
			return fmt.Errorf("Failed to validate conn %s", err)
		}


		log.Printf("successfully set DNS server %s", out)
	}
	return nil
}


