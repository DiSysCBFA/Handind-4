package peer

import (
	"context"
	"log"
	"net"
	"sync"
	"time"

	h4 "github.com/DiSysCBFA/Handind-4/h4"
	"google.golang.org/grpc"
)

// Peer represents a single peer in the network
type Peer struct {
	h4.UnimplementedH4Server
	Id         int
	port       string
	requested  bool
	grpcServer *grpc.Server
	mu         sync.Mutex
}

// NewPeer creates a new Peer instance
func NewPeer(id int, port string, totalPeers int) *Peer {
	peer := &Peer{
		Id:         id,
		port:       port,
		grpcServer: grpc.NewServer(),
		requested:  false,
	}
	h4.RegisterH4Server(peer.grpcServer, peer)
	return peer
}

// Multicast sends a request message to all specified peer ports
func (p *Peer) Multicast(ports []string) {
	req := &h4.Message{
		Id:        int64(p.Id),
		Timestamp: time.Now().UnixNano(),
	}
	p.requested = true
	for _, port := range ports {
		if port != p.port {
			go func(port string) {
				conn, err := grpc.Dial(port, grpc.WithInsecure())
				if err != nil {
					log.Printf("Failed to connect to %s: %v", port, err)
					return
				}
				defer conn.Close()

				client := h4.NewH4Client(conn)
				_, err = client.SendMessage(context.Background(), req)
				if err != nil {
					log.Printf("Error sending request to %s: %v", port, err)
				} else {
					log.Printf("Request sent to peer on %s", port)
				}

			}(port)
		}
	}
}

// SetupNode sets up and starts the gRPC server for the peer
func (p *Peer) SetupNode() error {
	listener, err := net.Listen("tcp", p.port)
	if err != nil {
		log.Printf("Failed to listen on %s: %v", p.port, err)
		return err
	}

	log.Printf("Node %d starting gRPC server on %s", p.Id, p.port)

	go func() {
		if err := p.grpcServer.Serve(listener); err != nil {
			log.Printf("Failed to serve gRPC server on %s: %v", p.port, err)
		}
	}()

	return nil
}
