package main

import (
	"bytes"
	"github.com/makeworld-the-better-one/go-gemini"
	"github.com/nhomble/gemini-server/server"
	"os"
	"strings"
	"testing"
	"time"
)

func TestMain(m *testing.M) {
	go start()
	time.Sleep(time.Second * 2)
	exitVal := m.Run()
	os.Exit(exitVal)
}

func TestGetFile(t *testing.T) {
	c := &gemini.Client{Insecure: true}
	resp, _ := c.Fetch("gemini://localhost:1965/resources/smaller.gmi")
	buf := new(bytes.Buffer)
	buf.ReadFrom(resp.Body)
	newStr := buf.String()
	if resp.Status != 20 {
		t.Fatalf("Didn't get 20 response, we got '%d'", resp.Status)
	}
	if !strings.Contains(newStr, "web") {
		t.Fatalf("Didn't see web, instead we got '%s'", newStr)
	}
}

func TestNoFile(t *testing.T) {
	c := &gemini.Client{Insecure: true}
	resp, _ := c.Fetch("gemini://localhost:1965/nope")
	if resp.Status != 40 {
		t.Fatalf("Didn't get a 40, instead we got '%d'", resp.Status)
	}
}

func start() {
	root, _ := os.Getwd()
	server.ListenAndServe("", "public.pem", "private.pem", server.FileServingRequestHandler{Root: root})
}
