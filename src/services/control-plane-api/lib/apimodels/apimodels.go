package apimodels

import "time"

type ServerStatus string

const (
	ServerStatusStopped  ServerStatus = "stopped"
	ServerStatusStopping ServerStatus = "stopping"
	ServerStatusStarting ServerStatus = "starting"
	ServerStatusRunning  ServerStatus = "running"
)

type HardwareSpecs struct {
	Cores         int      `json:"cores"`
	MemoryGB      int      `json:"memory_gb"`
	BandwidthGbps *float64 `json:"bandwidth_gbps,omitempty"`
}

type ServerSummary struct {
	ServerID string        `json:"server_id"`
	Name     string        `json:"name"`
	Game     string        `json:"game"`
	Specs    HardwareSpecs `json:"specs"`
	Status   ServerStatus  `json:"status"`
}

type ResourceUsage struct {
	CPUPercent    *float64 `json:"cpu_percent,omitempty"`
	MemoryPercent *float64 `json:"memory_percent,omitempty"`
	DiskUsedGB    *float64 `json:"disk_used_gb,omitempty"`
	NetworkRxMbps *float64 `json:"network_rx_mbps,omitempty"`
	NetworkTxMbps *float64 `json:"network_tx_mbps,omitempty"`
}

type ServerDetail struct {
	ServerSummary
	IPAddress        *string        `json:"ip_address,omitempty"`
	ConnectedClients *int           `json:"connected_clients,omitempty"`
	ResourceUsage    *ResourceUsage `json:"resource_usage,omitempty"`
}

type ProcessInfo struct {
	Name       string   `json:"name"`
	PID        int      `json:"pid"`
	CPUPercent *float64 `json:"cpu_percent,omitempty"`
	MemoryMB   *float64 `json:"memory_mb,omitempty"`
}

type OperationType string

const (
	OperationTypeStart  OperationType = "start"
	OperationTypeStop   OperationType = "stop"
	OperationTypeReboot OperationType = "reboot"
)

type OperationStatus string

const (
	OperationStatusPending    OperationStatus = "pending"
	OperationStatusInProgress OperationStatus = "in_progress"
	OperationStatusCompleted  OperationStatus = "completed"
	OperationStatusFailed     OperationStatus = "failed"
)

type Operation struct {
	OperationID string          `json:"operation_id"`
	Type        OperationType   `json:"type"`
	ServerID    string          `json:"server_id"`
	Status      OperationStatus `json:"status"`
	CreatedAt   time.Time       `json:"created_at"`
	CompletedAt *time.Time      `json:"completed_at,omitempty"`
	Error       *string         `json:"error,omitempty"`
}

type APIError struct {
	Error string `json:"error"`
	Code  string `json:"code"`
}
