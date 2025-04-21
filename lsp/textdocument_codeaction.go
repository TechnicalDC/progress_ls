package lsp

type CodeActionRequest struct {
	Request
	Params TextDocumentCodeActionParams `json:"params"`
}

type TextDocumentCodeActionParams struct {
	TextDocument TextDocumentIndentifier `json:"textDocument"`
	Range 		 Range						 `json:"range"`
	Context 		 CodeActionContext		 `json:"context"`
}

type CodeActionContext struct {}

type CodeActionResponse struct {
	Response
	Result 	[]CodeAction `json:"result"`
}

type CodeAction struct {
	Title   string         `json:"title"`
	Edit    *WorkspaceEdit `json:"edit,omitempty"`
	Command *Command       `json:"command,omitempty"`
}

type Command struct {
	Title     string        `json:"title"`
	Command   string        `json:"command"`
	Arguments []interface{} `json:"arguments"`
}
