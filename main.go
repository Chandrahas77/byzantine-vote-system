// main.go
package main

import (
    "flag"
)

func main() {
    var id int
    var port string
    var faulty bool
    flag.IntVar(&id, "id", 1, "Server ID")
    flag.StringVar(&port, "port", "8000", "Port to listen on")
    flag.BoolVar(&faulty, "faulty", false, "Is the server faulty?")
    flag.Parse()

    peers := []string{"localhost:8000", "localhost:8001", "localhost:8002"}
    // Remove self from peers
    peersWithoutSelf := make([]string, 0)
    for _, peer := range peers {
        if peer != "localhost:"+port {
            peersWithoutSelf = append(peersWithoutSelf, peer)
        }
    }

    server := NewServer(id, port, peersWithoutSelf, faulty)
    server.Start()
}
