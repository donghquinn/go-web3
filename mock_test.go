package web3

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
)

// newMockRPCServer returns a test server that responds to all RPC calls with result.
func newMockRPCServer(result interface{}) *httptest.Server {
	return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var req RPCRequest
		json.NewDecoder(r.Body).Decode(&req)
		respond(w, req.ID, result)
	}))
}

// newMockRPCErrorServer returns a test server that always responds with an RPC error.
func newMockRPCErrorServer(code int, message string) *httptest.Server {
	return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var req RPCRequest
		json.NewDecoder(r.Body).Decode(&req)
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(RPCResponse{
			ID:    req.ID,
			Error: &RPCError{Code: code, Message: message},
		})
	}))
}

// newDispatchServer returns a test server that routes each RPC method to a pre-set result.
// Methods not present in handlers fall back to "0x0".
func newDispatchServer(handlers map[string]interface{}) *httptest.Server {
	return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var req RPCRequest
		json.NewDecoder(r.Body).Decode(&req)
		result, ok := handlers[req.Method]
		if !ok {
			result = "0x0"
		}
		respond(w, req.ID, result)
	}))
}

func respond(w http.ResponseWriter, id uint64, result interface{}) {
	resultBytes, _ := json.Marshal(result)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(RPCResponse{
		ID:     id,
		Result: json.RawMessage(resultBytes),
	})
}
