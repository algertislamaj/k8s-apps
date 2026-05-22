package esxi

import (
    "context"
    "fmt"
    "net/url"
    "vm-manager/models"

    "github.com/vmware/govmomi"
    "github.com/vmware/govmomi/find"
    "github.com/vmware/govmomi/object"
    "github.com/vmware/govmomi/vim25/types"
)

type Client struct {
    client *govmomi.Client
    host   string
}

func NewClient(host, username, password string) (*Client, error) {
    ctx := context.Background()
    u := &url.URL{
        Scheme: "https",
        Host:   host,
        Path:   "/sdk",
        User:   url.UserPassword(username, password),
    }
    client, err := govmomi.NewClient(ctx, u, true) // true = skip SSL verify
    if err != nil {
        return nil, fmt.Errorf("esxi connect error: %w", err)
    }
    return &Client{client: client, host: host}, nil
}

func (c *Client) ListVMs(ctx context.Context) ([]models.VM, error) {
    finder := find.NewFinder(c.client.Client, true)
    vms, err := finder.VirtualMachineList(ctx, "*")
    if err != nil {
        return nil, err
    }

    var result []models.VM
    for _, vm := range vms {
        var moVM mo.VirtualMachine
        err := vm.Properties(ctx, vm.Reference(), []string{"summary", "config"}, &moVM)
        if err != nil {
            continue
        }
        result = append(result, models.VM{
            ID:         vm.Reference().Value,
            Name:       moVM.Summary.Config.Name,
            Status:     string(moVM.Summary.Runtime.PowerState),
            CPUCount:   int(moVM.Summary.Config.NumCpu),
            MemoryMB:   int(moVM.Summary.Config.MemorySizeMB),
            Hypervisor: "esxi",
            Host:       c.host,
        })
    }
    return result, nil
}

func (c *Client) UpdateCPU(ctx context.Context, vmID string, count int) error {
    vm := object.NewVirtualMachine(c.client.Client, types.ManagedObjectReference{
        Type: "VirtualMachine", Value: vmID,
    })
    spec := types.VirtualMachineConfigSpec{NumCPUs: int32(count)}
    task, err := vm.Reconfigure(ctx, spec)
    if err != nil {
        return err
    }
    return task.Wait(ctx)
}

func (c *Client) UpdateMemory(ctx context.Context, vmID string, memoryMB int) error {
    vm := object.NewVirtualMachine(c.client.Client, types.ManagedObjectReference{
        Type: "VirtualMachine", Value: vmID,
    })
    spec := types.VirtualMachineConfigSpec{MemoryMB: int64(memoryMB)}
    task, err := vm.Reconfigure(ctx, spec)
    if err != nil {
        return err
    }
    return task.Wait(ctx)
}