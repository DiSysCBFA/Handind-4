package peer

import (
	"fmt"
	critical "github.com/DiSysCBFA/Handind-4/Critical-section"
	"github.com/DiSysCBFA/Handind-4/h4"
	"google.golang.org/grpc"
	_ "google.golang.org/grpc/status"
	"log"
	"time"
)

type Peer struct {
	id         int
	port       int
	status     h4.Status
	requests   chan h4.RequestMessage
	grpcServer *grpc.Server
	grpcClient *grpc.ClientConn
}

func NewPeer(id int, port int) *Peer {
	return &Peer{
		id:         id,
		port:       port,
		status:     h4.Status_DENIED,
		requests:   make(chan h4.RequestMessage),
		grpcServer: grpc.NewServer(),
	}
}

func (p *Peer) multicast() {
	req := h4.RequestMessage{Id: int64(p.id), Timestamp: time.Now().UnixNano()}
	for _, port := range ports {
		conn, err := grpc.Dial(fmt.Sprintf("localhost:%d", peer.port), grpc.WithInsecure())
		if err != nil {
			log.Println(err)
			continue

	}


func (p *Peer) access() {
	if p.status == h4.Status_GRANTED {
		critical.Main()
	}


func (p *Peer) SetupNode(port string) error {
	log.Println("Setting up node on port:", port)

	listener, err := net.Listen("tcp", port)
	if err != nil {
		log.Printf("Failed to listen on port %s: %v", port, err)
		return err
	}

	// Create the gRPC server instance
	grpcServer, err := CreateNodeServer(p.NodeID, port)
	if err != nil {
		log.Printf("Failed to create node server: %v", err)
		return err
	}

	// Start serving gRPC requests
	log.Printf("Node %d starting gRPC server on %s", p.NodeID, port)
	if err := grpcServer.Serve(listener); err != nil {
		log.Printf("Failed to serve gRPC server on %s: %v", port, err)
		return err
	}

	return nil
}

func CreateNodeServer(nodeID int, port string) (*grpc.Server, error) {
	grpcServer := grpc.NewServer()
	peerServer, err := CreateServerForPeer(nodeID, port)
	if err != nil {
		return nil, err
	}

	h4.RegisterH4Server(grpcServer, peerServer)

	log.Printf("Created gRPC server for node ID: %d on port %s", nodeID, port)

	return grpcServer, nil
}

func CreateServerForPeer(nodeID int, port string) (*Peer, error) {
	peerServer := &Peer{
		NodeID: nodeID,
		Port:   port,
	}
	return peerServer, nil
}
