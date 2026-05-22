package models

type Hypervisor struct {
    ID       string `json:"id"`
    Name     string `json:"name"`
    Type     string `json:"type"` // "esxi" or "hyperv"
    Host     string `json:"host"`
    Username string `json:"username"`
    Password string `json:"password,omitempty"`
}

type VM struct {
    ID         string `json:"id"`
    Name       string `json:"name"`
    Status     string `json:"status"`
    CPUCount   int    `json:"cpu_count"`
    MemoryMB   int    `json:"memory_mb"`
    DiskGB     int    `json:"disk_gb"`
    Hypervisor string `json:"hypervisor"`
    Host       string `json:"host"`
}

type UpdateResourceRequest struct {
    Value int `json:"value"`
}