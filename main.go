package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
)

func userActivity(username string) {

	url := fmt.Sprintf("https://api.github.com/users/%s/events", username)

	resp, err := http.Get(url)
	if err != nil {
		log.Fatal(err)
	}
	defer resp.Body.Close()
}

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Please provide a command.")
		os.Exit(1)
	}

	username := os.Args[1]

	output, err := userActivity(username)
	if err != nil {
		fmt.Println("Error:", err)
		os.Exit(1)
	}

	fmt.Println(output)
}
