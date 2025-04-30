# Taxonomy gRPC Client

This is a Go service that provides a gRPC client to interact with the taxonomy-service.taxonomy.allen-stage.in API.

## Features

- Direct client that connects to the taxonomy service
- Transformed client that matches the exact output format of grpcurl
- REST API server that exposes the taxonomy service via HTTP endpoints
- Server that proxies requests to the taxonomy service
- Client that can make requests to the server
- Support for the `GetTaxonomyById` RPC method

## Prerequisites

- Go 1.16 or later
- Protocol Buffers compiler (protoc)
- Go gRPC plugins

## Installation

1. Clone the repository
2. Install dependencies:

```bash
go mod tidy
```

## Usage

### Using Docker (Recommended)

The easiest way to run both the API server and mock server is using Docker:

```bash
# Build and start the Docker containers
docker-compose up -d

# Check the logs
docker-compose logs -f
```

This will start:

1. A mock server on port 50051 that simulates the taxonomy service
2. An API server on port 8082 that connects to the mock server

The API will be accessible at http://localhost:8082/api/taxonomy

To stop the services:

```bash
docker-compose down
```

### Using the Mock Server

The mock server simulates the taxonomy service using local data:

```bash
# Build and start the mock server
go build -o bin/mock-server cmd/mock-server/main.go cmd/mock-server/importData.go
./bin/mock-server

# Import data from a specific file (optional)
./bin/mock-server --import taxonomy_raw.json
```

The mock server runs on port 50051 by default and responds to the same gRPC methods as the actual taxonomy service.

### Using the REST API Server

Alternatively, you can build and run the API server directly:

```bash
# Build and start the API server
go build -o bin/api-server cmd/api-server/main.go
./bin/api-server
```

Once running, you can access the API endpoints:

- Get taxonomy with default ID (1701181887VZ):

  ```
  GET http://localhost:8082/api/taxonomy
  ```

- Get taxonomy with specific ID:

  ```
  GET http://localhost:8082/api/taxonomy/{id}
  ```

- Get taxonomy with transformed format (matching grpcurl output):
  ```
  GET http://localhost:8082/api/taxonomy?format=transformed
  GET http://localhost:8082/api/taxonomy/{id}?format=transformed
  ```

### Using the Transformed Client (Recommended)

The transformed client connects directly to the taxonomy service and formats the output to exactly match the grpcurl command:

```bash
# Build the transformed client
go build -o bin/transform-client cmd/transform-client/main.go

# Run with default taxonomy ID (1701181887VZ)
./bin/transform-client

# Specify a custom taxonomy ID
./bin/transform-client YOUR_TAXONOMY_ID
```

### Using the Direct Client

The direct client connects directly to the taxonomy service:

```bash
# Build the direct client
go build -o bin/direct-client cmd/direct-client/main.go

# Run with default taxonomy ID (1701181887VZ)
./bin/direct-client

# Specify a custom taxonomy ID
./bin/direct-client YOUR_TAXONOMY_ID

# Show all nodes in the taxonomy (can be large output)
./bin/direct-client YOUR_TAXONOMY_ID --show-nodes
```

### Starting the Server

The server runs on port 8081 and forwards requests to the taxonomy service.

```bash
go run cmd/server/main.go
```

### Using the Client with Local Server

The client connects to the local server (which forwards to the taxonomy service):

```bash
# Use default taxonomy ID (1701181887VZ)
go run cmd/client/main.go

# Specify a custom taxonomy ID
go run cmd/client/main.go YOUR_TAXONOMY_ID

# Show all nodes in the taxonomy (can be large output)
go run cmd/client/main.go YOUR_TAXONOMY_ID --show-nodes
```

### Direct gRPCurl Command

For reference, here's the direct gRPCurl command that this client replaces:

```bash
grpcurl -plaintext -d '{"taxonomy_id": "1701181887VZ"}' taxonomy-service.taxonomy.allen-stage.in:80 taxonomy.v1.Taxonomy/GetTaxonomyById
```

## Output Format

The transformed client outputs JSON in the exact same format as the original grpcurl command, with these characteristics:

- Uses camelCase field naming (taxonomyInfo, nodeType, etc.)
- Represents NodeType as strings (e.g., "CLASS", "SUBJECT", "TOPIC")
- Identical structure and field organization

## Project Structure

- `cmd/api-server/`: REST API server implementation (exposes HTTP endpoints that call the gRPC service)
- `cmd/transform-client/`: Transformed client implementation (connects directly to the taxonomy service with formatted output)
- `cmd/raw-client/`: Raw client implementation (connects directly to the taxonomy service without format transformation)
- `cmd/direct-client/`: Direct client implementation (connects directly to the taxonomy service)
- `cmd/server/`: Server implementation (proxies requests to the taxonomy service)
- `cmd/client/`: Client implementation (connects to our local server)
- `proto/`: Protocol Buffer definitions
# grpc-comp-demo
