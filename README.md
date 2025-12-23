# MicroEKS

A lightweight local Kubernetes environment that mimics an AWS EKS cluster using MicroK8s running in a Multipass virtual machine.

## Overview

MicroEKS provides a simple way to deploy a local Kubernetes environment on macOS, Linux, and Windows that resembles an AWS EKS cluster. It's perfect for:

- Local development and testing of Kubernetes applications
- Learning Kubernetes without cloud costs
- Testing EKS-specific configurations locally

## Requirements

- **macOS** (Intel or Apple Silicon), **Linux**, or **Windows**
- **Multipass**: The tool will attempt to install it automatically if missing on supported platforms (macOS/Homebrew, Linux/Snap).

## Build and Install

### Homebrew (macOS/Linux)
You can install MicroEKS using Homebrew:

```bash
brew install arunsanna/tap/micro-eks
```

### Manual Installation
Download the latest binary for your operating system from the [Releases page](https://github.com/arunsanna/MicroEKS/releases).

### Build from Source
This project uses GitHub Actions for releases. To build locally:

```bash
go build -o bin/micro-eks ./cmd/micro-eks
```

## Quick Start

1. **Create** a new environment:
   ```bash
   ./bin/micro-eks create --memory 16G --disk 100G --channel 1.28/stable
   ```

2. **Access** the cluster:
   ```bash
   export KUBECONFIG=~/.kube/config-microk8s
   kubectl get nodes
   ```

## Commands

- **create**: Deploy a new VM with customizable memory, disk space, Kubernetes version, and addons.
- **start**: Start an existing MicroEKS VM.
- **stop**: Stop a running MicroEKS VM.
- **destroy**: Remove the VM and clean up resources.
- **status**: Check the current status of your MicroEKS environment.

## Customization Options

### Create Flags
- `--memory`: VM memory (default: 16G)
- `--disk`: VM disk size (default: 100G)
- `--channel`: MicroK8s channel (default: 1.28/stable)
- `--addons`: Comma-separated list of addons (default: dns,dashboard,storage,ingress)

## Usage Tips

### Accessing the Kubernetes Dashboard

```bash
kubectl port-forward -n kube-system service/kubernetes-dashboard 10443:443
```

Then visit: https://127.0.0.1:10443 in your browser

### Destroying the Environment

```bash
./bin/micro-eks destroy
```

This will:
- Delete the Multipass VM
- Remove the kubeconfig file
- Clean up all related resources

## Troubleshooting

### Connection Refused Errors

If you encounter "connection refused" errors when trying to use kubectl, ensure:
1. The VM is running (`./bin/micro-eks status` or `multipass info eks-vm`)
2. Your KUBECONFIG is correctly set
3. The IP address in `~/.kube/config-microk8s` matches the VM's IP (`multipass info eks-vm | grep IPv4`)

### Resetting the Environment

To completely reset your environment:

```bash
./bin/micro-eks destroy
./bin/micro-eks create
```

## License

MIT
