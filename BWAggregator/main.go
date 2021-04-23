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
func (s *endorserServer) ProcessProposal(ctx context.Context, req *userpb.BWTransaction) (*userpb.BWTransactionResponse, error) {

	functionName := req.Functionname
	Key := req.Key
	Fieldname := req.Fieldname
	Operator := req.Operator
	Operand := req.Operand
	Precondition := req.Precondition
	Postcondition := req.Postcondition
	fmt.Println(functionName)
	fmt.Println(Key)
	fmt.Println(Fieldname)
	fmt.Println(Operator)
	fmt.Println(Operand)
	fmt.Println(Precondition)
	fmt.Println(Postcondition)
	var message string
	message = "success"
	b := []byte(message)
	return &userpb.BWTransactionResponse{
		Response: 2,
		Payload:  b,
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
