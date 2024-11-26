package main

import (
	"bufio"
	"log"
	"os"
	"strconv"
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

	nop, err := r.ReadString('\n')
	var numberOfPeers, _ = strconv.Atoi(strings.TrimSpace(nop))
	log.Println("Number of peers is: ", numberOfPeers)
	if err != nil {
		log.Fatal(err)
		return
	}
	for i := 0; i < numberOfPeers; i++ {
		port, err := r.ReadString('\n') //! Make sure last line has a new line
		log.Println(port)
		if err != nil {
			break
		}

		port = strings.TrimSpace(port)
		log.Println(port) // Displaying each port read

		NodeID += 1
		peerNode = peer.NewPeer(NodeID, port, numberOfPeers) // Assign to the outer peerNode

		err = peerNode.SetupNode()

		if err == nil {
			break
		}

		// Node setup on port completed successfully
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
		peerNode.Request([]string{"localhost:4000", "localhost:4001", "localhost:4002"})
	} else {
		os.Exit(0)
	}

	defer file.Close()
	select {}
}
