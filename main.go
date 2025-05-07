package main

import (
	"bufio"
	"encoding/json"
	"io"
	"log"
	"progress_ls/analysis"
	"progress_ls/lsp"
	"progress_ls/progress"
	"progress_ls/rpc"
	"os"
)

func main() {
	logger := getLogger("/home/dilip/Gits/progress_ls/log.txt")
	logger.Println("Lsp Started")
	scanner := bufio.NewScanner(os.Stdin)
	scanner.Split(rpc.Split)
	state := analysis.NewState()
	builtin := progress.InitializeKeywords(logger)
	writer := os.Stdout

	for scanner.Scan() {
		msg := scanner.Bytes()
		method, content, err := rpc.DecodeMessage([]byte(msg))
		if err != nil {
			logger.Printf("Ooops! Got error: %s",err)
			continue
		}

		handleMessage(logger, writer, state, method, content, builtin)
	}
}

func handleMessage(logger *log.Logger, writer io.Writer, state analysis.State, method string, content []byte, builtin progress.ProgressKeywords) {
	logger.Printf("We recived mesage with method: %s", method)

	switch method {
	case "initialize":
		var request lsp.InitializeRequest
		var response lsp.InitializeResponse
		if err := json.Unmarshal(content, &request); err != nil {
			logger.Printf("hey, we conuldn't parse this: %s", err)
		}
		logger.Printf(
			"Connected to: %s %s",
			request.Params.ClientInfo.Name,
			request.Params.ClientInfo.Version)

		response = lsp.NewInitializeResponse(request.ID)
		logger.Printf("initialize: %s",rpc.EncodeMessage(response))
		writeResponse(writer, response)
		logger.Print("Sent the reply to neovim")
	case "textDocument/didOpen":
		var request lsp.DidOpenTextDocumentNotification
		if err := json.Unmarshal(content, &request); err != nil {
			logger.Printf("textDocument/didOpen: %s", err)
			return
		}
		logger.Printf("Opened: %s", request.Params.TextDocument.URI)
		diagnostics := state.OpenDocument(logger, request.Params.TextDocument.URI, request.Params.TextDocument.Text, builtin)
		writeResponse(writer, lsp.PublishDiagosticsNotification{
			Notification: lsp.Notification{
				RPC:    "2.0",
				Method: "textDocument/publishDiagnostics",
			},
			Params:       lsp.PublishDiagosticsParams{
				URI:         request.Params.TextDocument.URI,
				Diagnostics: diagnostics,
			},
		})
	case "textDocument/didChange":
		var request lsp.DidChangeTextDocumentNotification
		if err := json.Unmarshal(content, &request); err != nil {
			logger.Printf("textDocument/didChange: %s", err)
			return
		}
		logger.Printf("Changed: %s", request.Params.TextDocument.URI)
		logger.Printf("didChange: %s",content)
		for _, change := range request.Params.ContentChanges {
			diagnostics := state.UpdateDocument(logger, request.Params.TextDocument.URI, change, builtin)
			writeResponse(writer, lsp.PublishDiagosticsNotification{
				Notification: lsp.Notification{
					RPC:    "2.0",
					Method: "textDocument/publishDiagnostics",
				},
				Params:       lsp.PublishDiagosticsParams{
					URI:         request.Params.TextDocument.URI,
					Diagnostics: diagnostics,
				},
			})
		}
	case "textDocument/hover":
		var request lsp.HoverRequest
		var response lsp.HoverResponse
		if err := json.Unmarshal(content, &request); err != nil {
			logger.Printf("textDocument/hover: %s", err)
			return
		}
		response = state.Hover(request.ID, request.Params.TextDocument.URI, request.Params.Position, builtin)
		writeResponse(writer, response)
	case "textDocument/definition":
		var request lsp.DefinitionRequest
		var response lsp.DefinitionResponse
		if err := json.Unmarshal(content, &request); err != nil {
			logger.Printf("textDocument/definition: %s", err)
			return
		}
		response = state.Definition(request.ID, request.Params.TextDocument.URI, request.Params.Position)
		writeResponse(writer, response)
	case "textDocument/codeAction":
		var request lsp.CodeActionRequest
		var response lsp.CodeActionResponse
		if err := json.Unmarshal(content, &request); err != nil {
			logger.Printf("textDocument/codeAction: %s", err)
			return
		}
		response = state.CodeAction(request.ID, request.Params.TextDocument.URI)
		writeResponse(writer, response)
	case "textDocument/completion":
		var request lsp.CompletionRequest
		var response lsp.CompletionResponse
		if err := json.Unmarshal(content, &request); err != nil {
			logger.Printf("textDocument/completion: %s", err)
			return
		}
		response = state.Completion(request.ID, request.Params.TextDocument.URI)
		writeResponse(writer, response)
	}
}

func writeResponse(writer io.Writer, msg any) {
	reply := rpc.EncodeMessage(msg)
	writer.Write([]byte(reply))
}

func getLogger(filename string) *log.Logger {
	logfile, err := os.OpenFile(filename, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0666)

	if err != nil {
		panic("invalid file")
	}

	return log.New(logfile, "[lsp-go]", log.Ldate|log.Ltime|log.Lshortfile)
}
