package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/gorilla/mux"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	pb "taxonomy-client/taxonomy-client/proto/taxonomy/v1"
)

const (
	defaultTaxonomyID = "1701181887VZ"
	// Use the local mock service instead of the remote one
	taxonomyServiceAddr = "localhost:8083"
	apiPort             = "8082"
)

// Global gRPC client
var taxonomyClient pb.TaxonomyClient

func main() {
	// Set up a connection to the gRPC service
	log.Printf("Connecting to taxonomy service at %s", taxonomyServiceAddr)
	conn, err := grpc.Dial(taxonomyServiceAddr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("Failed to connect to taxonomy service: %v", err)
	}
	defer conn.Close()

	// Create a gRPC client
	taxonomyClient = pb.NewTaxonomyClient(conn)

	// Create a new router
	r := mux.NewRouter()

	// Register routes
	r.HandleFunc("/api/taxonomy/{id}", getTaxonomyHandler).Methods("GET")
	r.HandleFunc("/api/taxonomy", getTaxonomyHandler).Methods("GET") // Default ID

	// Start the HTTP server
	fmt.Printf("Starting API server on port %s...\n", apiPort)
	log.Fatal(http.ListenAndServe(":"+apiPort, r))
}

func getTaxonomyHandler(w http.ResponseWriter, r *http.Request) {
	// Get taxonomy ID from URL or use default
	vars := mux.Vars(r)
	taxonomyID, ok := vars["id"]
	if !ok {
		taxonomyID = defaultTaxonomyID
	}

	// Create a context with timeout
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	log.Printf("\n===== GRPC REQUEST =====")
	log.Printf("Method: GetTaxonomyById")
	log.Printf("Target Service: %s", taxonomyServiceAddr)
	log.Printf("Request: &pb.GetTaxonomyByIdRequest{TaxonomyId: %s}", taxonomyID)

	// Create a metadata collection context
	md := make(map[string]string)
	for k, v := range r.Header {
		if len(v) > 0 {
			md[k] = v[0]
		}
	}

	// Make the gRPC call
	log.Printf("Sending gRPC request to taxonomy service...")
	requestStartTime := time.Now()
	response, err := taxonomyClient.GetTaxonomyById(ctx, &pb.GetTaxonomyByIdRequest{
		TaxonomyId: taxonomyID,
	})
	requestDuration := time.Since(requestStartTime)
	log.Printf("gRPC request completed in %v", requestDuration)

	if err != nil {
		log.Printf("\n===== GRPC ERROR =====")
		log.Printf("Error calling taxonomy service: %v", err)
		http.Error(w, fmt.Sprintf("Error fetching taxonomy: %v", err), http.StatusInternalServerError)
		return
	}

	// Log complete gRPC response details
	log.Printf("\n===== GRPC RESPONSE =====")

	// Convert to JSON with indentation for better readability
	responseJSON, _ := json.MarshalIndent(response, "", "  ")

	// Print the full response with clear separation
	log.Printf("\nRESPONSE BODY:\n%s", string(responseJSON))

	// If there's protobuf metadata, display it
	log.Printf("\n===== GRPC METADATA =====")
	log.Printf("Message type: %T", response)
	log.Printf("Service: pb.TaxonomyClient")

	// Set content type
	w.Header().Set("Content-Type", "application/json")

	// Check if the client wants the raw format or transformed format
	format := r.URL.Query().Get("format")

	if format == "transformed" {
		// Return transformed response similar to grpcurl format
		transformedResponse := transformResponse(response)
		if err := json.NewEncoder(w).Encode(transformedResponse); err != nil {
			http.Error(w, fmt.Sprintf("Error encoding response: %v", err), http.StatusInternalServerError)
		}
	} else {
		// Return raw response
		if err := json.NewEncoder(w).Encode(response); err != nil {
			http.Error(w, fmt.Sprintf("Error encoding response: %v", err), http.StatusInternalServerError)
		}
	}
}

