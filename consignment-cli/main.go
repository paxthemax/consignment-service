package main

import (
	"context"
	"encoding/json"
	"io/ioutil"
	"log"
	"os"

	pb "github.com/paxthemax/consignment-service/proto/consignment"
	"google.golang.org/grpc"
)

const (
	address         = "localhost:50051"
	defaultFilename = "manifest.json"
)

func parse(file string) (*pb.Consignment, error) {
	consignment := pb.Consignment{}

	data, err := ioutil.ReadFile(file)
	if err != nil {
		return nil, err
	}

	json.Unmarshal(data, &consignment)
	return &consignment, err
}

func main() {
	// Parse the consignments manifest:
	file := defaultFilename
	if len(os.Args) > 1 {
		file = os.Args[1]
	}
	consignment, err := parse(file)
	if err != nil {
		log.Fatalf("Could not parse consignments manifest, error = %v", err)
	}

	// Set up the connection with the gRPC server:
	conn, err := grpc.Dial(address, grpc.WithInsecure())
	if err != nil {
		log.Fatalf("Failed to connect to grpc server, error = %v", err)
	}
	defer conn.Close()

	// Create a new client:
	client := pb.NewShippingServiceClient(conn)

	// Create a consignment:
	resp, err := client.CreateConsignment(context.Background(), consignment)
	if err != nil {
		log.Fatalf("Failed to process consignment, error = %v", err)
	}

	log.Printf("Processed consignment, response = %v", resp)

	getAll, err := client.GetConsignments(context.Background(), &pb.GetRequest{})
	if err != nil {
		log.Fatalf("Failed to fetch consignments, err = %v", err)
	}

	log.Printf("Consignments:")
	for _, v := range getAll.Consignments {
		log.Println(v)
	}
}
