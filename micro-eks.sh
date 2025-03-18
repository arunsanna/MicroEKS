#!/bin/bash

set -e

# Default values
VM_NAME="eks-vm"
DEFAULT_MEM="16G"
DEFAULT_DISK="100G"
DEFAULT_K8S_CHANNEL="1.28/stable"
DEFAULT_ADDONS="dns dashboard storage ingress"

# Function to destroy the VM and clean up
destroy() {
    echo "Destroying MicroEKS environment..."
    
    # Check if the VM exists before attempting to delete
    if multipass info "$VM_NAME" &>/dev/null; then
        echo "Deleting VM '$VM_NAME'..."
        multipass delete "$VM_NAME"
        multipass purge
        echo "VM deleted successfully."
    else
        echo "VM '$VM_NAME' does not exist. Nothing to delete."
    fi
    
    # Clean up kubeconfig
    if [ -f ~/.kube/config-microk8s ]; then
        echo "Removing kubeconfig file..."
        rm ~/.kube/config-microk8s
    fi
    
    echo "MicroEKS environment destroyed successfully."
    exit 0
}

# Function to start the VM
start() {
    # Check if the VM exists
    if ! multipass info "$VM_NAME" &>/dev/null; then
        echo "VM '$VM_NAME' does not exist. Please create it first."
        exit 1
    fi

    # Start the VM if it's not running
    if [ "$(multipass info "$VM_NAME" | grep State | awk '{print $2}')" != "Running" ]; then
        echo "Starting VM '$VM_NAME'..."
        multipass start "$VM_NAME"
        echo "VM '$VM_NAME' started successfully."
    else
        echo "VM '$VM_NAME' is already running."
    fi
}

# Function to stop the VM
stop() {
    # Check if the VM exists
    if ! multipass info "$VM_NAME" &>/dev/null; then
        echo "VM '$VM_NAME' does not exist. Nothing to stop."
        exit 1
    fi

    # Stop the VM if it's running
    if [ "$(multipass info "$VM_NAME" | grep State | awk '{print $2}')" == "Running" ]; then
        echo "Stopping VM '$VM_NAME'..."
        multipass stop "$VM_NAME"
        echo "VM '$VM_NAME' stopped successfully."
    else
        echo "VM '$VM_NAME' is not running."
    fi
}

