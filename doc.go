// Package rndz provides a simple rendezvous protocol implementation to help NAT traversal or hole punching.
//
// To connect a node behind a firewall or NAT (such as a home gateway), which only allows outbound connections,
// you not only need to know its gateway IP, but also have the node send you traffic first.
//
// This applies not only to IPv4 but also to IPv6. As IPv4 needs to deal with NAT, both need to deal with firewalls.
//
// # How rndz works
//
// Setup a publicly accessible server as a rendezvous point, to observe all peers' addresses and forward connection requests.
//
// Each peer needs a unique identity. The server will associate the identity with the observed address.
// A listening peer will keep pinging the server and receive forwarded requests. A connecting peer will
// request the server to forward its connection requests.
//
// After receiving a forwarded connection request from the server, the listening peer will send a dummy packet
// to the connecting peer. This will open the firewall or NAT rule for the connecting peer;
// otherwise, all packets from the peer will be blocked.
//
// After that, we return native socket types `net.Conn` and `net.UDPConn` to the caller.
//
// The essential part is that we must use the same port to communicate with the rendezvous server and peers.
//
// The implementation depends on socket options SO_REUSEADDR and SO_REUSEPORT, so it is OS dependent.
// For TCP, the OS should allow the listening socket and connecting socket to bind to the same port.
// For UDP, the OS should correctly dispatch traffic to connected and unconnected UDP sockets all binding to the same port.
package rndz
