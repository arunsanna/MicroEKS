package commands

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/arunsanna/MicroEKS/pkg/multipass"
	"github.com/spf13/cobra"
)

var createCmd = &cobra.Command{
	Use:   "create",
	Short: "Create a new MicroEKS environment",
	Long:  `Creates and configures a new MicroEKS environment using Multipass and MicroK8s.`,
	Run: func(cmd *cobra.Command, args []string) {
		mem, _ := cmd.Flags().GetString("memory")
		disk, _ := cmd.Flags().GetString("disk")
		channel, _ := cmd.Flags().GetString("channel")
		addons, _ := cmd.Flags().GetString("addons")

		fmt.Println("Initializing MicroEKS setup...")

		client := multipass.NewClient()
		vmName := "eks-vm"

		// Check if VM exists
		if client.Exists(vmName) {
			fmt.Printf("VM '%s' already exists.\n", vmName)
			// Simple prompt for recreation (optional, could be a flag to force)
			// For now, let's assume if it exists we might want to stop/start or just warn.
			// The bash script asked to delete.
			fmt.Println("Do you want to delete and recreate it? (y/N): ")
			var response string
			fmt.Scanln(&response)
			if response == "y" || response == "Y" {
				fmt.Println("Deleting existing VM...")
				client.Delete(vmName)
				client.Purge()
			} else {
				fmt.Println("Aborting creation.")
				return
			}
		}

		fmt.Printf("Deploying VM '%s' with %s memory and %s disk...\n", vmName, mem, disk)
		if err := client.Launch(vmName, mem, disk); err != nil {
			fmt.Printf("Error launching VM: %v\n", err)
			os.Exit(1)
		}
		fmt.Printf("VM '%s' deployed successfully.\n", vmName)

		// Wait for running state? content.Launch implementation waits?
		// Usually launch blocks until ready or fails.

		fmt.Printf("Installing MicroK8s (channel: %s)...\n", channel)
		installCmd := fmt.Sprintf("sudo snap install microk8s --classic --channel=%s", channel)
		if _, err := client.Exec(vmName, installCmd); err != nil {
			fmt.Printf("Error installing MicroK8s: %v\n", err)
			os.Exit(1)
		}

		fmt.Println("Waiting for MicroK8s to be ready...")
		client.Exec(vmName, "sudo microk8s status --wait-ready")

		// Addons
		addonList := strings.Split(addons, ",")
		for _, addon := range addonList {
			addon = strings.TrimSpace(addon)
			if addon == "" {
				continue
			}
			fmt.Printf("Enabling addon: %s...\n", addon)
			client.Exec(vmName, "sudo microk8s enable "+addon)
		}

		// Group permission
		client.Exec(vmName, "sudo usermod -a -G microk8s ubuntu")

		// Kubeconfig setup
		vmIP, err := client.GetIP(vmName)
		if err != nil {
			fmt.Printf("Error getting VM IP: %v\n", err)
			// Continue anyway, maybe just kubeconfig will fail
		}

		fmt.Println("Fetching kubeconfig...")
		config, err := client.Exec(vmName, "sudo microk8s config")
		if err != nil {
			fmt.Printf("Error fetching kubeconfig: %v\n", err)
		} else {
			// Replace 127.0.0.1 with VM IP
			if vmIP != "" {
				config = strings.ReplaceAll(config, "127.0.0.1", vmIP)
			}

			// Save to ~/.kube/config-microk8s
			home, _ := os.UserHomeDir()
			kubeDir := filepath.Join(home, ".kube")
			os.MkdirAll(kubeDir, 0755)
			configPath := filepath.Join(kubeDir, "config-microk8s")
			if err := os.WriteFile(configPath, []byte(config), 0600); err != nil {
				fmt.Printf("Error saving kubeconfig: %v\n", err)
			} else {
				fmt.Printf("Kubeconfig saved to %s\n", configPath)
				fmt.Println("You can use it via: export KUBECONFIG=" + configPath)
			}
		}

		fmt.Println("MicroEKS setup complete!")
	},
}

func init() {
	rootCmd.AddCommand(createCmd)
	createCmd.Flags().String("memory", "16G", "Memory size for the VM")
	createCmd.Flags().String("disk", "100G", "Disk size for the VM")
	createCmd.Flags().String("channel", "1.28/stable", "MicroK8s channel to install")
	createCmd.Flags().String("addons", "dns,dashboard,storage,ingress", "Comma separated list of addons")
}