# Function to create and configure the VM and MicroK8s
create() {
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

    # Prompt for VM configuration
    read -p "Enter memory size for VM (default: $DEFAULT_MEM): " MEM
    MEM=${MEM:-$DEFAULT_MEM}
    
    read -p "Enter disk size for VM (default: $DEFAULT_DISK): " DISK
    DISK=${DISK:-$DEFAULT_DISK}
    
    # Check if the VM already exists
    if multipass info "$VM_NAME" &>/dev/null; then
        echo "VM '$VM_NAME' already exists."
        read -p "Do you want to delete and recreate it? (y/N): " RECREATE
        if [[ "$RECREATE" == "y" || "$RECREATE" == "Y" ]]; then
            multipass delete "$VM_NAME"
            multipass purge
            echo "Creating new VM '$VM_NAME' with $MEM memory and $DISK disk..."
            multipass launch --name "$VM_NAME" --mem "$MEM" --disk "$DISK"
            echo "VM '$VM_NAME' deployed successfully."
        fi
    else
        echo "Deploying VM '$VM_NAME' with $MEM memory and $DISK disk..."
        multipass launch --name "$VM_NAME" --mem "$MEM" --disk "$DISK"
        echo "VM '$VM_NAME' deployed successfully."
    fi

    # Ensure VM is running
    if [ "$(multipass info "$VM_NAME" | grep State | awk '{print $2}')" != "Running" ]; then
        echo "Starting VM '$VM_NAME'..."
        multipass start "$VM_NAME"
        # Wait a bit for VM to fully initialize
        sleep 10
    fi

    # Prompt for MicroK8s channel
    echo ""
    echo "Available MicroK8s channels:"
    echo "1) 1.28/stable (default)"
    echo "2) 1.29/stable"
    echo "3) 1.27/stable"
    echo "4) latest/stable"
    echo "5) Custom channel"
    
    read -p "Select MicroK8s channel [1-5]: " CHANNEL_CHOICE
    
    case $CHANNEL_CHOICE in
        2) K8S_CHANNEL="1.29/stable" ;;
        3) K8S_CHANNEL="1.27/stable" ;;
        4) K8S_CHANNEL="latest/stable" ;;
        5) read -p "Enter custom channel: " K8S_CHANNEL ;;
        *) K8S_CHANNEL="$DEFAULT_K8S_CHANNEL" ;;
    esac
    
    echo "Installing MicroK8s (channel: $K8S_CHANNEL) in VM..."

    # Install MicroK8s inside the VM
    multipass exec "$VM_NAME" -- sudo snap install microk8s --classic --channel="$K8S_CHANNEL"

    # Wait for MicroK8s to be ready
    echo "Waiting for MicroK8s to be ready..."
    multipass exec "$VM_NAME" -- sudo microk8s status --wait-ready

    # Prompt for addons
    echo ""
    echo "MicroK8s Addons:"
    echo "1) Default addons (dns, dashboard, storage, ingress)"
    echo "2) Custom selection"
    
    read -p "Select option [1-2]: " ADDON_CHOICE
    
    if [ "$ADDON_CHOICE" == "2" ]; then
        ADDONS=""
        echo "Available addons:"
        echo "- dns: CoreDNS"
        echo "- dashboard: The Kubernetes dashboard"
        echo "- storage: Storage class and default storage pool"
        echo "- ingress: Ingress controller"
        echo "- metallb: Load balancer for bare metal"
        echo "- metrics-server: Metrics server for resource metrics"
        echo "- host-access: Allow pods to reach the host"
        echo "- registry: Private registry"
        
        read -p "Enter space-separated list of addons to enable: " ADDONS
    else
        ADDONS="$DEFAULT_ADDONS"
    fi
    
    # Enable necessary addons to simulate EKS functionality
    echo "Enabling selected addons: $ADDONS"
    for addon in $ADDONS; do
        echo "Enabling $addon..."
        multipass exec "$VM_NAME" -- sudo microk8s enable "$addon"
    done

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
    echo "MicroEKS has been successfully installed in the VM!"
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
    echo "To destroy this environment, run:"
    echo "./micro-eks.sh destroy"
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
}

# Function to display status
status() {
    if ! multipass info "$VM_NAME" &>/dev/null; then
        echo "VM '$VM_NAME' does not exist."
        exit 1
    fi
    
    echo "VM Status:"
    multipass info "$VM_NAME"
    
    if [ "$(multipass info "$VM_NAME" | grep State | awk '{print $2}')" == "Running" ]; then
        echo ""
        echo "Kubernetes Status:"
        KUBECONFIG=~/.kube/config-microk8s kubectl cluster-info
        echo ""
        KUBECONFIG=~/.kube/config-microk8s kubectl get nodes
    fi
}

# Main menu
if [ "$1" == "destroy" ]; then
    destroy
elif [ "$1" == "start" ]; then
    start
elif [ "$1" == "stop" ]; then
    stop
elif [ "$1" == "create" ]; then
    create
elif [ "$1" == "status" ]; then
    status
else
    echo "MicroEKS - Interactive EKS-like Kubernetes Environment"
    echo ""
    echo "Available commands:"
    echo "1) Create new environment"
    echo "2) Start environment"
    echo "3) Stop environment"
    echo "4) Destroy environment"
    echo "5) Show status"
    echo "6) Exit"
    echo ""
    read -p "Select an option [1-6]: " OPTION
    
    case $OPTION in
        1) create ;;
        2) start ;;
        3) stop ;;
        4) destroy ;;
        5) status ;;
        *) exit 0 ;;
    esac
fi