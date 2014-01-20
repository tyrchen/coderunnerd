package core

import (
	"bytes"
	"code.google.com/p/go.net/websocket"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
)

type Lang struct {
	exec string
	opt  string
}

type JsonCode struct {
	Suffix  string // suffix indicating the type of the code
	Content string
}

var (
	lang   map[string]Lang = make(map[string]Lang)
	tmpdir string
	uniq   = make(chan int)
)

func init() {
	var err error
	tmpdir, err = filepath.EvalSymlinks(os.TempDir())
	if err != nil {
		log.Fatal(err)
	}

	// generate uniq numbers
	go func() {
		for i := 0; ; i++ {
			uniq <- i
		}
	}()
}

func CodeHandler(ws *websocket.Conn) {
	var err error

	for {
		var data string

		if err = websocket.Message.Receive(ws, &data); err != nil {
			fmt.Println("Can't receive")
			break
		}

		fmt.Println("Received code from client: " + data)

		if out, err := runCode(data); err == nil {
			resultSend(ws, string(out))
		}
	}
}

// run executes the specified command and returns its output and an error.
func run(dir string, args ...string) ([]byte, error) {
	var buf bytes.Buffer
	cmd := exec.Command(args[0], args[1:]...)
	cmd.Dir = dir
	cmd.Stdout = &buf
	cmd.Stderr = cmd.Stdout
	err := cmd.Run()
	return buf.Bytes(), err
}

func runCode(data string) (out []byte, err error) {
	var code, ret JsonCode
	json.Unmarshal([]byte(data), &code)

	if handler, ok := lang[code.Suffix]; ok {
		x := filepath.Join(tmpdir, "compile"+strconv.Itoa(<-uniq))
		src := x + "." + code.Suffix
		defer os.Remove(src)
		if err = ioutil.WriteFile(src, []byte(code.Content), 0666); err != nil {
			return
		}
		// run the code
		dir, file := filepath.Split(src)

		data, err := run(dir, handler.exec, handler.opt, file)
		if err != nil {
			return out, err
		}
		ret.Suffix = code.Suffix
		ret.Content = string(data)
		out, _ = json.Marshal(ret)
		return out, nil
	}
	err = errors.New("Handler not registered")

	return
}

func resultSend(ws *websocket.Conn, msg string) {
	if err := websocket.Message.Send(ws, msg); err != nil {
		fmt.Println("Can't send")
	}
}
