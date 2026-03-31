package main

import (
	provisoning "hypha/app/worker-provisioning/provisioning"
)
func main() {
	provisoning.SpinUp(0.2)
}
