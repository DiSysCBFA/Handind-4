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

	r := bufio.NewReader(file)
	for {
		line, err := r.ReadString('\n')
		if err != nil {
			break
		}
		log.Println(line)
	}

	defer file.Close()
}
