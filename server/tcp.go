package server

import (
	"bufio"
	"crypto/tls"
	"fmt"
	"github.com/nhomble/gemini-server/spec"
	"log"
	"net"
	"net/url"
	"strconv"
	"time"
)

type RequestHandler interface {
	Handle(r *spec.Request) *spec.Response
}

func ListenAndServe(address string, pem string, key string, handler RequestHandler) {
	if address == "" {
		address = ":" + strconv.Itoa(spec.DEFAULT_PORT)
	}
	cert, err := tls.LoadX509KeyPair(pem, key)
	if err != nil {
		panic(err)
	}
	config := tls.Config{Certificates: []tls.Certificate{cert}}
	listener, err := tls.Listen("tcp", address, &config)
	defer listener.Close()
	if err != nil {
		panic(err)
	}
	for {
		c, err := listener.Accept()
		if err != nil {
			fmt.Println(err)
			return
		}
		go handleConnection(handler, c)
	}
}

func handleConnection(handler RequestHandler, c net.Conn) {
	defer c.Close()
	start := time.Now()
	addr := c.RemoteAddr().String()
	log.Printf("Received request from address=%s\n", addr)
	var status *spec.Status
	var body *string
	var target *url.URL
	tlscon, ok := c.(*tls.Conn)
	if !ok {
		status = spec.CertificateFailure("Did not establish tls connection")
	} else {
		line, isprefix, err := bufio.NewReader(tlscon).ReadLine()
		for {

			if isprefix {
				status = spec.PermanentFailure("Request was too long for buffer")
				break
			}
			if err != nil {
				status = spec.PermanentFailure(fmt.Sprintf("%v", err))
				break
			}
			isvalid, reason := spec.IsValidRequest(line)
			if !isvalid {
				status = spec.PermanentFailure(reason)
				break
			}

			target, err = url.Parse(string(line))
			if err != nil {
				status = spec.PermanentFailure(fmt.Sprintf("%v", err))
				break
			}

			if target.Scheme != spec.SCHEME {
				status = spec.PermanentFailure(fmt.Sprintf("Scheme=%s is invalid\n", target.Scheme))
				break
			}

			resp := handler.Handle(&spec.Request{
				URL: target,
			})
			status = resp.Status
			body = resp.Body
			break
		}
	}
	(&spec.Response{
		Status: status,
		Body:   body,
	}).WriteTo(tlscon)
	end := time.Now()
	t := ""
	if target != nil {
		t = target.String()
	}
	log.Printf("Completed request from addr=%s destination=%s in duration=%d ms\n", addr, t, end.Sub(start).Milliseconds())
}
