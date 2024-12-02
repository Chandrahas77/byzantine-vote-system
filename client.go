package main

import (
	"bytes"
	"encoding/json"
	"log"
	"net/http"
)

type Client struct {
	ServerURLs []string
}

func (c *Client) SendVote(voterID, choice string) {
	vote := Vote{VoterID: voterID, Choice: choice}
	data, _ := json.Marshal(vote)

	for _, serverURL := range c.ServerURLs {
		resp, err := http.Post(serverURL+"/cast", "application/json", bytes.NewBuffer(data))
		if err != nil {
			log.Printf("Error sending vote to %s: %v", serverURL, err)
			continue
		}
		resp.Body.Close()
	}
}
