// server.go
package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"sync"
)

type Server struct {
	ID         int
	Port       string
	Peers      []string
	Votes      map[string]string
	VotesMutex sync.Mutex
	Faulty     bool // Simulate faulty behavior
}

func NewServer(id int, port string, peers []string, faulty bool) *Server {
	return &Server{
		ID:     id,
		Port:   port,
		Peers:  peers,
		Votes:  make(map[string]string),
		Faulty: faulty,
	}
}

func (s *Server) handleVote(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Only POST method is allowed", http.StatusMethodNotAllowed)
		return
	}

	var vote Vote
	err := json.NewDecoder(r.Body).Decode(&vote)
	if err != nil {
		http.Error(w, "Invalid vote data", http.StatusBadRequest)
		return
	}

	s.VotesMutex.Lock()
	s.Votes[vote.VoterID] = vote.Choice
	s.VotesMutex.Unlock()

	go s.broadcastVote(vote)

	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, "Vote received by server %d", s.ID)
}

func (s *Server) broadcastVote(vote Vote) {
	for _, peer := range s.Peers {
		url := fmt.Sprintf("http://%s/consensus", peer)
		data, _ := json.Marshal(vote)
		_, err := http.Post(url, "application/json", bytes.NewBuffer(data))
		if err != nil {
			log.Printf("Server %d failed to send vote to %s: %v", s.ID, peer, err)
		}
	}
}

func (s *Server) handleConsensus(w http.ResponseWriter, r *http.Request) {
	if s.Faulty {
		// Simulate faulty behavior
		http.Error(w, "Server is faulty", http.StatusInternalServerError)
		return
	}

	var vote Vote
	err := json.NewDecoder(r.Body).Decode(&vote)
	if err != nil {
		http.Error(w, "Invalid data", http.StatusBadRequest)
		return
	}

	s.VotesMutex.Lock()
	s.Votes[vote.VoterID] = vote.Choice
	s.VotesMutex.Unlock()

	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, "Consensus updated on server %d", s.ID)
}

func (s *Server) handleResults(w http.ResponseWriter, r *http.Request) {
	s.VotesMutex.Lock()
	defer s.VotesMutex.Unlock()

	results := make(map[string]int)
	for _, choice := range s.Votes {
		results[choice]++
	}

	json.NewEncoder(w).Encode(results)
}

func (s *Server) Start() {
	http.HandleFunc("/vote", s.handleVote)
	http.HandleFunc("/consensus", s.handleConsensus)
	http.HandleFunc("/results", s.handleResults)

	log.Printf("Server %d starting on port %s", s.ID, s.Port)
	log.Fatal(http.ListenAndServe(":"+s.Port, nil))
}
