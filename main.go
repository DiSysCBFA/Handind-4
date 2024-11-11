package main

import (
	"bufio"
	"log"
	"os"
	"strings"

	"github.com/DiSysCBFA/Handind-4/peer"
	"github.com/manifoldco/promptui"
)

func main() {
	file, err := os.Open("ports.txt")
	if err != nil {
		log.Fatal(err)
		return
	}

	var peerNode *peer.Peer // Declare as a pointer

	var NodeID int = 0

	r := bufio.NewReader(file)
	for {
		port, err := r.ReadString('\n') //! Make sure last line has a new line
		if err != nil {
			break
		}

		port = strings.TrimSpace(port)
		log.Println("Initializing peer with port:", port) // Log port initialization

		NodeID += 1
		peerNode = peer.NewPeer(NodeID, port) // Assign to the outer peerNode

		err = peerNode.SetupNode()

		if err == nil {
			log.Printf("Node %d setup successfully on port %s", NodeID, port)
			break
		} else {
			log.Printf("Failed to set up node on port %s", port)
		}
	}

	if peerNode == nil {
		log.Fatal("No peer node was initialized")
		return
	}

	selection := promptui.Select{
		Label: "Select action",
		Items: []string{"Request", "Exit"},
	}

	_, result, err := selection.Run()
	if err != nil {
		log.Fatalf("Failed to run: %v", err)
	}

	if result == "Request" {
		log.Println("Requesting access to critical section")
		peerNode.Multicast([]string{"localhost:4000", "localhost:4001", "localhost:4002"})
	} else {
		os.Exit(0)
	}

	defer file.Close()
}
