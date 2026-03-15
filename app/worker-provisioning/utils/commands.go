package provisoning

import (
	"fmt"
)

type Commands struct {
	CheckPackageManager Command
	InstallPackageManager Command
	CheckDocker Command
	InstallDocker Command
	BuildImage Command
	SpinUpContainer Command
}

type Command struct {
		Header string
		ID string
		Commands []string
}

func CommandsDarwin() *Commands {
	return &Commands{
		CheckPackageManager : Command{
					ID: "Check Brew Install",
					Header: "brew", 
					Commands: []string{
					"--version",
					},
		},
		
		InstallPackageManager: Command{
			ID: "Install brew",
			Header : "/bin/bash",
			Commands: []string {
				"-c", "$(curl -fsSL https://raw.githubusercontent.com/Homebrew/install/HEAD/install.sh)",
			},
		},

		CheckDocker: Command{
						ID: "Check Docker Install",
						Header: "docker", 
						Commands: []string{
						"--version",
						},
			},	

		InstallDocker: Command{
			ID: "Install Docker",
			Header: "brew",
			Commands : []string{
				"install", "docker",
			},
		},
		
		BuildImage: Command{
			ID: "Build Image",
			Header: "docker",
			Commands: []string{
				"build", 
				"-t", 
				"alpine-ray", 
				".",
			},
		},
	}
}

func SetDockerCommand(mem string, disk string) *Commands {
	return &Commands{
		SpinUpContainer: Command{
			ID: "Spin Up Docker",
			Header : "docker",
			Commands: []string{
				"run",
				fmt.Sprintf("--memory=%s", mem),
				fmt.Sprintf("--shm-size=%s", disk),
				"--net=host",
				"-e", fmt.Sprintf("RAY_HEAD_ADDRESS=127.0.0.1:%v", 6379),
				"alpine-ray",
			},
		},
	}
}
