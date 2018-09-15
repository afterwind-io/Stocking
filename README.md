# Stocking

[WIP] Minimal Websocket Framework with Back-end(Go) and Front-end(Javascript)

> This is the server-side repository. For the front-end part, please visit [WIP](https://github.com/afterwind-io/stocking)

## Usage

```go
package main

import "github.com/afterwind-io/stocking"

// Start a standalone server
func main() {
  server := stocking.NewStocking("", "")
  server.Start()
}

// Done, serving on localhost:12345/ws
```
