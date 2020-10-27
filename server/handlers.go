package server

import (
	"fmt"
	"github.com/nhomble/gemini-server/spec"
	"io/ioutil"
	"os"
	"path"
)

type FileServingRequestHandler struct {
	Root string
}

func (handler FileServingRequestHandler) Handle(request *spec.Request) *spec.Response {
	complete := path.Join(handler.Root, request.URL.Path)
	data, err := ioutil.ReadFile(complete)
	resp := spec.Response{}

	if os.IsNotExist(err) {
		resp.Status = spec.TemporaryFailure("Does not exist")
		return &resp
	}

	if err != nil {
		resp.Status = spec.PermanentFailure(fmt.Sprintf("%v", err))
		return &resp
	}

	resp.Status = spec.SuccessStatus(ChooseMime(complete))
	loaded := string(data)
	resp.Body = &loaded
	return &resp
}
