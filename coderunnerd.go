package main

import (
	"code.google.com/p/go.net/websocket"
	. "coderunnerd/core"
	"fmt"
	"log"
	"net/http"
	"os"
	"runtime"
	"strconv"
)

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintf(os.Stderr, "Usage: %s <port>", os.Args[0])
		os.Exit(1)
	}
	runtime.GOMAXPROCS(runtime.NumCPU())

	sport, _ := strconv.Atoi(os.Args[1])
	port := uint(sport)

	http.Handle("/", websocket.Handler(CodeHandler))

	log.Printf("Starting listening to %d\n", port)
	if err := http.ListenAndServe(fmt.Sprintf(":%d", port), nil); err != nil {
		log.Fatal("ListenAndServe:", err)
	}

}
