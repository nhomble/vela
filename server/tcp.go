package server

import (
	"bufio"
	"crypto/tls"
	"fmt"
	"github.com/nhomble/gemini-server/spec"
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
	var target *url.URL
	resp := &spec.Response{}
	tlscon, ok := c.(*tls.Conn)
	if !ok {
		resp.Status = spec.CertificateFailure("Did not establish tls connection")
	} else {
		line, isprefix, err := bufio.NewReader(tlscon).ReadLine()
		for {

			if isprefix {
				resp.Status = spec.PermanentFailure("Request was too long for buffer")
				break
			}
			if err != nil {
				resp.Status = spec.PermanentFailure(fmt.Sprintf("%v", err))
				break
			}
			isvalid, reason := spec.IsValidRequest(line)
			if !isvalid {
				resp.Status = spec.PermanentFailure(reason)
				break
			}

			target, err = url.Parse(string(line))
			if err != nil {
				resp.Status = spec.PermanentFailure(fmt.Sprintf("%v", err))
				break
			}

			if target.Scheme != spec.SCHEME {
				resp.Status = spec.PermanentFailure(fmt.Sprintf("Scheme=%s is invalid\n", target.Scheme))
				break
			}

			resp = handler.Handle(&spec.Request{
				URL: target,
			})
			if resp.Body != nil {
				defer resp.Body.Close()
			}
			break
		}
	}
	sizeOf := resp.WriteTo(tlscon)

	clientAddr := ""
	endpoint := ""
	if target != nil {
		clientAddr = target.Hostname()
		endpoint = target.Path
	}
	accessLog := CommonLog{
		clientHost:   clientAddr,
		time:         start,
		url:          endpoint,
		status:       resp.Status.Code,
		responseSize: sizeOf,
	}.format()
	fmt.Println(accessLog)
}
