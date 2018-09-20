# Stocking

[WIP] Minimal Websocket Framework with Server(Go) and Client(Javascript)

> This is the server-side repository. For the client part, please visit [WIP](https://github.com/afterwind-io/stocking).

## Status

**Unstable** - **DO NOT USE IN PRODUCTION**

> The library is right now under early-stage development for experiment only. APIs may vary at any times.

## Installation

```bash
go get github.com/afterwind-io/stocking
```

## Documentation

- [WIP][APIs](https://godoc.org/github.com/afterwind-io/stocking)

## Usage

```go
package main

import "github.com/afterwind-io/stocking"

func main() {
  // Create a standalone server
  server := stocking.NewStocking("localhost:12345", "")

  // Setup routers
  server.On("echo", echo)
  server.On("greet", greet)
  
  // Start serving on ws://localhost:12345/ws
  server.Start()
}

func echo(p stocking.RouterPackage) (interface{}, error) {
  return p.Body, nil
}

func greet(p stocking.RouterPackage) (interface{}, error) {
  body := p.Body

  if body, ok := body.(string); ok {
    return "Hello " + body, nil
  }

  return "Who are you?", nil
}
```

## Protocol

The client and the server communicates based on JSON structures as following:

```typescript
// Client -> Server
interface Inbound {
  route: string
  body: any
}

//Server -> Client
interface Outbound {
  error: string
  body: any
}
```

More details see [protocol.go](protocol.go)

## Roadmap

[WIP]

- [x] At least it is kinda working
- [ ] Broadcast
- [ ] Room
- [ ] Auto-Reconnect

## Trivia

- Why name it "`Stocking`"?
  > ~~I'm a fan of white stocking~~ Just put a "t" in "Sock". Besides, Christmas is coming.

- How do you implement the cascade in the middleware?
  > The cascading control flow is largely inspired by [Koa](https://koajs.com/#application). I managed to mimic similar syntax using `goroute` and `channel`. Checkout [hub.go](hub.go)

## License

[MIT](LICENSE)
