
# rndz-go

A simple rendezvous protocol implementation to help NAT traversal or hole punching.

Golang implementation for [rndz](https://github.com/optman/rndz)


### tcp

client1
```go
import "github.com/optman/rndz-go/client/tcp"

c := tcp.NewClient(rndzServer, "c1", netip.AddrPort{})
defer c.Close()
l, _ := c.Listen(context.Background())
defer l.Close()
for {
	conn, _ := l.Accept()
    defer conn.Close()
    ...
}
```

client2
```go
import "github.com/optman/rndz-go/client/tcp"

c := tcp.NewClient(rndzServer, "c2", netip.AddrPort{})
defer c.Close()
conn, _ := c.Connect(context.Background(), "c1")
defer conn.Close()
```

### udp 

```go
import "github.com/optman/rndz-go/client/udp"

c := udp.NewClient(rndzServer, id, netip.AddrPort{})
...
```
