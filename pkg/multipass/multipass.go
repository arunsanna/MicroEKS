package multipass

import (
	"bytes"
	"fmt"
	"os/exec"
	"strings"
)

// Runner interface for executing commands
type Runner interface {
	Launch(name, memory, disk string) error
	Info(name string) (string, error)
	Start(name string) error
	Stop(name string) error
	Delete(name string) error
	Purge() error
	Exec(name, command string) (string, error)
	Exists(name string) bool
	IsRunning(name string) bool
	GetIP(name string) (string, error)
}

// Client implements Runner using the multipass CLI
type Client struct{}

func NewClient() *Client {
	return &Client{}
}

func (c *Client) run(args ...string) (string, error) {
	cmd := exec.Command("multipass", args...)
	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr
	err := cmd.Run()
	if err != nil {
		return "", fmt.Errorf("multipass command failed: %s, stderr: %s", err, stderr.String())
	}
	return strings.TrimSpace(stdout.String()), nil
}

func (c *Client) Launch(name, memory, disk string) error {
	_, err := c.run("launch", "--name", name, "--mem", memory, "--disk", disk)
	return err
}

func (c *Client) Info(name string) (string, error) {
	return c.run("info", name)
}

func (c *Client) Start(name string) error {
	_, err := c.run("start", name)
	return err
}

func (c *Client) Stop(name string) error {
	_, err := c.run("stop", name)
	return err
}

func (c *Client) Delete(name string) error {
	_, err := c.run("delete", name)
	return err
}

func (c *Client) Purge() error {
	_, err := c.run("purge")
	return err
}

func (c *Client) Exec(name, command string) (string, error) {
	// exec command structure: multipass exec <name> -- <command>
	// We need to split the command string if it contains arguments, but exec.Command expects separate args.
	// However, multipass exec expects the command to run inside the VM.
	// simpler to pass the entire command string as arguments to bash -c inside?
	// "multipass exec name -- command"
	// Actually, exec.Command("multipass", "exec", name, "--", "bash", "-c", command) might be safer for complex commands.
	// Let's try direct execution first as per the bash script.
	// Bash script: multipass exec "$VM_NAME" -- sudo snap install ...

	// For specific complex commands, we might need flexibility.
	// Let's use a shell execution inside the VM to handle pipes/redirects easily.

	cmd := exec.Command("multipass", "exec", name, "--", "bash", "-c", command)
	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr
	err := cmd.Run()
	if err != nil {
		return "", fmt.Errorf("exec in VM failed: %v, stderr: %s", err, stderr.String())
	}
	return strings.TrimSpace(stdout.String()), nil
}

func (c *Client) Exists(name string) bool {
	// multipass info name returns error if not found
	_, err := c.Info(name)
	return err == nil
}

func (c *Client) IsRunning(name string) bool {
	output, err := c.Info(name)
	if err != nil {
		return false
	}
	return strings.Contains(output, "State:          Running")
}

func (c *Client) GetIP(name string) (string, error) {
	output, err := c.Info(name)
	if err != nil {
		return "", err
	}
	return parseIP(output)
}

func parseIP(output string) (string, error) {
	lines := strings.Split(output, "\n")
	for _, line := range lines {
		if strings.Contains(line, "IPv4:") {
			parts := strings.Fields(line)
			if len(parts) >= 2 {
				return parts[1], nil
			}
		}
	}
	return "", fmt.Errorf("IP not found in info output")
}
