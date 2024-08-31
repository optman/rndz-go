# rndz-go

A simple rendezvous protocol implementation to help NAT traversal or hole punching.

Golang implementation for [rndz](https://github.com/optman/rndz).

## Table of Contents

- [Features](#features)
- [Installation](#installation)
- [Usage](#usage)
  - [TCP](#tcp)
  - [UDP](#udp)
- [License](#license)

## Features

- Simple and lightweight implementation
- Supports both TCP and UDP protocols
- Easy to integrate with existing Go projects

## Installation

To install `rndz-go`, use `go get`:

```sh
go get github.com/optman/rndz-go
```

## Usage

Here's how you can use `rndz-go` to perform NAT traversal or hole punching.

### TCP

#### Client 1

```go
import (
	"context"
	"net/netip"

	"github.com/optman/rndz-go/client/tcp"
)

func main() {
	rndzServer := "your-rendezvous-server"
	c := tcp.NewClient(rndzServer, "c1", netip.AddrPort{})
	defer c.Close()
	l, _ := c.Listen(context.Background())
	defer l.Close()
	for {
		conn, _ := l.Accept()
		defer conn.Close()
		// Handle connection
	}
}
```

#### Client 2

```go
import (
	"context"
	"net/netip"

	"github.com/optman/rndz-go/client/tcp"
)

func main() {
	rndzServer := "your-rendezvous-server"
	c := tcp.NewClient(rndzServer, "c2", netip.AddrPort{})
	defer c.Close()
	conn, _ := c.Connect(context.Background(), "c1")
	defer conn.Close()
	// Use connection
}
```

### UDP

```go
import (
	"net/netip"

	"github.com/optman/rndz-go/client/udp"
)

func main() {
	rndzServer := "your-rendezvous-server"
	id := "your-client-id"
	c := udp.NewClient(rndzServer, id, netip.AddrPort{})
	// Use client
}
```

## License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.