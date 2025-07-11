package lsp

type InitializeRequest struct {
	Request
	Params InitializeRequestParams `json:"params"`
}

type InitializeRequestParams struct {
	ClientInfo *ClientInfo `json:"clientInfo"`
}

type ClientInfo struct {
	Name    string `json:"name"`
	Version string `json:"version"`
}

type ServerInfo struct {
	Name    string `json:"name"`
	Version string `json:"version"`
}

type InitializeResponse struct {
 	Response
	Result InitializeResult `json:"result"`
}

type InitializeResult struct {
	Capabilities ServerCapabilities `json:"capabilities"`
	ServerInfo   *ServerInfo        `json:"serverInfo"`
}

type ServerCapabilities struct {
	TextDocumentSync   TextDocumentSyncKind `json:"textDocumentSync"`
	HoverProvider      bool                 `json:"hoverProvider"`
	DefinitionProvider bool                 `json:"definitionProvider"`
	CodeActionProvider bool                 `json:"codeActionProvider"`
	CompletionProvider map[string]any       `json:"completionProvider"`
}

func NewInitializeResponse(id int) InitializeResponse {
	return InitializeResponse{
		Response: Response{
			RPC: "2.0",
			ID:  &id,
		},
		Result:   InitializeResult{
			Capabilities: ServerCapabilities{
				TextDocumentSync:   Full,
				HoverProvider:      true,
				DefinitionProvider: true,
				CodeActionProvider: true,
				CompletionProvider: map[string]any{},
			},
			ServerInfo:   &ServerInfo{
				Name:    "progress_ls",
				Version: "0.0.1-beta",
			},
		},
	}
}
