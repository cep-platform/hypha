package dns

type Commands struct {
	CheckPackageManager Command
	InstallPackageManager Command
	CheckNetworkUtil Command
	InstallNetworkUtil Command
	StartNetworkUtil Command
	RetrieveActiveConnection Command
	SetNewDNS Command
}

type Command struct {
		Header string
		ID string
		Commands []string
}

func CommandsLinux() *Commands {
	return &Commands{
		CheckPackageManager : Command{
					ID: "Check Pacman Install",
					Header: "pacman", 
					Commands: []string{
					"--version",
					},
		},
		
		InstallPackageManager: Command{
			ID: "Install Pacman",
			Header : "sudo",
			Commands: []string {
				"apt", "install", "pacman-package-manager",
			},
		},

		CheckNetworkUtil: Command{
						ID: "Check nmcli Install",
						Header: "nmcli", 
						Commands: []string{
						"--version",
						},
			},	

		InstallNetworkUtil: Command{
			ID: "Install nmcli",
			Header: "sudo",
			Commands : []string{
				"pacman", "-S", "install", "networkmanager",
			},
		},

		StartNetworkUtil: Command{
					ID: "Start nmcli",
					Header: "sudo",
					Commands : []string{
						"systemctl", "start", "NetworkManager",
					},
				},
		
		RetrieveActiveConnection: Command{
			ID: "retrieve active connection",
			Header: "bash",
			Commands: []string{
				"nmcli -t -f NAME connection show --active",
			},
		},

		// SetNewDNS: Command{
		// 	ID: "Setting new DNS",
		// 	Header: "nmcli",
		// 	Commands: []string{
		// 		"connection",
		// 		"modify",
		// 		"%s",
		// 		"ipv4.dns",
		// 		"%s",
		// 	},
		// },
	}
}
