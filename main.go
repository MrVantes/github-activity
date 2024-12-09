package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
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
	repoEventCounts := make(map[string]map[string]int)

	for _, event := range events {
		repoName := event.Repo.Name
		if _, exists := repoEventCounts[repoName]; !exists {
			repoEventCounts[repoName] = make(map[string]int)
		}

		switch event.Type {
		case "PushEvent":
			repoEventCounts[repoName]["Push"] += len(event.Payload.Commits)
		case "PullRequestEvent":
			repoEventCounts[repoName]["Pull"] += 1
		case "IssuesEvent":
			repoEventCounts[repoName]["Issue"] += 1
		case "WatchEvent":
			repoEventCounts[repoName]["Star"] += 1
		case "ForkEvent":
			repoEventCounts[repoName]["Fork"] += 1
		case "CreateEvent":
			repoEventCounts[repoName]["Create"] += 1
		case "DeleteEvent":
			repoEventCounts[repoName]["Delete"] += 1
		case "ReleaseEvent":
			repoEventCounts[repoName]["Release"] += 1
		case "GollumEvent":
			repoEventCounts[repoName]["Wiki"] += 1
		}
	}

	for repo, eventCounts := range repoEventCounts {
		if pushCount, exists := eventCounts["Push"]; exists && pushCount > 0 {
			fmt.Printf("- Pushed %d commits to %s\n", pushCount, repo)
		}
		if pullCount, exists := eventCounts["Pull"]; exists && pullCount > 0 {
			fmt.Printf("- Pulled %d from %s\n", pullCount, repo)
		}
		if issueCount, exists := eventCounts["Issue"]; exists && issueCount > 0 {
			fmt.Printf("- Opened %d issue/s in %s\n", issueCount, repo)
		}
		if starCount, exists := eventCounts["Star"]; exists && starCount > 0 {
			fmt.Printf("- Starred %s\n", repo)
		}
		if forkCount, exists := eventCounts["Fork"]; exists && forkCount > 0 {
			fmt.Printf("- Forked %s\n", repo)
		}
		if createCount, exists := eventCounts["Create"]; exists && createCount > 0 {
			fmt.Printf("- Created %s\n", repo)
		}
		if deleteCount, exists := eventCounts["Delete"]; exists && deleteCount > 0 {
			fmt.Printf("- Deleted %s\n", repo)
		}
		if releaseCount, exists := eventCounts["Release"]; exists && releaseCount > 0 {
			fmt.Printf("- Released %s\n", repo)
		}
		if wikiCount, exists := eventCounts["Wiki"]; exists && wikiCount > 0 {
			fmt.Printf("- Wiki-ed %s\n", repo)
		}
	}
}

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Please provide a GitHub username.")
		os.Exit(1)
	}

	username := os.Args[1]

	events, err := userActivity(username)
	if err != nil {
		fmt.Println("Error: username does not exist")
	}

	summarizeActivity(events)
}
