# Paper Storage Server

A distributed system for storing, retrieving, and managing academic papers using RPC (Remote Procedure Calls) and RabbitMQ for message passing notifications.

## Features

- **Store Papers**: Upload academic papers with metadata such as author and title.
- **Retrieve Papers**: Fetch metadata or content of stored papers by their ID.
- **List Papers**: View all stored papers in the system.
- **Paper Notifications**: Get real-time notifications when a new paper is added, powered by RabbitMQ.
- **Distributed Architecture**: Implements a client-server model using gRPC and message-passing mechanisms.

---

## Project Structure

```plaintext
.
├── README.md
├── TS.pdf
├── cmd
│   ├── paperclient
│   │   ├── clientNotification
│   │   │   └── client_rabbitMQ.go
│   │   └── main.go
│   └── paperserver
│       └── main.go
├── go.mod
├── go.sum
├── paper-storage-server
│   └── paperpb
│       ├── paperStorageServer.pb.go
│       └── paperStorageServer_grpc.pb.go
├── proto
│   └── paperStorageServer.proto
├── screenshots
│   ├── add_paper_notification_multiclient.png
│   ├── add_paper_single_client.png
│   ├── fetch_content.png
│   ├── fetch_content_pre.png
│   ├── get_detail_0.png
│   ├── get_detail_1.png
│   └── rabbiyMQ_screen.png
├── server
│   ├── rabbitmq.go
│   └── server.go
└── srs.doc
```

---

## Installation

### Prerequisites
1. **Go** (1.16 or higher)
2. **RabbitMQ** (Installed and running)
3. **Protocol Buffers Compiler** (protoc)

### Steps
1. Clone the repository:
   ```bash
   git clone https://github.com/your-username/paper-storage-server.git
   cd paper-storage-server
   ```

2. Install dependencies:
   ```bash
   go mod tidy
   ```

3. Generate gRPC code (if needed):
   ```bash
   protoc --go_out=. --go-grpc_out=. paperpb/paper.proto
   ```

---

## Usage

### Start the Server
1. Navigate to the `paperserver` directory:
   ```bash
   cd cmd/paperserver
   ```

2. Run the server:
   ```bash
   go run main.go
   ```

### Use the Client
1. Navigate to the `paperclient` directory:
   ```bash
   cd cmd/paperclient
   ```

2. Use the client to interact with the server:
### Start The Client
```bash
go run main.go <server-address> 
```

#### Add a Paper
```bash
add <server-address> "<author>" "<title>" <file-path>
```
Example:
```bash
add localhost:50051 "John Doe" "Distributed Systems" paper.pdf
```

#### List All Papers
```bash
list <server-address>
```
Example:
```bash
list localhost:50051
```

#### Get Paper Details
```bash
detail <server-address> <paper-id>
```
Example:
```bash
detail localhost:50051 1
```

#### Fetch Paper Content
```bash
fetch <server-address> <paper-id>
```
Example:
```bash
fetch localhost:50051 1
```

---

## Notifications

When a new paper is added to the server, clients subscribed to the notification channel will receive a message via RabbitMQ. Notifications are handled asynchronously.

---

## Author

- ** Bisrat Asaye ** 

