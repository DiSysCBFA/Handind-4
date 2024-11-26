package critical

import (
	"fmt"
	"log"
	"os"

	"github.com/manifoldco/promptui"
)

func Main() {
	fmt.Println("Critical section has been accessed")

	selectionQuit := promptui.Select{
		Label: "Do you want to quit?",
		Items: []string{"Yes", "No"},
	}

	_, result, err := selectionQuit.Run()
	if err != nil {
		log.Fatalf("Prompt failed: %v\n", err)
	}

	// Check the user's selection
	if result == "Yes" {
		log.Println("Exiting the program")
		os.Exit(0)
	} else if result == "No" {
		fmt.Println("Continuing the program")
	}
}
