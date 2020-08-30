# golang-jsonrpc2

Implementation for JSON-RPC2 protocol

Full specification: https://www.jsonrpc.org/specification


# HTTP example

## Server
```go
package main

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	jrpc "github.com/gumeniukcom/golang-jsonrpc2"
)

func main() {

	serv := jrpc.New()

	if err := serv.RegisterMethod("sum", sum); err != nil {
		panic(err)
	}

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		ctx := context.Background()
		body, err := ioutil.ReadAll(r.Body)
		if err != nil {
			panic(err)
		}
		defer r.Body.Close()

		w.Header().Set("Content-Type", "applicaition/json")
		w.WriteHeader(http.StatusOK)
		if _, err = w.Write(serv.HandleRPCJsonRawMessage(ctx, body)); err != nil {
			panic(err)
		}
	})

	if err := http.ListenAndServe(":8088", nil); err != nil {
		panic(err)
	}
}

type income struct {
	A int `json:"a"`
	B int `json:"b"`
}
type outcome struct {
	Sum int `json:"sum"`
}

func sum(ctx context.Context, data json.RawMessage) (json.RawMessage, int, error) {
	if data == nil {
		return nil, jrpc.InvalidRequestErrorCode, fmt.Errorf("empty request")
	}
	inc := &income{}
	err := json.Unmarshal(data, inc)
	if err != nil {
		return nil, jrpc.InvalidRequestErrorCode, err
	}

	C := outcome{
		Sum: inc.A + inc.B,
	}

	mdata, err := json.Marshal(C)
	if err != nil {
		return nil, jrpc.InternalErrorCode, err
	}
	return mdata, jrpc.OK, nil
}
```

## Request


```bash
curl -d '{"jsonrpc":"2.0", "id":"qwe", "method":"sum", "params":{"a":5, "b":3}}' -H "Content-Type: application/json" -X POST http://localhost:8088/
``` 
## Response

```json
{"jsonrpc":"2.0","result":{"sum":8},"id":"qwe"}
```

