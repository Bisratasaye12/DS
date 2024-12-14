package main

import (
	"fmt"
	"log"
	"net"
	"net/rpc"
	"os"
	"sync"
)

type Replica struct {
	data    map[string]string
	mu      sync.Mutex
	peers   []string // List of peer addresses
	ackLock sync.Mutex
	acks    map[string]int // Track acknowledgments
}

type Args struct {
	Key   string
	Value string
}

// Update method updates the replica's data and replies true upon success
func (r *Replica) Update(args *Args, reply *bool) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.data[args.Key] = args.Value
	log.Printf("[Update] Data updated: {%s: %s}\n", args.Key, args.Value)
	*reply = true
	return nil
}

// propagateUpdates sends the update to all peers and tracks acknowledgments
func (r *Replica) propagateUpdates(key, value string) {
	log.Printf("[Propagation] Propagating update {%s: %s} to peers\n", key, value)
	r.ackLock.Lock()
	r.acks[key] = 0
	r.ackLock.Unlock()

	for _, peer := range r.peers {
		go func(peer string) {
			client, err := rpc.Dial("tcp", peer)
			if err != nil {
				log.Printf("[Error] Unable to connect to peer %s: %v\n", peer, err)
				return
			}
			defer client.Close()

			args := &Args{Key: key, Value: value}
			var reply bool
			if err := client.Call("Replica.Update", args, &reply); err == nil && reply {
				r.ackLock.Lock()
				r.acks[key]++
				r.ackLock.Unlock()
				log.Printf("[Acknowledgment] Received acknowledgment from peer %s for {%s: %s}\n", peer, key, value)
			} else {
				log.Printf("[Error] Failed to get acknowledgment from peer %s: %v\n", peer, err)
			}
		}(peer)
	}
}

// waitForAcknowledgments waits until a quorum of peers have acknowledged the update
func (r *Replica) waitForAcknowledgments(key string, quorum int) {
	log.Printf("[Acknowledgment] Waiting for quorum of %d acknowledgments for key: %s\n", quorum, key)
	for {
		r.ackLock.Lock()
		if r.acks[key] >= quorum {
			r.ackLock.Unlock()
			break
		}
		r.ackLock.Unlock()
	}
	log.Printf("[Commit] Update committed for key: %s after receiving %d acknowledgments\n", key, quorum)
}

func main() {
	if len(os.Args) < 4 {
		fmt.Println("Usage: go run replica_strong.go <machine_ip:port> <quorum> <peer1_ip:port> [<peer2_ip:port>...]")
		return
	}

	// Parse command-line arguments
	machineAddr := os.Args[1]
	peers := os.Args[3:]
	quorum := len(peers)/2 + 1

	// Initialize the replica
	replica := &Replica{
		data:  make(map[string]string),
		peers: peers,
		acks:  make(map[string]int),
	}

	rpc.Register(replica)

	// Start the RPC server
	listener, err := net.Listen("tcp", machineAddr)
	if err != nil {
		log.Fatalf("[Error] Failed to start server on %s: %v\n", machineAddr, err)
	}
	defer listener.Close()
	log.Printf("[Server] Replica RPC Server listening on %s\n", machineAddr)

	// Start accepting connections
	go func() {
		for {
			conn, err := listener.Accept()
			if err != nil {
				continue
			}
			go rpc.ServeConn(conn)
		}
	}()

	// Simulate a strong consistency update
	key, value := "key1", "value1"
	log.Printf("[Simulation] Initiating update {%s: %s}\n", key, value)
	var reply bool
	replica.Update(&Args{Key: key, Value: value}, &reply)
	replica.propagateUpdates(key, value)
	replica.waitForAcknowledgments(key, quorum)

	log.Println("[Complete] Update committed successfully")
}
