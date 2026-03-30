package dns

import (
	"net"
  "github.com/godbus/dbus/v5"
	"bytes"
	"log"
	"os/exec"
	"fmt"
)

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

