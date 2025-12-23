package commands

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/arunsanna/MicroEKS/pkg/multipass"
	"github.com/spf13/cobra"
)

var destroyCmd = &cobra.Command{
	Use:   "destroy",
	Short: "Destroy the MicroEKS environment",
	Long:  `Destroys the MicroEKS VM and cleans up related configuration files.`,
	Run: func(cmd *cobra.Command, args []string) {
		client := multipass.NewClient()
		vmName := "eks-vm"

		fmt.Println("Destroying MicroEKS environment...")

		if client.Exists(vmName) {
			fmt.Printf("Deleting VM '%s'...\n", vmName)
			if err := client.Delete(vmName); err != nil {
				fmt.Printf("Error deleting VM: %v\n", err)
			}
			client.Purge()
			fmt.Println("VM deleted successfully.")
		} else {
			fmt.Printf("VM '%s' does not exist. Nothing to delete.\n", vmName)
		}

		// Clean up kubeconfig
		home, _ := os.UserHomeDir()
		configPath := filepath.Join(home, ".kube", "config-microk8s")
		if _, err := os.Stat(configPath); err == nil {
			fmt.Println("Removing kubeconfig file...")
			os.Remove(configPath)
		}

		fmt.Println("MicroEKS environment destroyed successfully.")
	},
}

func init() {
	rootCmd.AddCommand(destroyCmd)
}
