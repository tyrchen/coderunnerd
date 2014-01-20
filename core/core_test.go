package core

import (
	"code.google.com/p/go.net/websocket"
	"encoding/json"
	"strings"
	"testing"
)

const (
	HOST = "localhost:8210/"
)

const (
	MAX_BYTES = 64000
)

const (
	EXPECTED = "Hello world!\n"
	GOCODE   = `
package main
import (
	"fmt"
)
func main() {
	fmt.Printf("Hello world!\n")
}	
`
	PYCODE = `
import sys
if __name__ == "__main__":
	print("Hello world!")	
`
)

func getCode(suffix string, content string) string {
	var val JsonCode
	val.Suffix = suffix
	val.Content = content

	b, _ := json.Marshal(val)
	return string(b)
}

func testCode(suffix string, content string) bool {
	var ret JsonCode
	msg := getCode(suffix, content)
	output, _ := runCode(msg)
	json.Unmarshal(output, &ret)

	return ret.Suffix == suffix && strings.Contains(ret.Content, EXPECTED)
}

func TestGolang(t *testing.T) {
	if !testCode("go", GOCODE) {
		t.Errorf("Failed to run golang code")
	}
}

func TestPython(t *testing.T) {
	if !testCode("py", PYCODE) {
		t.Errorf("Failed to run python code")
	}
}

func TestWebsocket(t *testing.T) {
	var ret JsonCode

	ws, err := websocket.Dial("ws://"+HOST, "", "http://"+HOST)
	if err != nil {
		panic(err)
	}

	resp := make([]byte, MAX_BYTES)

	msg := getCode("go", GOCODE)

	n, err := ws.Write([]byte(msg))
	if err != nil {
		panic(err)
	}

	n, err = ws.Read(resp)
	if err != nil {
		panic(err)
	}

	json.Unmarshal(resp[:n], &ret)
	if !strings.Contains(ret.Content, EXPECTED) {
		t.Errorf("Failed to run golang code through websocket")
	}
}
