package lsp

type Request struct {
	RPC    string `json:"jsonrpc"`
	ID     int    `json:"id"`
	Method string `json:"method"`

	// Params: will be handled later
}

type Response struct {
	RPC string `json:"jsonrpc"`
	ID  *int   `json:"id,omitempty"`

	// Results: will be handled later
	// Error: will be handled later
}

type Notification struct {
	RPC    string `json:"jsonrpc"`
	Method string `json:"method"`
}
