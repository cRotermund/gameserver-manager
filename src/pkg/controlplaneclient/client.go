package controlplaneclient

import (
	"context"
	"fmt"
	"log/slog"
	"net"
	"net/http"
	"net/url"
	"strconv"
	"time"

	"resty.dev/v3"
)

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

type ClientError struct {
	StatusCode int
	APIError   *APIError
	Err        error
}

func (e *ClientError) Error() string {
	if e.APIError != nil && e.APIError.Error != "" {
		return fmt.Sprintf("gsm api returned %d: %s", e.StatusCode, e.APIError.Error)
	}
	if e.Err != nil {
		return fmt.Sprintf("gsm api request failed: %v", e.Err)
	}
	return fmt.Sprintf("gsm api returned HTTP status %d", e.StatusCode)
}

func (e *ClientError) Unwrap() error { return e.Err }

type GsmApiClient struct {
	logger *slog.Logger
	rc     *resty.Client
}

func NewClient(logger *slog.Logger, host string, port int, rc *resty.Client) *GsmApiClient {
	if rc == nil {
		rc = resty.New()
	}
	rc.SetBaseURL("http://" + net.JoinHostPort(host, strconv.Itoa(port)))
	rc.SetHeader("Accept", "application/json")
	rc.SetResultError(&APIError{})

	return &GsmApiClient{
		logger: logger,
		rc:     rc,
	}
}

func (c *GsmApiClient) Healthcheck(ctx context.Context) error {
	_, err := c.do(ctx, http.MethodGet, "/health", nil)
	return err
}

func (c *GsmApiClient) ListServers(ctx context.Context) ([]ServerSummary, error) {
	var result []ServerSummary
	_, err := c.do(ctx, http.MethodGet, "/v1/servers", &result)
	return result, err
}

func (c *GsmApiClient) GetServer(ctx context.Context, serverID string) (*ServerDetail, error) {
	result := new(ServerDetail)
	_, err := c.do(ctx, http.MethodGet, serverPath(serverID), result)
	return result, err
}

func (c *GsmApiClient) ListProcesses(ctx context.Context, serverID string) ([]ProcessInfo, error) {
	var result []ProcessInfo
	_, err := c.do(ctx, http.MethodGet, serverPath(serverID)+"/processes", &result)
	return result, err
}

func (c *GsmApiClient) StartServer(ctx context.Context, serverID string) (*Operation, error) {
	return c.controlServer(ctx, serverID, "start")
}

func (c *GsmApiClient) StopServer(ctx context.Context, serverID string) (*Operation, error) {
	return c.controlServer(ctx, serverID, "stop")
}

func (c *GsmApiClient) RebootServer(ctx context.Context, serverID string) (*Operation, error) {
	return c.controlServer(ctx, serverID, "reboot")
}

func (c *GsmApiClient) GetOperation(ctx context.Context, operationID string) (*Operation, error) {
	result := new(Operation)
	_, err := c.do(ctx, http.MethodGet, "/v1/operations/"+url.PathEscape(operationID), result)
	return result, err
}

func (c *GsmApiClient) controlServer(ctx context.Context, serverID, action string) (*Operation, error) {
	result := new(Operation)
	_, err := c.do(ctx, http.MethodPost, serverPath(serverID)+"/"+action, result)
	return result, err
}

func (c *GsmApiClient) do(ctx context.Context, method, path string, result any) (*resty.Response, error) {
	request := c.rc.R().SetContext(ctx)
	if result != nil {
		request.SetResult(result)
	}

	response, err := request.Execute(method, path)
	if err != nil {
		return response, &ClientError{Err: err}
	}
	if response.IsStatusSuccess() {
		return response, nil
	}

	apiErr := new(APIError)
	if response.ResultError() != nil {
		if decoded, ok := response.ResultError().(*APIError); ok {
			apiErr = decoded
		}
	}
	if apiErr.Error == "" {
		apiErr = nil
	}

	return response, &ClientError{StatusCode: response.StatusCode(), APIError: apiErr}
}

func serverPath(serverID string) string {
	return "/v1/servers/" + url.PathEscape(serverID)
}

func (c *GsmApiClient) Close() {
	if c.rc != nil {
		c.rc.Close()
	}
}

var _ error = (*ClientError)(nil)
