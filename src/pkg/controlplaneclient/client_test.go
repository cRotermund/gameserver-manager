package controlplaneclient

import (
	"context"
	"encoding/json"
	"errors"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"resty.dev/v3"
)

func TestGsmApiClientListServers(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet || r.URL.Path != "/v1/servers" {
			t.Fatalf("unexpected request: %s %s", r.Method, r.URL.Path)
		}
		if got := r.Header.Get("Accept"); got != "application/json" {
			t.Fatalf("Accept header = %q", got)
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode([]ServerSummary{{ServerID: "server-1", Name: "Survival", Status: ServerStatusRunning}})
	}))
	defer server.Close()

	client := newTestClient(server)
	servers, err := client.ListServers(context.Background())
	if err != nil {
		t.Fatalf("ListServers() error = %v", err)
	}
	if len(servers) != 1 || servers[0].ServerID != "server-1" {
		t.Fatalf("ListServers() = %#v", servers)
	}
}

func TestGsmApiClientControlServer(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost || r.URL.Path != "/v1/servers/server-1/start" {
			t.Fatalf("unexpected request: %s %s", r.Method, r.URL.Path)
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(Operation{
			OperationID: "operation-1",
			Type:        OperationTypeStart,
			ServerID:    "server-1",
			Status:      OperationStatusPending,
			CreatedAt:   time.Now().UTC(),
		})
	}))
	defer server.Close()

	client := newTestClient(server)
	operation, err := client.StartServer(context.Background(), "server-1")
	if err != nil {
		t.Fatalf("StartServer() error = %v", err)
	}
	if operation.OperationID != "operation-1" || operation.Status != OperationStatusPending {
		t.Fatalf("StartServer() = %#v", operation)
	}
}

func TestGsmApiClientReturnsStructuredAPIError(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusConflict)
		json.NewEncoder(w).Encode(APIError{Error: "Server is already running", Code: "STATE_CONFLICT"})
	}))
	defer server.Close()

	client := newTestClient(server)
	_, err := client.StartServer(context.Background(), "server-1")
	if err == nil {
		t.Fatal("StartServer() error = nil")
	}

	var clientErr *ClientError
	if !errors.As(err, &clientErr) {
		t.Fatalf("error type = %T, want *ClientError", err)
	}
	if clientErr.StatusCode != http.StatusConflict || clientErr.APIError.Code != "STATE_CONFLICT" {
		t.Fatalf("client error = %#v", clientErr)
	}
}

func newTestClient(server *httptest.Server) *GsmApiClient {
	client := NewClient(slog.Default(), "placeholder", 80, resty.New())
	client.rc.SetBaseURL(server.URL)
	return client
}
