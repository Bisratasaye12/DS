// package main

// import (
// 	"fmt"
// 	"math"
// 	"net"
// 	"os"
// 	"strings"
// 	"sync"
// )

// type Replica struct {
// 	value float64
// 	mu    sync.Mutex
// 	peers []string // List of peer replica addresses
// }

// // Update updates the replica value if the new value is within the delta threshold
// func (r *Replica) Update(newValue, delta float64) bool {
// 	r.mu.Lock()
// 	defer r.mu.Unlock()
// 	if math.Abs(newValue-r.value) <= delta {
// 		r.value = newValue
// 		return true
// 	}
// 	return false
// }

// // propagateUpdates sends the updated value to all peers
// func (r *Replica) propagateUpdates() {
// 	for _, peer := range r.peers {
// 		go func(peer string) {
// 			conn, err := net.Dial("tcp", peer)
// 			if err != nil {
// 				fmt.Println("Error connecting to peer:", peer, err)
// 				return
// 			}
// 			defer conn.Close()

// 			r.mu.Lock()
// 			message := fmt.Sprintf("%.2f\n", r.value)
// 			r.mu.Unlock()
// 			conn.Write([]byte(message))
// 		}(peer)
// 	}
// }

// // handleConnection handles incoming connections from peers
// func handleConnection(conn net.Conn, replica *Replica, delta float64) {
// 	defer conn.Close()
// 	for {
// 		buffer := make([]byte, 1024)
// 		n, err := conn.Read(buffer)
// 		if err != nil {
// 			break
// 		}
// 		newValue := strings.TrimSpace(string(buffer[:n]))
// 		var value float64
// 		fmt.Sscanf(newValue, "%f", &value)
// 		replica.Update(value, delta)
// 	}
// }

// func main() {

// 	if len(os.Args) < 4 {
// 		fmt.Println("Usage: go run replica_numerical.go <machine_ip:port> <delta> <peer1_ip:port> [<peer2_ip:port>...]")
// 		return
// 	}

// 	// Parse command-line arguments
// 	machineAddr := os.Args[1]
// 	delta := 0.0
// 	fmt.Sscanf(os.Args[2], "%f", &delta)
// 	peers := os.Args[3:]

// 	// Initialize the replica
// 	replica := &Replica{
// 		value: 10.0,
// 		peers: peers,
// 	}

// 	// Start the TCP server
// 	listener, err := net.Listen("tcp", machineAddr)
// 	if err != nil {
// 		panic(err)
// 	}
// 	defer listener.Close()

// 	fmt.Printf("Replica listening on %s\n", machineAddr)

// 	go func() {
// 		for {
// 			conn, err := listener.Accept()
// 			if err != nil {
// 				continue
// 			}
// 			go handleConnection(conn, replica, delta)
// 		}
// 	}()

// 	// Simulate an update
// 	replica.value = 12.0
// 	replica.propagateUpdates()

// 	fmt.Println("Replica Value:", replica.value)
// }

package main

import (
	"fmt"
	"math"
	"net"
	"os"
	"strings"
	"sync"
)

type Replica struct {
	value float64
	mu    sync.Mutex
	peers []string // List of peer replica addresses
}

func (r *Replica) Update(newValue, delta float64) bool {
	r.mu.Lock()
	defer r.mu.Unlock()
	if math.Abs(newValue-r.value) <= delta {
		r.value = newValue
		return true
	}
	return false
}

func (r *Replica) propagateUpdates(delta float64) {
	for _, peer := range r.peers {
		go func(peer string) {
			conn, err := net.Dial("tcp", peer)
			if err != nil {
				fmt.Println("Error connecting to peer:", peer, err)
				return
			}
			defer conn.Close()

			r.mu.Lock()
			message := fmt.Sprintf("%.2f\n", r.value)
			r.mu.Unlock()
			conn.Write([]byte(message))
		}(peer)
	}
}

func handleConnection(conn net.Conn, replica *Replica, delta float64) {
	defer conn.Close()
	for {
		buffer := make([]byte, 1024)
		n, err := conn.Read(buffer)
		if err != nil {
			break
		}
		newValue := strings.TrimSpace(string(buffer[:n]))
		var value float64
		fmt.Sscanf(newValue, "%f", &value)
		replica.Update(value, delta)
	}
}

func main() {
	if len(os.Args) < 3 {
		fmt.Println("Usage: go run replica_numerical.go <machine_ip:port> <peer1_ip:port> [<peer2_ip:port>...]")
		return
	}

	address := os.Args[1]
	peers := os.Args[2:]

	replica := &Replica{
		value: 10.0,
		peers: peers,
	}

	delta := 5.0
	listener, err := net.Listen("tcp", address)
	if err != nil {
		panic(err)
	}
	defer listener.Close()

	fmt.Printf("Replica listening on %s\n", address)

	// Handle incoming connections
	go func() {
		for {
			conn, err := listener.Accept()
			if err != nil {
				continue
			}
			go handleConnection(conn, replica, delta)
		}
	}()

	// Simulate an update
	replica.value = 12.0
	replica.propagateUpdates(delta)

	fmt.Println("Replica Value:", replica.value)
}
