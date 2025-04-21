package lsp

type PublishDiagosticsNotification struct {
	Notification
	Params PublishDiagosticsParams `json:"params"`
}

type PublishDiagosticsParams struct {
	URI         string        `json:"uri"`
	Diagnostics []Diagnostics `json:"diagnostics"`
}

type Diagnostics struct {
	Range    Range  `json:"range"`
	Severity int    `json:"severity"`
	Source   string `json:"source"`
	Message  string `json:"message"`
}
