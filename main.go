package main

import (
	"context"
	"log"
	"net"
	"sync"

	pb "github.com/paxthemax/consignment-service/proto/consignment"
	"google.golang.org/grpc"
)

const (
	port = ":50051"
)

// Creator can create a consignment in itself.
type Creator interface {
	Create(consignment *pb.Consignment) (*pb.Consignment, error)
}

// AllGetter can get all of it's consignments.
type AllGetter interface {
	GetAll() []*pb.Consignment
}

// Repository is a simple in memory repository of Consignments. Has mutex lock.
type Repository struct {
	mx           sync.RWMutex
	consignments []*pb.Consignment
}

// Create will create a new consignment in the repository.
// It will return an error if consignment could not be created.
func (repo *Repository) Create(consignment *pb.Consignment) (*pb.Consignment, error) {
	updated := append(repo.consignments, consignment)

	repo.mx.Lock()
	repo.consignments = updated
	repo.mx.Unlock()

	return consignment, nil
}

// GetAll will return all consignments.
func (repo *Repository) GetAll() []*pb.Consignment {
	return repo.consignments
}

// Service is a ShippingService.
type Service struct {
	repo Repository
}

// NewService returns a new service
func NewService() *Service {
	return &Service{Repository{}}
}

// CreateConsignment will attempt to create a save a new consignment and return a response of that action.
// Returns error if the consignment was not created.
func (svc *Service) CreateConsignment(ctx context.Context, req *pb.Consignment) (*pb.CreateResponse, error) {
	cs, err := svc.repo.Create(req)
	if err != nil {
		return nil, err
	}

	return &pb.CreateResponse{
		Created:     true,
		Consignment: cs,
	}, nil
}

// GetConsignments will return all consignments from the repo.
func (svc *Service) GetConsignments(ctx context.Context, req *pb.GetRequest) (*pb.GetResponse, error) {
	result := svc.repo.GetAll()
	return &pb.GetResponse{Consignments: result}, nil
}

func main() {
	lis, err := net.Listen("tcp", port)
	if err != nil {
		log.Fatalf("Failed to listen on port, error = %v", err)
	}

	s := grpc.NewServer()
	pb.RegisterShippingServiceServer(s, NewService())

	log.Printf("Running server on port %s", port)
	if err := s.Serve(lis); err != nil {
		log.Fatalf("Internal server error = %v", err)
	}
}
