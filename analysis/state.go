package analysis

import (
	"fmt"
	"log"
	"progress_ls/lsp"
	"progress_ls/progress"
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

func getWord(idx int, line string) string {
	var text string
	for _, char := range line[idx:] {
		if string(char) != "." && string(char) != "(" && string(char) != " " {
			text = text + string(char)
		} else {
			break
		}
	}
	return text
}

func getDiagnostics(logger *log.Logger,row int, line string, types progress.ProgressTypes) []lsp.Diagnostics {
	diagostics := []lsp.Diagnostics{}
	found := false
	trimmed_text := strings.Trim(line, " ")
	var datatype string
	var propertyMethod string
	var idx int
	var text string

	if len(strings.Split(trimmed_text, ". ")) > 1 {
		diagostics = append(diagostics, createError(row, 0, len(trimmed_text), "single line cannot contains multiple statements"))
	}

	for _, splits := range strings.Split(trimmed_text[:len(trimmed_text) - 2]," ") {
		if splits != "" {
			text = text + " " + splits
		}
	}

	if strings.Contains(text, "define") {
		if strings.Contains(text, "property") || strings.Contains(text, "variable") {
			if !strings.Contains(text, "no-undo") {
				idx = strings.Index(line, "define") - 1
				diagostics = append(diagostics, createError(row, idx, idx + len(text), "no-undo is missing"))
			}

			split := strings.Split(strings.Trim(text," "), " ")

			if split[1] != progress.ProgressPrivate && split[1] != progress.ProgressProtected && split[1] != progress.ProgressPublic {
				datatype = split[4]
			} else {
				datatype = split[5]
			}

			if progress.IsDefaultDataType(datatype) {
				found = true
			} else {
				for _, class := range types.Classes {
					if strings.Contains(class, datatype) {
						found = true
						break
					}
				}
			}

			if !found {
				idx = strings.Index(line, datatype) - 1
				diagostics = append(diagostics, createError(row, idx, idx + len(datatype), "class is not imported. Import the class with using statement"))
			}
		}
	}

	if idx = progress.IndexOfRestrictedText(line) ; idx > 0 {
		diagostics = append(diagostics, createWarning(row, idx - 1, idx + len(getWord(idx, line)) - 1, "incorrect format for variable name, use PascalCase or camelCase format"))
	}

	if strings.Contains(text, "this-object:") {
		for _, this := range types.Methods {
			if strings.Contains(text, this) {
				found = true
				break
			}
		}

		for _, this := range types.Properties {
			if strings.Contains(text, this) {
				found = true
				break
			}
		}

		if !found {
			for _, char := range line[strings.Index(line, ":") + 1:] {
				if string(char) != "." && string(char) != "(" && string(char) != " " {
					propertyMethod = propertyMethod + string(char)
				} else {
					break
				}
			}
			idx = strings.Index(line, propertyMethod) - 1
			diagostics = append(diagostics, createError(row, idx, idx + len(propertyMethod), "undefined property/method"))
		}
	}

	return diagostics
}

func (s *State) getKeyword(uri string, postion lsp.Position) string {
	text := strings.Split(s.Documents[uri],"\n")[postion.Line]
	text = strings.Split(text[postion.Character:], " ")[0]
	return text
}

func ProcessDocument(logger *log.Logger, uri string, text string, builtin progress.ProgressKeywords) []lsp.Diagnostics {
	var new_line   string
	var end_char   string
	var classes    []string
	var methods    []string
	var properties []string
	var types 		progress.ProgressTypes
	diagostics := []lsp.Diagnostics{}

	text = strings.Trim(text, " ")

	// Process document for listing classes, methods, properties
	for _, line := range strings.Split(text, "\n") {
		line = strings.Trim(line, " ")
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

	types = progress.ProgressTypes{
		URI:        uri,
		Classes:    classes,
		Methods:    methods,
		Properties: properties,
	}

	// Process document for diagostics
	for row, line := range strings.Split(text, "\n") {
		if line != "" {
			end_char = line[len(line) - 1:]
			if end_char == "." || end_char == ":" {
				new_line = new_line + " " + line + "\n"
				diagostics = append(diagostics, getDiagnostics(logger, row, new_line, types)...)
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
	return ProcessDocument(logger, uri, s.Documents[uri], builtin)
}

func (s *State) UpdateDocument(logger *log.Logger, uri string, change lsp.TextDocumentContentChangeEvent, builtin progress.ProgressKeywords) []lsp.Diagnostics {
	s.Documents[uri] = change.Text
	return ProcessDocument(logger, uri, s.Documents[uri], builtin)
}

func (s *State) Hover(id int, uri string, position lsp.Position, builtin progress.ProgressKeywords) lsp.HoverResponse {
	keyword := s.getKeyword(uri, position)
	var content string
	var desc string
	if builtin.IsBuiltin(keyword) {
		desc = builtin.GetDescription(keyword)
	}
	content = fmt.Sprintf("# %s\n---\n%s", keyword, desc)

	return lsp.HoverResponse{
		Response: lsp.Response{
			RPC: "2.0",
			ID:  &id,
		},
		Result:   lsp.HoverResult{
			Contents: content,
			// Contents: fmt.Sprintf("Text: %s ", s.getKeyword(uri, position) ),
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

func createError(row, start, end int, message string) lsp.Diagnostics {
	return lsp.Diagnostics{
		Range:    LineRange(row, start, end),
		Severity: 1,
		Source:   "progress_ls",
		Message:  message,
	}
}

func createWarning(row, start, end int, message string) lsp.Diagnostics {
	return lsp.Diagnostics{
		Range:    LineRange(row, start, end),
		Severity: 2,
		Source:   "progress_ls",
		Message:  message,
	}
}

func LineRange(line, start, end int) lsp.Range {
	return lsp.Range{
		Start: lsp.Position{
			Line:      line,
			Character: start,
		},
		End:   lsp.Position{
			Line:      line,
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
					Range: LineRange(row, idx, idx + len("VS Code")),
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
					Range:   LineRange(row, idx, idx + len("VS Code")),
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
