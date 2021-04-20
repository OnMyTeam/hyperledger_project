package main

import (
	"context"
	"fmt"
	"log"
	"net"

	userpb "hyperledger_project/BWAggregator/protos"

	"google.golang.org/grpc"
)

const portNumber = "9000"

type endorserServer struct {
	userpb.EndorserServer
}

// GetUser returns user message by user_id
func (s *endorserServer) ProcessProposal(ctx context.Context, req *userpb.SignedProposal) (*userpb.BWTransactionResponse, error) {

	userID := req.ProposalBytes
	fmt.Println(userID)
	return &userpb.BWTransactionResponse{
		Response: 2,
	}, nil
}

func main() {
	lis, err := net.Listen("tcp", ":"+portNumber)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	grpcServer := grpc.NewServer()
	userpb.RegisterEndorserServer(grpcServer, &endorserServer{})

	log.Printf("start gRPC server on !!%s port", portNumber)
	if err := grpcServer.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %s", err)
	}
}
