package main

import (
    "bytes"
    "encoding/json"
    "fmt"
    "net/http"
    "time"
)

// Proposal represents the data structure sent in the POST request.
type Proposal struct {
    ProposalNumber int    `json:"ProposalNumber"`
    Value          string `json:"Value"`
}

func main() {
    // List of server endpoints (all resolved through the Kubernetes service)
    servers := []string{
        "http://localhost:8080/propose", // Using localhost for testing with port forwarding
    }

    // Initial Proposal data
    proposal := Proposal{
        ProposalNumber: 5, 
        Value:          "Hello from Paxos client!",
    }

    // Serialize proposal to JSON
    data, err := json.Marshal(proposal)
    if err != nil {
        fmt.Printf("Error marshaling proposal: %v\n", err)
        return
    }

    // Send proposals to all servers
    for _, server := range servers {
        go sendProposal(server, data)
    }

    // Wait to ensure all goroutines complete (for simplicity)
    time.Sleep(5 * time.Second)
}

func sendProposal(server string, data []byte) {
    // Retry logic with ProposalNumber increment
    for i := 1; i <= 3; i++ { // Retry up to 3 times
        fmt.Printf("Attempt %d: Sending to %s...\n", i, server)

        resp, err := http.Post(server, "application/json", bytes.NewBuffer(data))
        if err != nil {
            fmt.Printf("Attempt %d: Error sending to %s: %v\n", i, server, err)
            time.Sleep(2 * time.Second) // Wait before retrying
            continue
        }

        // Print response
        defer resp.Body.Close()
        if resp.StatusCode == http.StatusOK {
            fmt.Printf("Success: Received response from %s\n", server)
            break
        } else if resp.StatusCode == http.StatusConflict {
            fmt.Printf("Attempt %d: Conflict error from %s. Incrementing ProposalNumber...\n", i, server)
            // Increment ProposalNumber and resend
            var proposal Proposal
            json.Unmarshal(data, &proposal)
            proposal.ProposalNumber++
            data, _ = json.Marshal(proposal) // Update data with new ProposalNumber
        } else {
            fmt.Printf("Failed: Received status code %d from %s\n", resp.StatusCode, server)
            break
        }
    }
}
