package dns

import (
	"net"
  "github.com/godbus/dbus/v5"
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

