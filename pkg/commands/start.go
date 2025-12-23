package commands

import (
	"fmt"
	"os"

	"github.com/arunsanna/MicroEKS/pkg/multipass"
	"github.com/spf13/cobra"
)

var startCmd = &cobra.Command{
	Use:   "start",
	Short: "Start the MicroEKS environment",
	Run: func(cmd *cobra.Command, args []string) {
		client := multipass.NewClient()
		vmName := "eks-vm"

		if !client.Exists(vmName) {
			fmt.Printf("VM '%s' does not exist. Please create it first.\n", vmName)
			os.Exit(1)
		}

		if client.IsRunning(vmName) {
			fmt.Printf("VM '%s' is already running.\n", vmName)
			return
		}

		fmt.Printf("Starting VM '%s'...\n", vmName)
		if err := client.Start(vmName); err != nil {
			fmt.Printf("Error starting VM: %v\n", err)
			os.Exit(1)
		}
		fmt.Printf("VM '%s' started successfully.\n", vmName)
	},
}

func init() {
	rootCmd.AddCommand(startCmd)
}
