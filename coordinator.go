package main

import (
	"encoding/json"
	"log"
	"net/http"
	"sync"
)

type Coordinator struct {
	Servers []string
}

func NewCoordinator(servers []string) *Coordinator {
	return &Coordinator{Servers: servers}
}

func (c *Coordinator) CollectVotes() map[string]int {
	voteCounts := make(map[string]int)
	voteResponses := make([]map[string]string, len(c.Servers))
	var wg sync.WaitGroup

	log.Println("Coordinator: Collecting votes from all servers...")
	for i, server := range c.Servers {
		wg.Add(1)
		go func(idx int, url string) {
			defer wg.Done()
			resp, err := http.Get(url + "/votes")
			if err != nil {
				log.Printf("Coordinator: Error fetching votes from %s: %v", url, err)
				return
			}
			defer resp.Body.Close()

			var votes map[string]string
			if err := json.NewDecoder(resp.Body).Decode(&votes); err != nil {
				log.Printf("Coordinator: Error decoding votes from %s: %v", url, err)
				return
			}
			voteResponses[idx] = votes
			log.Printf("Coordinator: Received votes from %s: %v", url, votes)
		}(i, server)
	}

	wg.Wait()

	for _, response := range voteResponses {
		for _, choice := range response {
			voteCounts[choice]++
		}
	}

	log.Printf("Coordinator: Aggregated vote counts: %v", voteCounts)
	return voteCounts
}

func (c *Coordinator) FindConsensus(voteCounts map[string]int) string {
	var winner string
	maxVotes := 0

	for choice, count := range voteCounts {
		if count > maxVotes {
			maxVotes = count
			winner = choice
		}
	}

	log.Printf("Coordinator: Final consensus - Winner is %s with %d votes", winner, maxVotes)
	return winner
}
