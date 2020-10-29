vela
====
![Go](https://github.com/nhomble/vela/workflows/Go/badge.svg)

a [gemini](https://gemini.circumlunar.space/docs/specification.html) server

# Usage
```go
package main

import (
	"github.com/nhomble/gemini-server/server"
	"net/http"
"os"
)

func main(){
    root, _ := os.Getwd()
	server.ListenAndServe("", "path/to/public.pem", "path/to/private.pem", server.FileServingRequestHandler{Root: root})
}
```

