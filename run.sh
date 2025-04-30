#!/bin/bash

echo "Current directory: $(pwd)"
echo "Listing files:"
ls -la

# Check if taxonomy file exists
if [ -f "taxonomy_raw.json" ]; then
  echo "Found taxonomy_raw.json file ($(du -h taxonomy_raw.json | cut -f1))."
else
  echo "WARNING: taxonomy_raw.json file not found!"
fi

# In Docker, we use pre-built binaries, but locally we might need to build
if [ ! -f "bin/mock-server" ] || [ ! -f "bin/api-server" ]; then
  echo "Building applications..."
  mkdir -p bin
  go build -o bin/mock-server cmd/mock-taxonomy-service/main.go
  go build -o bin/api-server cmd/api-server/main.go
fi

# Start the mock taxonomy service in the background
echo "Starting mock taxonomy service on port 8083..."
./bin/mock-server --port 8083 --json taxonomy_raw.json &
MOCK_PID=$!

# Wait for the mock service to start
echo "Waiting for mock service to initialize..."
sleep 2

# Start the API server
echo "Starting API server on port 8082..."
./bin/api-server &
API_PID=$!

# Wait for any signal to terminate
echo "Services are running. Binaries for Services have been created"
trap "echo 'Shutting down services...'; kill $MOCK_PID $API_PID; echo 'All services stopped.'; exit 0" SIGINT SIGTERM

# Keep the script running
kill $API_PID $MOCK_PID 