package analysis

import (
	"fmt"
	"log"
	"lsp-go/lsp"
	"lsp-go/progress"
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

func getDiagnostics(row int, text string) []lsp.Diagnostics {
	diagostics := []lsp.Diagnostics{}

	if len(strings.Split(text, ". ")) > 1 {
		diagostics = append(diagostics, createDiagnostics(row, row, len(text), "single line cannot contains multiple statements"))
	}

	if strings.Contains(text, "define variable") {
		if !strings.Contains(text, "no-undo") {
			diagostics = append(diagostics, createDiagnostics(row, row, len(text), "no-undo is missing"))
		}
	}
	return diagostics
}

func (s *State) getKeyword(uri string, postion lsp.Position) string {
	text := strings.Split(s.Documents[uri],"\n")[postion.Line]
	text = text[postion.Character:]
	return text
}

func ProcessDocument(logger *log.Logger, text string, builtin progress.ProgressKeywords) []lsp.Diagnostics {
	var new_line   string
	var end_char   string
	var classes    []string
	var methods    []string
	var properties []string
	diagostics := []lsp.Diagnostics{}

	text = strings.Trim(text, " ")

	// Process document for listing classes, methods, properties
	for _, line := range strings.Split(text, "\n") {
		if strings.HasPrefix(line, "using") {
			class := line[len("using "):]
			classes = append(classes, class[:len(class) - 1])
		}

		if strings.HasPrefix(line, "method") {
			method := strings.Split(line, " ")[3]
			method = method[:strings.Index(method, "(")]
			methods = append(methods, method)
		}

		if strings.Contains(line, "property") {
			property := strings.Split(line, " ")[3]
			properties = append(properties, property)
		}
	}

	// Process document for diagostics
	for row, line := range strings.Split(text, "\n") {
		if line != "" {
			end_char = line[len(line) - 1:]
			if end_char == "." || end_char == ":" {
				new_line = new_line + " " + line + "\n"
				diagostics = append(diagostics, getDiagnostics(row, new_line)...)
				new_line = ""
			} else if end_char == " " {
				new_line = new_line + line
			} else{
				new_line = new_line + " " + line
			}
		}
	}
	return diagostics
}

func (s *State) OpenDocument(logger *log.Logger, uri, text string, builtin progress.ProgressKeywords) []lsp.Diagnostics {
	s.Documents[uri] = text
	return ProcessDocument(logger, s.Documents[uri], builtin)
}

func (s *State) UpdateDocument(logger *log.Logger, uri string, change lsp.TextDocumentContentChangeEvent, builtin progress.ProgressKeywords) []lsp.Diagnostics {
	s.Documents[uri] = change.Text
	return ProcessDocument(logger, s.Documents[uri], builtin)
}

func (s *State) Hover(id int, uri string, position lsp.Position) lsp.HoverResponse {
	// document := s.Documents[uri]

	return lsp.HoverResponse{
		Response: lsp.Response{
			RPC: "2.0",
			ID:  &id,
		},
		Result:   lsp.HoverResult{
			Contents: fmt.Sprintf("Text: %s ", s.getKeyword(uri, position) ),
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

func createDiagnostics(srow, erow, line int, message string) lsp.Diagnostics {
	return lsp.Diagnostics{
		Range:    LineRange(srow, erow, 0, line),
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
