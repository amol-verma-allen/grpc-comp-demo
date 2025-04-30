package main

import (
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"net"
	"os"

	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"

	pb "taxonomy-client/taxonomy-client/proto/taxonomy/v1"
)

type mockTaxonomyServer struct {
	pb.UnimplementedTaxonomyServer
	mockData *pb.GetTaxonomyByIdResponse
}

func (s *mockTaxonomyServer) GetTaxonomyById(ctx context.Context, req *pb.GetTaxonomyByIdRequest) (*pb.GetTaxonomyByIdResponse, error) {
	log.Printf("Mock service received request for taxonomy ID: %s", req.TaxonomyId)
	return s.mockData, nil
}

func main() {
	port := flag.Int("port", 8083, "The server port")
	jsonFile := flag.String("json", "taxonomy_raw.json", "Path to JSON file with mock taxonomy data")
	flag.Parse()

	// Load mock data from JSON file
	mockData, err := loadMockData(*jsonFile)
	if err != nil {
		log.Fatalf("Failed to load mock data: %v", err)
	}

	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", *port))
	if err != nil {
		log.Fatalf("Failed to listen: %v", err)
	}

	grpcServer := grpc.NewServer()
	pb.RegisterTaxonomyServer(grpcServer, &mockTaxonomyServer{mockData: mockData})

	// Register reflection service for debugging with grpcurl
	reflection.Register(grpcServer)

	log.Printf("Mock taxonomy service is running on port %d", *port)
	log.Printf("Using mock data from: %s", *jsonFile)
	log.Printf("Server is ready to accept requests!")

	if err := grpcServer.Serve(lis); err != nil {
		log.Fatalf("Failed to serve: %v", err)
	}
}

func loadMockData(file string) (*pb.GetTaxonomyByIdResponse, error) {
	// Check if file exists
	if _, err := os.Stat(file); os.IsNotExist(err) {
		// Create an empty response if file doesn't exist
		log.Printf("Mock data file %s not found, creating a basic empty response", file)
		return &pb.GetTaxonomyByIdResponse{
			TaxonomyInfo: &pb.TaxonomyInfo{
				TaxonomyId: "mock-taxonomy-id",
				Name:       "Mock Taxonomy",
				Nodes:      make(map[string]*pb.Node),
			},
		}, nil
	}

	// Read JSON file
	jsonData, err := ioutil.ReadFile(file)
	if err != nil {
		return nil, fmt.Errorf("could not read file: %v", err)
	}

	// Parse JSON into response object
	var response pb.GetTaxonomyByIdResponse
	if err := json.Unmarshal(jsonData, &response); err != nil {
		return nil, fmt.Errorf("could not unmarshal json: %v", err)
	}

	return &response, nil
}
