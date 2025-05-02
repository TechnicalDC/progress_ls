package analysis

import (
	"fmt"
	"log"
	"lsp-go/lsp"
	"strings"
)

type State struct {
	Documents map[string]string
}

func NewState() State {
	return State{
		Documents: map[string]string{},
	}
}

func getDiagnostics(logger *log.Logger, row int, text string) []lsp.Diagnostics {
	diagostics := []lsp.Diagnostics{}
	for _, line := range strings.Split(text, ". ") {
		if strings.Contains(line, "define variable") {
			if !strings.Contains(line, "no-undo") {
				diagostics = append(diagostics, createDiagnostics(row, len(line), "no-undo is missing"))
			}
		}
	}
	return diagostics
}

func ProcessDocument(logger *log.Logger, text string) []lsp.Diagnostics {
	var new_content string
	var end_char string
	diagostics := []lsp.Diagnostics{}
	for row, line := range strings.Split(text, "\n") {
		if line != "" {
			end_char = line[len(line) - 1:]
			if end_char == "." || end_char == ":" {
				new_content = new_content + " " + line + "\n"
				diagostics = append(diagostics, getDiagnostics(logger, row, new_content)...)
				new_content = ""
			} else if end_char == " " {
				new_content = new_content + line
			} else{
				new_content = new_content + " " + line
			}
		}
	}
	return diagostics
}

func (s *State) OpenDocument(logger *log.Logger, uri, text string) []lsp.Diagnostics {
	s.Documents[uri] = text
	return ProcessDocument(logger, s.Documents[uri])
}

func (s *State) UpdateDocument(logger *log.Logger, uri string, change lsp.TextDocumentContentChangeEvent) []lsp.Diagnostics {
	s.Documents[uri] = change.Text
	return ProcessDocument(logger, s.Documents[uri])
}

func (s *State) Hover(id int, uri string, position lsp.Position) lsp.HoverResponse {
	document := s.Documents[uri]

	return lsp.HoverResponse{
		Response: lsp.Response{
			RPC: "2.0",
			ID:  &id,
		},
		Result:   lsp.HoverResult{
			Contents: fmt.Sprintf("File: %s, Length: %d", uri, len(document)),
		},
	}
}

func (s *State) Definition(id int, uri string, position lsp.Position) lsp.DefinitionResponse {
	return lsp.DefinitionResponse{
		Response: lsp.Response{
			RPC: "2.0",
			ID:  &id,
		},
		Result:   lsp.Location{
			URI:   uri,
			Range: lsp.Range{
				Start: lsp.Position{
					Line:      position.Line - 1,
					Character: 0,
				},
				End:   lsp.Position{
					Line:      position.Line - 1,
					Character: 0,
				},
			},
		},
	}
}

func createDiagnostics(row, line int, message string) lsp.Diagnostics {
	return lsp.Diagnostics{
		Range:    LineRange(row, row, 0, line),
		Severity: 1,
		Source:   "progress_ls",
		Message:  message,
	}
}

func LineRange(sline, eline, start, end int) lsp.Range {
	return lsp.Range{
		Start: lsp.Position{
			Line:      sline,
			Character: start,
		},
		End:   lsp.Position{
			Line:      eline,
			Character: end,
		},
	}
}

func (s *State) CodeAction(id int, uri string) lsp.CodeActionResponse {
	text := s.Documents[uri]
	actions := []lsp.CodeAction{}

	for row, line := range strings.Split(text, "\n") {
		idx := strings.Index(line, "VS Code")
		if idx >= 0 {
			replaceChange := map[string][]lsp.TextEdit{}
			replaceChange[uri] = []lsp.TextEdit{
				{
					Range: LineRange(row, row, idx, idx + len("VS Code")),
					NewText: "Neovim",
				},
			}

			actions = append(actions, lsp.CodeAction{
				Title: "Replace VS Code with a superior editor",
				Edit: &lsp.WorkspaceEdit{Changes: replaceChange},
			})

			censorChanges := map[string][]lsp.TextEdit{}
			censorChanges[uri] = []lsp.TextEdit{
				{
					Range:   LineRange(row, row, idx, idx + len("VS Code")),
					NewText: "VS C*de",
				},
			}

			actions = append(actions, lsp.CodeAction{
				Title: "Censor to VS C*de",
				Edit: &lsp.WorkspaceEdit{Changes: censorChanges},
			})
		}
	}

	return lsp.CodeActionResponse{
		Response: lsp.Response{
			RPC: "2.0",
			ID:  &id,
		},
		Result:   actions,
	}
}

func (s *State) Completion(id int, uri string) lsp.CompletionResponse {
	items := []lsp.CompletionItem{
		{
			Label:         "Neovim (BTW)",
			Detail:        "I use neovim btw",
			Documentation: "Completion for the most superior text editor",
			Kind: 			lsp.Keyword,
		},
		{
			Label:         "Dilip Chauhan",
			Detail:        "I use neovim btw",
			Documentation: "Completion for the most superior text editor",
			Kind: 			lsp.Class,
		},
	}

	return lsp.CompletionResponse{
		Response: lsp.Response{
			RPC: "2.0",
			ID:  &id,
		},
		Result:   items,
	}
}
