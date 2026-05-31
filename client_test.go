package web3

import (
	"context"
	"testing"
)

func TestNewClient(t *testing.T) {
	c := NewClient("http://localhost:8545")
	if c == nil {
		t.Fatal("NewClient() returned nil")
	}
}

func TestClientEth(t *testing.T) {
	c := NewClient("http://localhost:8545")
	if c.Eth() == nil {
		t.Fatal("Client.Eth() returned nil")
	}
}

func TestClientCall(t *testing.T) {
	srv := newMockRPCServer("0x64")
	defer srv.Close()

	result, err := NewClient(srv.URL).Call(context.Background(), "eth_blockNumber", []interface{}{})
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

	_, err := NewClient(srv.URL).Call(context.Background(), "unknown_method", []interface{}{})
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
	cancel()

	_, err := NewClient(srv.URL).Call(ctx, "eth_blockNumber", []interface{}{})
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
