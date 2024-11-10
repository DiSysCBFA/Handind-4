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
}
