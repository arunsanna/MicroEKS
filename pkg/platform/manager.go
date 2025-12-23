package platform

import (
	"fmt"
	"os/exec"
	"runtime"
)

type Manager interface {
	CheckMultipass() bool
	InstallMultipass() error
}

func NewManager() Manager {
	switch runtime.GOOS {
	case "darwin":
		return &DarwinManager{}
	case "linux":
		return &LinuxManager{}
	default:
		return &GenericManager{}
	}
}

type DarwinManager struct{}

func (m *DarwinManager) CheckMultipass() bool {
	_, err := exec.LookPath("multipass")
	return err == nil
}

func (m *DarwinManager) InstallMultipass() error {
	fmt.Println("Installing multipass using Homebrew...")
	cmd := exec.Command("brew", "install", "--cask", "multipass")
	cmd.Stdout = exec.Command("cat").Stdout
	cmd.Stderr = exec.Command("cat").Stderr
	return cmd.Run()
}

type LinuxManager struct{}

func (m *LinuxManager) CheckMultipass() bool {
	_, err := exec.LookPath("multipass")
	return err == nil
}

func (m *LinuxManager) InstallMultipass() error {
	fmt.Println("Installing multipass using Snap...")
	cmd := exec.Command("sudo", "snap", "install", "multipass")
	cmd.Stdout = exec.Command("cat").Stdout
	cmd.Stderr = exec.Command("cat").Stderr
	return cmd.Run()
}

type GenericManager struct{}

func (m *GenericManager) CheckMultipass() bool {
	_, err := exec.LookPath("multipass")
	return err == nil
}

func (m *GenericManager) InstallMultipass() error {
	return fmt.Errorf("automatic installation not supported for OS: %s", runtime.GOOS)
}
