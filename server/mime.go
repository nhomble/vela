package server

import "path/filepath"

var mime = map[string]string{
	".gmi":    "text/gemini",
	".gemini": "text/gemini",
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
