package server

import (
	"fmt"
	"github.com/nhomble/gemini-server/spec"
	"os"
	"path"
)

type FileServingRequestHandler struct {
	Root string
}

func (handler FileServingRequestHandler) Handle(request *spec.Request) *spec.Response {
	complete := path.Join(handler.Root, request.URL.Path)
	f, err := os.Open(complete)
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
	resp.Body = f
	return &resp
}
