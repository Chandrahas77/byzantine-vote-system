package main

import (
	"encoding/json"
	"log"
	"math/rand"
	"net/http"
	"sync"
)

type Vote struct {
	VoterID string `json:"voter_id"`
	Choice  string `json:"choice"`
}

type Server struct {
	ID         int
	VoteStore  map[string]string
	mutex      sync.Mutex
	isFaulty   bool
	httpServer *http.Server
}

func NewServer(id int, isFaulty bool) *Server {
	return &Server{
		ID:        id,
		VoteStore: make(map[string]string),
		isFaulty:  isFaulty,
	}
}

func (s *Server) CastVote(w http.ResponseWriter, r *http.Request) {
	var vote Vote
	if err := json.NewDecoder(r.Body).Decode(&vote); err != nil {
		http.Error(w, "Invalid request", http.StatusBadRequest)
		return
	}

	s.mutex.Lock()
	defer s.mutex.Unlock()

	if s.isFaulty {
		// Malicious behavior: Flip vote choice randomly
		log.Printf("Server %d (Faulty): Received vote from %s, altering vote choice!", s.ID, vote.VoterID)
		vote.Choice = []string{"A", "B", "C"}[rand.Intn(3)]
	} else {
		log.Printf("Server %d: Received vote from %s for choice %s", s.ID, vote.VoterID, vote.Choice)
	}
	s.VoteStore[vote.VoterID] = vote.Choice
	w.WriteHeader(http.StatusOK)
}

func (s *Server) GetVotes(w http.ResponseWriter, r *http.Request) {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	if s.isFaulty {
		// Malicious behavior: Provide incorrect votes
		log.Printf("Server %d (Faulty): Returning altered vote results!", s.ID)
		faultyVotes := map[string]string{}
		for voter := range s.VoteStore {
			faultyVotes[voter] = []string{"A", "B", "C"}[rand.Intn(3)]
		}
		json.NewEncoder(w).Encode(faultyVotes)
		return
	}

	log.Printf("Server %d: Returning correct vote results", s.ID)
	json.NewEncoder(w).Encode(s.VoteStore)
}

func (s *Server) StartServer(port string) {
	mux := http.NewServeMux()
	mux.HandleFunc("/cast", s.CastVote)
	mux.HandleFunc("/votes", s.GetVotes)

	s.httpServer = &http.Server{
		Addr:    ":" + port,
		Handler: mux,
	}

	log.Printf("Server %d starting on port %s (Faulty: %v)", s.ID, port, s.isFaulty)
	log.Fatal(s.httpServer.ListenAndServe())
}
