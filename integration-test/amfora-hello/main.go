package main

import (
	"bytes"
	"github.com/makeworld-the-better-one/go-gemini"
	"github.com/nhomble/gemini-server/server"
	"os"
	"strings"
	"time"
)

func main() {
	go start()
	time.Sleep(time.Second * 2)
	c := &gemini.Client{Insecure: true}
	resp, _ := c.Fetch("gemini://localhost:1965/integration-test/resources/smaller.gmi")
	buf := new(bytes.Buffer)
	buf.ReadFrom(resp.Body)
	newStr := buf.String()

	if !strings.Contains(newStr, "web") {
		panic("Didn't see web, instead we go '" + newStr + "'")
	}
}

func start() {
	root, _ := os.Getwd()
	server.ListenAndServe("", "integration-test/amfora-hello/public.pem", "integration-test/amfora-hello/private.pem", server.FileServingRequestHandler{Root: root})
}
