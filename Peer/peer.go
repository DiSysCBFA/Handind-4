package peer

import (
	"context"
	"log"
	"net"
	"time"

	critical "github.com/DiSysCBFA/Handind-4/critical"
	h4 "github.com/DiSysCBFA/Handind-4/h4"
	"google.golang.org/grpc"
)

type Peer struct {
	h4.UnimplementedH4Server
	Id         int
	port       string
	status     h4.Status
	grpcServer *grpc.Server
}

func NewPeer(id int, port string) *Peer {
	peer := &Peer{
		Id:         id,
		port:       port,
		status:     h4.Status_DENIED,
		grpcServer: grpc.NewServer(),
	}
	h4.RegisterH4Server(peer.grpcServer, peer) // Register the peer as the service handler
	return peer
}

// Multicast sends a request message to all specified peer ports
// Multicast sends a request message to all specified peer ports
func (p *Peer) Multicast(ports []string) {
	req := h4.RequestMessage{Id: int64(p.Id), Timestamp: time.Now().UnixNano()}
	for _, port := range ports {
		if port != p.port {
			conn, err := grpc.Dial(port, grpc.WithInsecure())
			if err != nil {
				log.Printf("Failed to connect to %s: %v", port, err)
				continue
			}
			defer conn.Close()

			client := h4.NewH4Client(conn)
			// Send the request
			_, err = client.Request(context.Background(), &req)
			if err != nil {
				log.Printf("Error sending request to %s: %v", port, err)
			} else {
				log.Printf("Request sent to peer on %s", port)
			}
		}
	}

	// Log the current status of the peer
	log.Printf("Peer %d status: %v", p.Id, p.status)

	// Try to access the critical section if the status is granted
	p.access()
}

// access attempts to enter the critical section if the status is GRANTED
// access attempts to enter the critical section if the status is GRANTED
func (p *Peer) access() {
	if p.status == h4.Status_GRANTED {
		critical.Main() // Enter critical section
	} else {
		log.Printf("Peer %d not granted access. Current status: %v", p.Id, p.status)
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

func (p *Peer) Request(ctx context.Context, req *h4.RequestMessage) (*h4.ReplyMessage, error) {
	// Log the received request for debugging
	log.Printf("Received request from peer %d with timestamp %d", req.Id, req.Timestamp)

	// Use timestamp to decide whether to grant access
	// You can further refine this logic, for example, by adding priority if necessary
	if p.status == h4.Status_DENIED || req.Timestamp < time.Now().UnixNano() {
		// Grant access only if the timestamp is earlier or the peer is in DENIED status
		p.status = h4.Status_GRANTED
		log.Printf("Granting access to peer %d", req.Id)
	} else {
		log.Printf("Denying access to peer %d", req.Id)
	}

	// Call access() to try entering the critical section
	p.access()

	// Respond with the current status of the peer (granted or denied)
	response := &h4.ReplyMessage{
		Status: p.status, // Send back the current status
	}
	return response, nil
}
