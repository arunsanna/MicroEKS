package commands

import (
	"fmt"
	"os"

	"github.com/arunsanna/MicroEKS/pkg/multipass"
	"github.com/spf13/cobra"
)

var statusCmd = &cobra.Command{
	Use:   "status",
	Short: "Show the status of the MicroEKS environment",
	Run: func(cmd *cobra.Command, args []string) {
		client := multipass.NewClient()
		vmName := "eks-vm"

		if !client.Exists(vmName) {
			fmt.Printf("VM '%s' does not exist.\n", vmName)
			os.Exit(1)
		}

		fmt.Println("VM Status:")
		info, err := client.Info(vmName)
		if err != nil {
			fmt.Printf("Error getting VM info: %v\n", err)
		} else {
			fmt.Println(info)
		}

		if client.IsRunning(vmName) {
			fmt.Println("\nKubernetes Status:")
			// We can run kubectl inside the VM or rely on local kubectl if configured
			// The original script ran: KUBECONFIG=~/.kube/config-microk8s kubectl cluster-info
			// But that requires kubectl on the host. Let's try inside VM for portability?
			// Original script: KUBECONFIG=... kubectl ...
			// If we want this to be standalone, maybe check inside the VM:
			// "microk8s kubectl cluster-info"

			out, err := client.Exec(vmName, "sudo microk8s kubectl cluster-info")
			if err != nil {
				fmt.Printf("Error getting k8s status: %v\n", err)
			} else {
				fmt.Println(out)
			}

			fmt.Println("")
			outNodes, err := client.Exec(vmName, "sudo microk8s kubectl get nodes")
			if err != nil {
				fmt.Printf("Error getting nodes: %v\n", err)
			} else {
				fmt.Println(outNodes)
			}
		}
	},
}

func init() {
	rootCmd.AddCommand(statusCmd)
}
