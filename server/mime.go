package server

import (
	"github.com/nhomble/gemini-server/spec"
	"path/filepath"
)

var mime = map[string]string{
	".gmi":    spec.GEMINI_MIME,
	".gemini": spec.GEMINI_MIME,
}

func ChooseMime(path string) string {
	extension := filepath.Ext(path)
	val, ok := mime[extension]
	if ok {
		return val
	} else {
		return "plain/text"
	}
}
