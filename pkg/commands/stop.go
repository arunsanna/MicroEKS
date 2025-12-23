package commands

import (
	"fmt"
	"os"

	"github.com/arunsanna/MicroEKS/pkg/multipass"
	"github.com/spf13/cobra"
)

var stopCmd = &cobra.Command{
	Use:   "stop",
	Short: "Stop the MicroEKS environment",
	Run: func(cmd *cobra.Command, args []string) {
		client := multipass.NewClient()
		vmName := "eks-vm"

		if !client.Exists(vmName) {
			fmt.Printf("VM '%s' does not exist. Nothing to stop.\n", vmName)
			os.Exit(1)
		}

		if !client.IsRunning(vmName) {
			fmt.Printf("VM '%s' is not running.\n", vmName)
			return
		}

		fmt.Printf("Stopping VM '%s'...\n", vmName)
		if err := client.Stop(vmName); err != nil {
			fmt.Printf("Error stopping VM: %v\n", err)
			os.Exit(1)
		}
		fmt.Printf("VM '%s' stopped successfully.\n", vmName)
	},
}

func init() {
	rootCmd.AddCommand(stopCmd)
}
