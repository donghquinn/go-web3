package web3

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

func newMockRPCServer(result interface{}) *httptest.Server {
	return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var req RPCRequest
		json.NewDecoder(r.Body).Decode(&req)

		resultBytes, _ := json.Marshal(result)
		resp := RPCResponse{
			ID:     req.ID,
			Result: json.RawMessage(resultBytes),
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(resp)
	}))
}

func newMockRPCErrorServer(code int, message string) *httptest.Server {
	return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var req RPCRequest
		json.NewDecoder(r.Body).Decode(&req)

		resp := RPCResponse{
			ID:    req.ID,
			Error: &RPCError{Code: code, Message: message},
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(resp)
	}))
}

func TestNewClient(t *testing.T) {
	c := NewClient("http://localhost:8545")
	if c == nil {
		t.Fatal("NewClient() returned nil")
	}
}

func TestClientEth(t *testing.T) {
	c := NewClient("http://localhost:8545")
	eth := c.Eth()
	if eth == nil {
		t.Fatal("Client.Eth() returned nil")
	}
}

func TestClientCall(t *testing.T) {
	srv := newMockRPCServer("0x64")
	defer srv.Close()

	c := NewClient(srv.URL)
	result, err := c.Call(context.Background(), "eth_blockNumber", []interface{}{})
	if err != nil {
		t.Fatalf("Call() error = %v", err)
	}
	if result == nil {
		t.Fatal("Call() result is nil")
	}
}

func TestClientCallRPCError(t *testing.T) {
	srv := newMockRPCErrorServer(-32601, "method not found")
	defer srv.Close()

	c := NewClient(srv.URL)
	_, err := c.Call(context.Background(), "unknown_method", []interface{}{})
	if err == nil {
		t.Fatal("expected RPC error")
	}
	if err.Error() != "RPC error -32601: method not found" {
		t.Errorf("error message = %q", err.Error())
	}
}

func TestClientCallContextCanceled(t *testing.T) {
	srv := newMockRPCServer("0x1")
	defer srv.Close()

	ctx, cancel := context.WithCancel(context.Background())
	cancel() // already canceled

	c := NewClient(srv.URL)
	_, err := c.Call(ctx, "eth_blockNumber", []interface{}{})
	if err == nil {
		t.Fatal("expected error for canceled context")
	}
}

func TestRPCErrorMessage(t *testing.T) {
	e := &RPCError{Code: -32000, Message: "execution reverted"}
	want := "RPC error -32000: execution reverted"
	if e.Error() != want {
		t.Errorf("RPCError.Error() = %q, want %q", e.Error(), want)
	}
}
