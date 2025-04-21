package rpc_test

import (
	"lsp-go/rpc"
	"testing"
)

type EncodingExample struct {
	Testing bool
}

func TestEncoding(t *testing.T) {
	expected := "Content-Length: 16\r\n\r\n{\"Testing\":true}"
	actual := rpc.EncodeMessage(EncodingExample{Testing:true})

	if expected != actual {
		t.Fatalf("Expected %s, Actual: %s", expected, actual)
	}
}

func TestDecodeing(t *testing.T) {
	incomingMessage := "Content-Length: 15\r\n\r\n{\"Method\":\"hi\"}"
	method, content, err := rpc.DecodeMessage([]byte(incomingMessage))
	contentLen := len(content)

	if err != nil {
		t.Fatal(err)
	}

	if contentLen != 15 {
		t.Fatalf("Expected: 15, Got: %d",contentLen)
	}

	if method != "hi" {
		t.Fatalf("Expected: 'hi', Got: %s", method)
	}
}
