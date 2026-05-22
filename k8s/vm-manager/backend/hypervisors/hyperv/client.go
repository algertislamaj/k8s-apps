package hyperv

import (
	"context"
	"encoding/json"
	"fmt"
	"vm-manager/models"

	"github.com/masterzen/winrm"
)

type Client struct {
	client *winrm.Client
	host   string
}

func NewClient(host, username, password string) (*Client, error) {
	endpoint := winrm.NewEndpoint(host, 5985, false, false, nil, nil, nil, 0)
	client, err := winrm.NewClient(endpoint, username, password)
	if err != nil {
		return nil, fmt.Errorf("hyperv connect error: %w", err)
	}
	return &Client{client: client, host: host}, nil
}

func (c *Client) ListVMs(ctx context.Context) ([]models.VM, error) {
	script := `
        Get-VM | Select-Object Name, State, ProcessorCount, MemoryAssigned, VMId |
        ConvertTo-Json -Compress
    `
	stdout, _, _, err := c.client.RunWithContextWithString(ctx, winrm.Powershell(script), "")
	if err != nil {
		return nil, err
	}

	var raw []map[string]interface{}
	if err := json.Unmarshal([]byte(stdout), &raw); err != nil {
		return nil, err
	}

	var vms []models.VM
	for _, v := range raw {
		vms = append(vms, models.VM{
			ID:         fmt.Sprintf("%v", v["VMId"]),
			Name:       fmt.Sprintf("%v", v["Name"]),
			Status:     fmt.Sprintf("%v", v["State"]),
			CPUCount:   int(v["ProcessorCount"].(float64)),
			MemoryMB:   int(v["MemoryAssigned"].(float64)) / 1024 / 1024,
			Hypervisor: "hyperv",
			Host:       c.host,
		})
	}
	return vms, nil
}

func (c *Client) UpdateCPU(ctx context.Context, vmName string, count int) error {
	script := fmt.Sprintf(`Set-VMProcessor -VMName "%s" -Count %d`, vmName, count)
	_, _, _, err := c.client.RunWithContextWithString(ctx, winrm.Powershell(script), "")
	return err
}

func (c *Client) UpdateMemory(ctx context.Context, vmName string, memoryMB int) error {
	script := fmt.Sprintf(`Set-VMMemory -VMName "%s" -StartupBytes %dMB`, vmName, memoryMB)
	_, _, _, err := c.client.RunWithContextWithString(ctx, winrm.Powershell(script), "")
	return err
}
