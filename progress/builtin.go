package progress

import (
	"encoding/json"
	"log"
	"progress_ls/lsp"
	"os"
)

type ProgressKeywords struct {
	ProgressKeyword []ProgressKeyword `json:"builtin"`
}

type ProgressKeyword struct {
	Name        string                 `json:"name"`
	Kind        lsp.CompletionItemKind `json:"kind"`
	Description string                 `json:"description"`
}

func (p *ProgressKeywords) IsBuiltin(keyword string) bool {
	for _, builtin := range p.ProgressKeyword {
		if builtin.Name == keyword {
			return true
		}
	}
	return false
}
func (p *ProgressKeywords) GetDescription(keyword string) string {
	for _, builtin := range p.ProgressKeyword {
		if builtin.Name == keyword {
			if builtin.Description != "" {
				return builtin.Description
			} else {
				return builtin.Name
			}
		}
	}
	return ""
}

func InitializeKeywords(logger *log.Logger) ProgressKeywords {
	var keywords ProgressKeywords
	content, err := os.ReadFile("/home/dilip/Gits/progress_ls/res/builtin.json")
	if err != nil {
		logger.Println("Error when opening file: ", err)
	}
	err = json.Unmarshal(content, &keywords)
	if err != nil {
		logger.Println("Error during Unmarshal(): ", err)
	}
	return keywords
}
