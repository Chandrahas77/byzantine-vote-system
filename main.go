package main

import (
	"log"
	"time"
)

func main() {
	// Start servers
	servers := []string{"http://localhost:8000", "http://localhost:8001", "http://localhost:8002"}
	go func() {
		NewServer(1, false).StartServer("8000")
	}()
	go func() {
		NewServer(2, true).StartServer("8001")
	}()
	go func() {
		NewServer(3, false).StartServer("8002")
	}()

	time.Sleep(2 * time.Second) // Wait for servers to start

	// Cast votes
	client := &Client{ServerURLs: servers}
	client.SendVote("voter1", "A")
	client.SendVote("voter2", "B")
	client.SendVote("voter3", "A")

	time.Sleep(1 * time.Second) // Allow time for votes to propagate

	// Collect results
	coordinator := NewCoordinator(servers)
	voteCounts := coordinator.CollectVotes()
	winner := coordinator.FindConsensus(voteCounts)

	log.Printf("Election Winner: %s", winner)

	// Keep main running
	select {}
}
