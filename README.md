# Stocking

[WIP] Minimal Websocket Framework with Server(Go) and Client(Javascript)

> This is the server-side repository. For the client part, please visit [WIP](https://github.com/afterwind-io/stocking).

## Status

**Unstable** - **DO NOT USE IN PRODUCTION**

> The library is right now under early-stage development for experiment only. APIs may vary at any time.

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
  server.On("echo", echo, nil)
  server.On("greet", greet, person{})
  server.Otherwise(otherwise)
  
  // Start serving on ws://localhost:12345/ws
  server.Start()
}

type person struct {
  Name string `json:"name"`
}

func echo(p stocking.RouterPackage) (interface{}, error) {
  return p.Body, nil
}

func greet(p stocking.RouterPackage) (interface{}, error) {
  body, _ := p.Body.(*person)

  if body.Name == "doge" {
    return "Hello " + body.Name, nil
  }

  return "Who are you?", nil
}

func otherwise(p stocking.RouterPackage) (interface{}, error) {
  return "oops", nil
}
```

## Protocol

The client and the server communicates based on string structure as following:

> `[MessageType],[Ack],[Content]`
>
> eg: `6,0,{r: "echo", p: "doge"}`

Where:

- MessageType: `string`
- Ack: `string`
- Content: [See below]

```go
// Content type when sent to server
type TextMessageInboundProtocol struct {
  Route   string          `json:"r"`
  Payload json.RawMessage `json:"p"`
}

// Content type when sent to client
type TextMessageOutboundProtocol struct {
  Code    int         `json:"c"`
  Payload interface{} `json:"p"`
}
```

> In fact, the type of inbound `Payload` is `interface{}`.
>
> `json.RawMessage` ensures that we can unmarshal it to the actual type we need later.

More details see [protocol.go](protocol.go)

## Middleware

[WIP]

Middleware design in `Stocking` is largely inspired by [Koa](https://koajs.com/#application). I managed to mimic similar syntax using `goroute` and `channel`. The following logger example shows how basic middleware looks like:

```go
func (me *mLogger) Handle(p *HubPackge, next MiddlewareStepFunc) {
  // [Downstream Jobs] Do some inbound logs here
  log.Println(fmt.Sprintf("<-- [%v] %v, %v, %v", p.client.id, p.mtype, p.ack, p.content))

  // Temporarily suspended and pass the control to the next middleware
  done := <-next(nil)
  
  // Meanwhile downstream middlewares are doing their own jobs now

  // After they finished their jobs, the control flow rewinds,
  // and middlewares are resumed to perform its upstream jobs, in reverse order

  // [Upstream Jobs] Now grab outbound messages and log it
  log.Println(fmt.Sprintf("--> [%v] %v, %v, %v", p.client.id, p.mtype, p.ack, p.content))

  // duty ends here
  done <- nil
}
```

## Roadmap

- [x] At least it is kinda working
- [ ] Broadcast
- [ ] Room
- [ ] Auto-Reconnect

## Known Flaws

- If an error occurs and breaks the `Read`/`Write` loop, there's currently no graceful way to stop the other one. -> [client.go](client.go)
- If Client route handler panics, the whole server just faint away :( -> [mRouter.go](mRouter.go)
- Route handling stops the world. -> [hub.go](hub.go)

## Trivia

- Why name it "`Stocking`"?
  > ~~I'm a fan of white stocking~~ Just put a "t" in "Sock". Besides, Christmas is coming.

## License

[MIT](LICENSE)
