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
	h4.H4Client
	Id           int
	port         string
	requested    bool
	Reqtimestamp int64
	grpcServer   *grpc.Server
	mu           sync.Mutex
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
	p.Reqtimestamp = time.Now().UnixNano()
	p.requested = true
	req := &h4.Message{
		Timestamp: p.Reqtimestamp,
		Answer:    0,
	}
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

func (p *Peer) logic() {
	// Logic for the peer

}
func (p *Peer) SendMessage(ctx context.Context, req *h4.Message) (*h4.Message, error) {
	resp := &h4.Message{}
	log.Println("Received request with timestamp: ", req.Timestamp)
	if !p.requested {
		resp = &h4.Message{
			Answer: 1,
		}
		log.Println(resp)
	} else {
		if p.Reqtimestamp > req.Timestamp {
			resp = &h4.Message{
				Answer: 1,
			}
			log.Println(resp)
		} else {
			resp = &h4.Message{
				Answer: 2,
			}
			log.Println(resp)
		}
	}
	return resp, nil
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
