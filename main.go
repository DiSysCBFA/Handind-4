package main

import (
	"bufio"
	"log"
	"os"

	"github.com/DiSysCBFA/Handind-4/peer"
)

func main() {
	file, err := os.Open("ports.txt")
	if err != nil {
		log.Fatal(err)
		return
	}

	peerNode := peer.Peer{}

	var NodeID int = 0

	r := bufio.NewReader(file)
	for {
		port, err := r.ReadString('\n') //! Make sure last line has a new line
		if err != nil {
			break
		}
		log.Println(port) //TODO:  to be replaced with node attempt setup

		NodeID++

		err = peerNode.SetupNode(port)
		if err == nil {
			break
		}

		// TODO: Implement node setup on port.
	}

	defer file.Close()
}
