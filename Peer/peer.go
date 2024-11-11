package peer

import (
	"context"
	"fmt"
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
	Id                  int
	port                string
	status              h4.Status
	grpcServer          *grpc.Server
	isInCriticalSection bool
	grantCount          int
	totalPeers          int
	mu                  sync.Mutex
}

// NewPeer creates a new Peer instance
func NewPeer(id int, port string, totalPeers int) *Peer {
	peer := &Peer{
		Id:                  id,
		port:                port,
		status:              h4.Status_DENIED,
		grpcServer:          grpc.NewServer(),
		isInCriticalSection: false,
		totalPeers:          totalPeers,
	}
	h4.RegisterH4Server(peer.grpcServer, peer)
	return peer
}

// Multicast sends a request message to all specified peer ports
func (p *Peer) Multicast(ports []string) {
	req := &h4.RequestMessage{
		Id:        int64(p.Id),
		Timestamp: time.Now().UnixNano(),
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
				_, err = client.Request(context.Background(), req)
				if err != nil {
					log.Printf("Error sending request to %s: %v", port, err)
				} else {
					log.Printf("Request sent to peer on %s", port)
				}
			}(port)
		}
	}
}

// Request handles incoming access requests from other nodes
func (p *Peer) Request(ctx context.Context, req *h4.RequestMessage) (*h4.ReplyMessage, error) {
	log.Printf("Received request from peer %d with timestamp %d", req.Id, req.Timestamp)

	// Decide whether to grant or deny access based on current state
	responseStatus := h4.Status_GRANTED
	if p.isInCriticalSection {
		responseStatus = h4.Status_DENIED
	}

	// Prepare the reply message
	reply := &h4.ReplyMessage{Id: int64(p.Id), Status: responseStatus}
	log.Printf("Sending reply to peer %d with status %v", req.Id, responseStatus)

	return reply, nil
}

// Reply handles incoming replies to previous requests
func (p *Peer) Reply(ctx context.Context, req *h4.RequestMessage) (*h4.ReplyMessage, error) {
	log.Printf("Received reply from peer %d", req.Id)

	// Temporarily bypass grantCount logic for testing
	if p.status == h4.Status_GRANTED {
		log.Printf("Grant acknowledged from peer %d", req.Id)
	}

	return &h4.ReplyMessage{Id: int64(p.Id), Status: h4.Status_GRANTED}, nil
}

// SendReply sends a reply (RequestMessage) to a specific peer with status embedded in internal logic
func (p *Peer) SendReply(requesterId int64, replyMessage *h4.RequestMessage, status h4.Status) error {
	requesterPort := fmt.Sprintf("localhost:%d", 4000+int(requesterId)) // Assuming peer IDs map to specific ports
	conn, err := grpc.Dial(requesterPort, grpc.WithInsecure())
	if err != nil {
		return fmt.Errorf("failed to connect to peer %d at %s: %v", requesterId, requesterPort, err)
	}
	defer conn.Close()

	client := h4.NewH4Client(conn)

	// Simulate sending a request message back as a reply, while maintaining the status internally
	if status == h4.Status_GRANTED {
		p.grantCount++
	}

	_, err = client.Reply(context.Background(), replyMessage)
	if err != nil {
		return fmt.Errorf("failed to send reply to peer %d: %v", requesterId, err)
	}
	return nil
}

// Access attempts to enter the critical section if the status is GRANTED
func (p *Peer) Access() {
	p.mu.Lock()
	defer p.mu.Unlock()

	if p.status == h4.Status_GRANTED && !p.isInCriticalSection {
		log.Printf("Node %d entering critical section", p.Id)
		p.isInCriticalSection = true

		// Simulate critical section work
		time.Sleep(2 * time.Second) // Replace with actual critical section logic if needed

		// Reset after exiting critical section
		p.isInCriticalSection = false
		p.status = h4.Status_DENIED
		p.grantCount = 0
		log.Printf("Node %d leaving critical section", p.Id)
	} else {
		log.Printf("Node %d cannot enter critical section, status: %v", p.Id, p.status)
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
