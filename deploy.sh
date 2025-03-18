#!/bin/bash
# filepath: /Users/megamind/code/eks-local/deploy.sh

set -e

# Check if running on macOS
if [[ "$(uname)" != "Darwin" ]]; then
    echo "This script only supports macOS."
    exit 1
fi

# Check if multipass is already installed
if ! command -v multipass &>/dev/null; then
    # Determine CPU architecture
    ARCH=$(uname -m)
    echo "Detected architecture: $ARCH"

    # For Apple Silicon (arm64) use Homebrew to install multipass
    if [[ "$ARCH" == "arm64" ]]; then
        echo "Installing multipass via Homebrew (Cask) for Apple Silicon..."
        brew install --cask multipass
    else
        echo "Installing multipass via Homebrew (Cask) for Intel..."
        brew install --cask multipass
    fi

    echo "Multipass installation complete."
else
    echo "Multipass is already installed."
fi

# Define VM name
VM_NAME="eks-vm"

# Check if the VM already exists
if multipass info "$VM_NAME" &>/dev/null; then
    echo "VM '$VM_NAME' already exists."
else
    echo "Deploying VM '$VM_NAME' with 16GB memory and 100GB disk..."
    multipass launch --name "$VM_NAME" --mem 16G --disk 100G
    echo "VM '$VM_NAME' deployed successfully."
fi

# Ensure VM is running
if [ "$(multipass info "$VM_NAME" | grep State | awk '{print $2}')" != "Running" ]; then
    echo "Starting VM '$VM_NAME'..."
    multipass start "$VM_NAME"
    # Wait a bit for VM to fully initialize
    sleep 10
fi

echo "Installing MicroK8s (EKS distribution) in VM..."

# Install MicroK8s inside the VM
multipass exec "$VM_NAME" -- sudo snap install microk8s --classic --channel=1.28/stable

# Wait for MicroK8s to be ready
echo "Waiting for MicroK8s to be ready..."
multipass exec "$VM_NAME" -- sudo microk8s status --wait-ready

# Enable necessary addons to simulate EKS functionality
echo "Enabling required addons (DNS, dashboard, storage, ingress)..."
multipass exec "$VM_NAME" -- sudo microk8s enable dns dashboard storage ingress

# Add the current user to microk8s group inside VM
multipass exec "$VM_NAME" -- sudo usermod -a -G microk8s ubuntu

# Get VM's IP address
VM_IP=$(multipass info "$VM_NAME" | grep IPv4 | awk '{print $2}')

# Create .kube directory locally if it doesn't exist
mkdir -p ~/.kube

# Get the kubeconfig from the VM and save it locally
echo "Fetching kubeconfig from the VM..."
multipass exec "$VM_NAME" -- sudo microk8s config > ~/.kube/config-microk8s-temp

# Fix the server URL in the kubeconfig to use the VM's IP address instead of localhost
sed "s/127.0.0.1:16443/$VM_IP:16443/g" ~/.kube/config-microk8s-temp > ~/.kube/config-microk8s
rm ~/.kube/config-microk8s-temp

# Make a copy of the current kubectl config
if [ -f ~/.kube/config ]; then
    echo "Backing up existing kubectl config..."
    cp ~/.kube/config ~/.kube/config.backup.$(date +%Y%m%d%H%M%S)
fi

# Set proper permissions for the kubeconfig
chmod 600 ~/.kube/config-microk8s

echo "======================================================================================"
echo "MicroK8s (EKS-like distribution) has been successfully installed in the VM!"
echo ""
echo "To use kubectl with this cluster, run:"
echo "export KUBECONFIG=~/.kube/config-microk8s"
echo ""
echo "Or to permanently configure it, run:"
echo "cp ~/.kube/config-microk8s ~/.kube/config"
echo ""
echo "To access the Kubernetes dashboard, run:"
echo "kubectl port-forward -n kube-system service/kubernetes-dashboard 10443:443"
echo "Then visit: https://127.0.0.1:10443 in your browser"
echo ""
echo "VM IP address: $VM_IP"
echo "======================================================================================"

# Quick test to verify the connection works
echo ""
echo "Testing connection to the Kubernetes cluster..."
KUBECONFIG=~/.kube/config-microk8s kubectl get nodes

echo ""
echo "If you see your node above, the connection is working properly."
echo "If not, there might be networking issues between your host and the VM."