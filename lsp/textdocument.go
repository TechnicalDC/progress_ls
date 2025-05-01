package lsp

type TextDocumentSyncKind int

const (
	None TextDocumentSyncKind = iota
	Full
	Incremental
)

type TextDocmentItem struct {
	URI        string `json:"uri"`
	LanguageId string `json:"languageId"`
	Version    int		`json:"version"`
	Text       string `json:"text"`
}

type TextDocumentIndentifier struct {
	URI string `json:"uri"`
}

type VersionedTextDocumentIndentifier struct {
	TextDocumentIndentifier
	Version int `json:"version"`
}

type TextDocumentPositionParams struct {
	TextDocument TextDocumentIndentifier `json:"textDocument"`
	Position 	Position 					`json:"position"`
}

type Position struct {
	Line      int `json:"line"`
	Character int `json:"character"`
}

type Location struct {
	URI string `json:"uri"`
	Range Range `json:"range"`
}

type Range struct {
	Start Position `json:"start"`
	End 	Position	`json:"end"`
}

type WorkspaceEdit struct {
	Changes map[string][]TextEdit `json:"changes"`
}

type TextEdit struct {
	Range   Range  `json:"range"`
	NewText string `json:"newText"`
}
