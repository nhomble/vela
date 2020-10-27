package spec

import (
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
	Body   *string
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
		c.Write([]byte(*resp.Body))
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

//The requested resource accepts a line of textual user input.
//The <META> line is a prompt which should be displayed to the user.
//The same resource should then be requested again with the user's input included as a query component.
//Queries are included in requests as per the usual generic URL definition in RFC3986, i.e.
//separated from the path by a ?. Reserved characters used in the user's input must be "percent-encoded" as
//per RFC3986, and space characters should also be percent-encoded.
func InputStatus(prompt string) *Status {
	return &Status{
		Code:     "10",
		Metadata: prompt,
	}
}

//The request was handled successfully and a response body will follow the response header.
//The <META> line is a MIME media type which applies to the response body.
func SuccessStatus(mime string) *Status {
	return &Status{
		Code:     "20",
		Metadata: mime,
	}
}

//The server is redirecting the client to a new location for the requested resource.
//There is no response body. <META> is a new URL for the requested resource.
//The URL may be absolute or relative. The redirect should be considered temporary,
//i.e. clients should continue to request the resource at the original address and should not
//performance convenience actions like automatically updating bookmarks. There is no response body.
func RedirectStatus(info string) *Status {
	return &Status{
		Code:     "30",
		Metadata: info,
	}
}

//The request has failed. There is no response body. The nature of the failure is
//temporary, i.e. an identical request MAY succeed in the future. The contents of
//<META> may provide additional information on the failure, and should be displayed to human users.
func TemporaryFailure(info string) *Status {
	return &Status{
		Code:     "40",
		Metadata: info,
	}
}

//The request has failed. There is no response body. The nature of the failure is permanent, i.e. identical
//future requests will reliably fail for the same reason. The contents of <META> may provide additional
//information on the failure, and should be displayed to human users. Automatic clients such as
//aggregators or indexing crawlers should not repeat this request.
func PermanentFailure(info string) *Status {
	return &Status{
		Code:     "50",
		Metadata: info,
	}
}

//The requested resource requires a client certificate to access. If the request was made without
//a certificate, it should be repeated with one. If the request was made with a certificate, the server
//did not accept it and the request should be repeated with a different certificate. The contents
//of <META> (and/or the specific 6x code) may provide additional information on certificate requirements or
//the reason a certificate was rejected.
func CertificateFailure(info string) *Status {
	return &Status{
		Code:     "60",
		Metadata: info,
	}
}
