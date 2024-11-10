package peer

import (
	"context"
	"fmt"
	"log"
	"net"
	"time"

	critical "github.com/DiSysCBFA/Handind-4/Critical-Section"
	h4 "github.com/DiSysCBFA/Handind-4/H4"
	"google.golang.org/grpc"
)

type Peer struct {
	h4.UnimplementedH4Server
	Id         int
	port       string
	status     h4.Status
	requests   chan h4.RequestMessage
	grpcServer *grpc.Server
}

func NewPeer(id int, port string) *Peer {
	return &Peer{
		Id:         id,
		port:       port,
		status:     h4.Status_DENIED,
		requests:   make(chan h4.RequestMessage),
		grpcServer: grpc.NewServer(),
	}
}

// multicast sends a request message to all specified peer ports
func (p *Peer) multicast(ports []int) {
	req := h4.RequestMessage{Id: int64(p.Id), Timestamp: time.Now().UnixNano()}
	for _, port := range ports {
		address := fmt.Sprintf("localhost:%d", port)
		conn, err := grpc.Dial(address, grpc.WithInsecure())
		if err != nil {
			log.Printf("Failed to connect to %s: %v", address, err)
			continue
		}
		defer conn.Close()

		client := h4.NewH4Client(conn)
		_, err = client.Request(context.Background(), &req)
		if err != nil {
			log.Printf("Error sending request to %s: %v", address, err)
		} else {
			log.Printf("Request sent to peer on %s", address)
		}
	}
}

// access attempts to enter the critical section if the status is GRANTED
func (p *Peer) access() {
	if p.status == h4.Status_GRANTED {
		critical.Main()
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
	log.Println("Node setup on port: ", 4)

	go func() {
		if err := p.grpcServer.Serve(listener); err != nil {
			log.Printf("Failed to serve gRPC server on %s: %v", p.port, err)
		}
	}()

	return nil
}

// CreateNodeServer initializes a gRPC server and registers the Peer as the service handler
func CreateNodeServer(nodeID int, port string) (*grpc.Server, error) {
	peerServer := NewPeer(nodeID, port)
	grpcServer := grpc.NewServer()

	h4.RegisterH4Server(grpcServer, peerServer)

	log.Printf("Created gRPC server for node ID: %d on port %d", nodeID, port)

	return grpcServer, nil
}
