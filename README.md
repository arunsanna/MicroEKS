# EKS Local

A lightweight local Kubernetes environment that mimics an AWS EKS cluster using MicroK8s running in a Multipass virtual machine.

## Overview

EKS Local provides a simple way to deploy a local Kubernetes environment on macOS that resembles an AWS EKS cluster. It's perfect for:

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
   git clone https://github.com/yourusername/eks-local.git
   cd eks-local
   ```

2. Make the deployment script executable:
   ```
   chmod +x deploy.sh
   ```

3. Run the deployment script:
   ```
   ./deploy.sh
   ```

4. Configure kubectl to use the local cluster:
   ```
   export KUBECONFIG=~/.kube/config-microk8s
   ```

5. Verify the cluster is working:
   ```
   kubectl get nodes
   ```

## Features

- Automatic installation of Multipass VM manager
- Single VM deployment with MicroK8s Kubernetes
- Pre-configured with essential Kubernetes addons:
  - DNS (CoreDNS)
  - Dashboard
  - Storage
  - Ingress
- Automatic kubectl configuration

## Usage Tips

### Accessing the Kubernetes Dashboard

```bash
kubectl port-forward -n kube-system service/kubernetes-dashboard 10443:443
```

Then visit: https://127.0.0.1:10443 in your browser

### Managing the VM

```bash
# Stop the VM
multipass stop eks-vm

# Start the VM
multipass start eks-vm

# Get VM info
multipass info eks-vm

# Shell into the VM
multipass shell eks-vm
```

### Permanently Configure kubectl

To permanently set the cluster as your default:

```bash
cp ~/.kube/config-microk8s ~/.kube/config
```

## Troubleshooting

### Connection Refused Errors

If you encounter "connection refused" errors when trying to use kubectl, ensure:
1. The VM is running (`multipass info eks-vm`)
2. Your KUBECONFIG is correctly set
3. The IP address in ~/.kube/config-microk8s matches the VM's IP (`multipass info eks-vm | grep IPv4`)

### Resetting the Environment

To completely reset your environment:

```bash
multipass delete eks-vm
multipass purge
./deploy.sh
```

## License

MIT