// Transformation structures to match the original grpcurl output
type TransformedResponse struct {
	TaxonomyInfo *TransformedTaxonomyInfo `json:"taxonomyInfo"`
}

type TransformedTaxonomyInfo struct {
	TaxonomyId  string                      `json:"taxonomyId"`
	Name        string                      `json:"name"`
	Description string                      `json:"description"`
	Version     string                      `json:"version"`
	TenantId    string                      `json:"tenantId"`
	CreatedAt   string                      `json:"createdAt"`
	UpdatedAt   string                      `json:"updatedAt"`
	Nodes       map[string]*TransformedNode `json:"nodes"`
	RootNodes   []string                    `json:"rootNodes"`
	Levels      []string                    `json:"levels"`
}

type TransformedNode struct {
	Id           string                    `json:"id"`
	Name         string                    `json:"name"`
	Description  string                    `json:"description,omitempty"`
	ShortCode    string                    `json:"shortCode,omitempty"`
	ParentNode   string                    `json:"parentNode,omitempty"`
	NodeType     string                    `json:"nodeType"`
	RelatedNodes []*TransformedRelatedNode `json:"relatedNodes,omitempty"`
	Children     []string                  `json:"children,omitempty"`
	Ancestors    []string                  `json:"ancestors,omitempty"`
	Inactive     bool                      `json:"inactive,omitempty"`
}

type TransformedRelatedNode struct {
	TaxonomyId string `json:"taxonomyId"`
	NodeId     string `json:"nodeId"`
}

// Map to convert node type integers to strings
var nodeTypeMap = map[int32]string{
	0: "NODETYPE_UNDEFINED",
	1: "CLASS",
	2: "SUBJECT",
	3: "SUPER_TOPIC",
	4: "TOPIC",
	5: "SUBTOPIC",
	6: "CONCEPT",
}

func transformResponse(response *pb.GetTaxonomyByIdResponse) *TransformedResponse {
	info := response.GetTaxonomyInfo()

	// Create the transformed taxonomy info
	transformedInfo := &TransformedTaxonomyInfo{
		TaxonomyId:  info.GetTaxonomyId(),
		Name:        info.GetName(),
		Description: info.GetDescription(),
		Version:     info.GetVersion(),
		TenantId:    info.GetTenantId(),
		CreatedAt:   fmt.Sprintf("%d", info.GetCreatedAt()),
		UpdatedAt:   fmt.Sprintf("%d", info.GetUpdatedAt()),
		Nodes:       make(map[string]*TransformedNode),
		RootNodes:   info.GetRootNodes(),
		Levels:      info.GetLevels(),
	}

	// Transform each node
	for id, node := range info.GetNodes() {
		transformedNode := &TransformedNode{
			Id:          node.GetId(),
			Name:        node.GetName(),
			Description: node.GetDescription(),
			ShortCode:   node.GetShortCode(),
			ParentNode:  node.GetParentNode(),
			NodeType:    nodeTypeMap[int32(node.GetNodeType())], // Convert int to string
			Children:    node.GetChildren(),
			Ancestors:   node.GetAncestors(),
			Inactive:    node.GetInactive(),
		}

		// Transform related nodes
		if len(node.GetRelatedNodes()) > 0 {
			transformedNode.RelatedNodes = make([]*TransformedRelatedNode, 0, len(node.GetRelatedNodes()))
			for _, relNode := range node.GetRelatedNodes() {
				transformedNode.RelatedNodes = append(transformedNode.RelatedNodes, &TransformedRelatedNode{
					TaxonomyId: relNode.GetTaxonomyId(),
					NodeId:     relNode.GetNodeId(),
				})
			}
		}

		transformedInfo.Nodes[id] = transformedNode
	}

	return &TransformedResponse{
		TaxonomyInfo: transformedInfo,
	}
}
