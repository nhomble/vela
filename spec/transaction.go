package spec

import (
	"io"
	"io/ioutil"
	"net"
	"net/url"
	"strconv"
)

const SCHEME = "gemini"
const CRLF = "\r\n"
const MAX_LEN = 1024

func IsValidRequest(line []byte) (bool, string) {
	l := len(line)
	if l < 2 {
		return false, "line is too short, must at least contain <CRLF>"
	}
	if l > MAX_LEN {
		return false, "line too long"
	}

	return true, ""
}

type Response struct {
	Status *Status
	Body   io.ReadCloser
}

type Request struct {
	URL *url.URL
}

func (resp *Response) WriteTo(c net.Conn) {
	c.Write([]byte(resp.Status.Code))
	c.Write([]byte{0x20})
	c.Write([]byte(resp.Status.Metadata))
	c.Write([]byte(CRLF))
	if resp.Status.isSuccess() {
		b, _ := ioutil.ReadAll(resp.Body)
		c.Write(b)
		c.Write([]byte(CRLF))
	}
}

type Status struct {
	Code     string
	Metadata string
}

func (status *Status) isSuccess() bool {
	i, _ := strconv.ParseInt(status.Code, 10, 64)
	return i >= 20 && i < 30
}
