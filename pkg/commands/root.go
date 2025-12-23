package commands

import (
	"fmt"
	"os"

	"github.com/arunsanna/MicroEKS/pkg/platform"
	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "micro-eks",
	Short: "MicroEKS - A lightweight Kubernetes environment using Multipass",
	Long:  `MicroEKS allows you to easily create, manage, and destroy a local Kubernetes cluster using Multipass and MicroK8s.`,
	PersistentPreRun: func(cmd *cobra.Command, args []string) {
		// Skip check for help and completion commands
		if cmd.Name() == "help" || cmd.Name() == "completion" {
			return
		}

		pm := platform.NewManager()
		if !pm.CheckMultipass() {
			fmt.Println("Multipass prerequisite not found. Attempting to install...")
			if err := pm.InstallMultipass(); err != nil {
				fmt.Printf("Error installing multipass: %v\n", err)
				fmt.Println("Please install multipass manually for your operating system.")
				os.Exit(1)
			}
			fmt.Println("Multipass installed successfully.")
		}
	},
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
