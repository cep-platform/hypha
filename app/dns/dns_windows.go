package dns

import (
    "golang.org/x/sys/windows/registry"
    "strings"
)

//TODO: from Claude needs to be tested on windows machine @Me
func GetGlobalDNS() ([]string, error) {
    k, err := registry.OpenKey(registry.LOCAL_MACHINE,
        `SYSTEM\CurrentControlSet\Services\Tcpip\Parameters`,
        registry.QUERY_VALUE)
    if err != nil {
        return nil, err
    }
    defer k.Close()

    val, _, err := k.GetStringValue("NameServer")
    if err != nil {
        return nil, err
    }

    // NameServer can be comma or space separated
    servers := strings.FieldsFunc(val, func(r rune) bool {
        return r == ',' || r == ' '
    })
    return servers, nil
}
