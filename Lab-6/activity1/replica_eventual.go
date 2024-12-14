package main

import (
	"bufio"
	"fmt"
	"math/rand"
	"net"
	"os"
	"strings"
	"sync"
	"time"
)

type Replica struct {
	data  map[string]string
	mu    sync.Mutex
	peers []string // List of peer addresses
}

func (r *Replica) Update(key, value string) {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.data[key] = value
}

func (r *Replica) propagateUpdates(key, value string) {
	for _, peer := range r.peers {
		go func(peer string) {
			time.Sleep(time.Duration(rand.Intn(1000)) * time.Millisecond) // Simulating delay Exercise 1 Question 1
			conn, err := net.Dial("tcp", peer)
			if err != nil {
				fmt.Println("Error connecting to peer:", peer, err)
				return
			}
			defer conn.Close()
			fmt.Fprintf(conn, "%s:%s\n", key, value)
		}(peer)
	}
}
func handleConnection(conn net.Conn, replica *Replica) {
	defer conn.Close()
	reader := bufio.NewReader(conn)
	for {
		message, err := reader.ReadString('\n')
		if err != nil {
			break
		}
		parts := strings.Split(strings.TrimSpace(message), ":")
		if len(parts) == 2 {
			replica.Update(parts[0], parts[1])
		}
	}
}
func main() {
	if len(os.Args) < 3 {
		fmt.Println("Usage: go run replica_eventual.go <machine_ip:port> <peer1_ip:port> [<peer2_ip:port>...]")

		return
	}
	// Parse command-line arguments
	machineAddr := os.Args[1]
	peers := os.Args[2:]
	// Initialize the replica
	replica := &Replica{

		data:  make(map[string]string),
		peers: peers,
	}
	// Start the server
	listener, err := net.Listen("tcp", machineAddr)
	if err != nil {
		panic(err)
	}
	defer listener.Close()
	fmt.Printf("Replica listening on %s\n", machineAddr)
	go func() {
		for {
			conn, err := listener.Accept()
			if err != nil {
				continue
			}
			go handleConnection(conn, replica)
		}
	}()
	// Simulate an update
	start := time.Now() // Exercise 1 Question 2
	replica.Update("key2", "value2")
	replica.propagateUpdates("key2", "value2")

	time.Sleep(5 * time.Second) // Allow propagation
	replica.mu.Lock()
	fmt.Println("Final Replica Data:", replica.data)
	fmt.Printf("Convergence Time: %v\n", time.Since(start)) // Exercise 1 Question 2
	replica.mu.Unlock()

}
