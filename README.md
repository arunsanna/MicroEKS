# MicroEKS

A lightweight local Kubernetes environment that mimics an AWS EKS cluster using MicroK8s running in a Multipass virtual machine.

## Overview

MicroEKS provides a simple way to deploy a local Kubernetes environment on macOS that resembles an AWS EKS cluster. It's perfect for:

- Local development and testing of Kubernetes applications
- Learning Kubernetes without cloud costs
- Testing EKS-specific configurations locally

## Requirements

- macOS (Intel or Apple Silicon)
- Homebrew (will be used to install Multipass)
- 16GB+ RAM recommended
- 100GB+ free disk space

## Quick Start

1. Clone this repository:
   ```
   git clone https://github.com/yourusername/microeks.git
   cd microeks
   ```

2. Make the script executable:
   ```
   chmod +x micro-eks.sh
   ```

3. Run the script:
   ```
   ./micro-eks.sh
   ```

4. From the interactive menu, select option 1 to create a new environment.

5. Follow the prompts to configure:
   - VM memory (default: 16GB)
   - VM disk space (default: 100GB)
   - Kubernetes version/channel
   - MicroK8s addons to enable

6. Once setup completes, use kubectl with the cluster:
   ```
   export KUBECONFIG=~/.kube/config-microk8s
   kubectl get nodes
   ```

## Interactive Features

The `micro-eks.sh` script provides an interactive menu for managing your MicroEKS environment:

- **Create**: Deploy a new VM with customizable memory, disk space, Kubernetes version, and addons
- **Start**: Start an existing MicroEKS VM
- **Stop**: Stop a running MicroEKS VM
- **Destroy**: Remove the VM and clean up resources
- **Status**: Check the current status of your MicroEKS environment

You can also use command-line arguments for direct access to functions:

```bash
# Create a new environment with interactive prompts
./micro-eks.sh create

# Start the environment
./micro-eks.sh start

# Stop the environment
./micro-eks.sh stop

# Destroy the environment
./micro-eks.sh destroy

# Show status
./micro-eks.sh status
```

## Customization Options

During the creation process, you can customize:

### VM Resources
- Memory allocation (default: 16GB)
- Disk space (default: 100GB)

### Kubernetes Version
Choose from:
- 1.28/stable (default)
- 1.29/stable
- 1.27/stable
- latest/stable
- Custom channel

### MicroK8s Addons
Default addons:
- dns: CoreDNS
- dashboard: The Kubernetes dashboard
- storage: Storage class and default storage pool
- ingress: Ingress controller

Additional available addons:
- metallb: Load balancer for bare metal
- metrics-server: Metrics server for resource metrics
- host-access: Allow pods to reach the host
- registry: Private registry

## Usage Tips

### Accessing the Kubernetes Dashboard

```bash
kubectl port-forward -n kube-system service/kubernetes-dashboard 10443:443
```

Then visit: https://127.0.0.1:10443 in your browser

### Managing the VM

```bash
# Using the interactive menu
./micro-eks.sh

# Or direct commands
./micro-eks.sh start
./micro-eks.sh stop
./micro-eks.sh status
```

### Permanently Configure kubectl

To permanently set the cluster as your default:

```bash
cp ~/.kube/config-microk8s ~/.kube/config
```

### Destroying the Environment

When you're done using MicroEKS, you can completely remove it with:

```bash
./micro-eks.sh destroy
```

This will:
- Delete the Multipass VM
- Remove the kubeconfig file
- Clean up all related resources

## Troubleshooting

### Connection Refused Errors

If you encounter "connection refused" errors when trying to use kubectl, ensure:
1. The VM is running (`./micro-eks.sh status` or `multipass info eks-vm`)
2. Your KUBECONFIG is correctly set
3. The IP address in ~/.kube/config-microk8s matches the VM's IP (`multipass info eks-vm | grep IPv4`)

### Resetting the Environment

To completely reset your environment:

```bash
./micro-eks.sh destroy
./micro-eks.sh create
```

## License

MIT
