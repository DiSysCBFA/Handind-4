package main

import (
	"bufio"
	"log"
	"os"
)

func main() {
	file, err := os.Open("ports.txt")
	if err != nil {
		log.Fatal(err)
		return
	}

	var NodeID int = 0

	r := bufio.NewReader(file)
	for {
		line, err := r.ReadString('\n')
		if err != nil {
			break
		}
		log.Println(line) // to be replaced with node attempt setup

		NodeID++

		//TODO: Implement node setup on port.
	}

	defer file.Close()
}
