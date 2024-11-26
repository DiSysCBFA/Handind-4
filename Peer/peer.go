package peer

import (
	"context"
	critical "github.com/DiSysCBFA/Handind-4/Critical"
	h4 "github.com/DiSysCBFA/Handind-4/h4"
	"google.golang.org/grpc"
	"log"
	"net"
	"sync"
	"time"
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

// Request ensures that timestamp remains the same throughout the request
func (p *Peer) Request(ports []string) {
	p.Reqtimestamp = time.Now().UnixNano()
	p.Multicast(ports)
}

// Multicast sends a request message to all specified peer ports
func (p *Peer) Multicast(ports []string) {
	p.requested = true
	req := &h4.Message{
		Timestamp: p.Reqtimestamp,
		Answer:    0,
	}

	respChan := make(chan int32, len(ports))
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
				resp, err := client.SendMessage(context.Background(), req)
				if err != nil {
					log.Printf("Error sending request to %s: %v", port, err)
				} else {
					log.Printf("Request sent to peer on %s", port)
					respChan <- resp.Answer
				}
			}(port)
		}
	}

	// logic for the response
	go func() {
		var receivedTwos bool
		responses := 0

		for range ports {
			select {
			case answer := <-respChan:
				responses++
				if answer == 2 {
					receivedTwos = true
				}
			case <-time.After(3 * time.Second): //  To avoid timeout
				log.Println("Timeout waiting for responses")
				break
			}
		}

		if receivedTwos {
			log.Println("Received at least one 2(no), waiting 5 seconds and retrying...")
			time.Sleep(5 * time.Second)
			p.Multicast(ports)
		} else if responses == len(ports)-1 {
			log.Println("Received only 1's(yes), accesing critical section")
			critical.Main()
			p.requested = false
		}
	}()
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
