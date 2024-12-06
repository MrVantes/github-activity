package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"strings"
)

// Define structs for GitHub event data
type Commit struct {
	Sha    string `json:"sha"`
	Author struct {
		Email string `json:"email"`
		Name  string `json:"name"`
	} `json:"author"`
	Message string `json:"message"`
}

type Payload struct {
	Commits []Commit `json:"commits"`
}

type Repo struct {
	Name string `json:"name"`
	Url  string `json:"url"`
}

type Actor struct {
	Login string `json:"login"`
	Url   string `json:"url"`
}

type GitHubEvent struct {
	Id        string  `json:"id"`
	Type      string  `json:"type"`
	Actor     Actor   `json:"actor"`
	Repo      Repo    `json:"repo"`
	Payload   Payload `json:"payload"`
	Public    bool    `json:"public"`
	CreatedAt string  `json:"created_at"`
}

func userActivity(username string) ([]GitHubEvent, error) {
	url := fmt.Sprintf("https://api.github.com/users/%s/events", username)

	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var events []GitHubEvent
	err = json.NewDecoder(resp.Body).Decode(&events)
	if err != nil {
		return nil, err
	}

	return events, nil
}

func summarizeActivity(events []GitHubEvent) {
	eventSummary := map[string]int{
		"Push":        0,
		"PullRequest": 0,
		"Issue":       0,
		"Star":        0,
		"Fork":        0,
		"Create":      0,
		"Delete":      0,
		"Release":     0,
		"Wiki":        0,
	}

	for _, event := range events {
		switch event.Type {
		case "PushEvent":
			eventSummary["Push"] += len(event.Payload.Commits) // Count commits pushed
		case "PullRequestEvent":
			eventSummary["PullRequest"] += 1 // Count pull requests
		case "IssuesEvent":
			eventSummary["Issue"] += 1 // Count issues opened/closed
		case "WatchEvent":
			eventSummary["Star"] += 1 // Count starred repositories
		case "ForkEvent":
			eventSummary["Fork"] += 1 // Count forked repositories
		case "CreateEvent":
			eventSummary["Create"] += 1 // Count new branches or tags created
		case "DeleteEvent":
			eventSummary["Delete"] += 1 // Count branches or tags deleted
		case "ReleaseEvent":
			eventSummary["Release"] += 1 // Count releases created or published
		case "GollumEvent":
			eventSummary["Wiki"] += 1 // Count wiki page edits
		}
	}

	// Print the summarized activity
	if eventSummary["Push"] > 0 {
		fmt.Printf("- Pushed %d commits to %s\n", eventSummary["Push"], events[0].Repo.Name)
	}
	if eventSummary["PullRequest"] > 0 {
		fmt.Printf("- Opened a pull request in %s\n", events[0].Repo.Name)
	}
	if eventSummary["Issue"] > 0 {
		fmt.Printf("- Opened a new issue in %s\n", events[0].Repo.Name)
	}
	if eventSummary["Star"] > 0 {
		fmt.Printf("- Starred %s\n", events[0].Repo.Name)
	}
	if eventSummary["Fork"] > 0 {
		fmt.Printf("- Forked %s\n", events[0].Repo.Name)
	}
	if eventSummary["Create"] > 0 {
		fmt.Printf("- Created a new branch or tag in %s\n", events[0].Repo.Name)
	}
	if eventSummary["Delete"] > 0 {
		fmt.Printf("- Deleted a branch or tag in %s\n", events[0].Repo.Name)
	}
	if eventSummary["Release"] > 0 {
		fmt.Printf("- Created or published a release in %s\n", events[0].Repo.Name)
	}
	if eventSummary["Wiki"] > 0 {
		fmt.Printf("- Edited a wiki page in %s\n", events[0].Repo.Name)
	}
}

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Please provide a GitHub username.")
		os.Exit(1)
	}

	username := strings.Join(os.Args[1:], " ")

	events, err := userActivity(username)
	if err != nil {
		fmt.Println("Error:", err)
		os.Exit(1)
	}

	summarizeActivity(events)
}
