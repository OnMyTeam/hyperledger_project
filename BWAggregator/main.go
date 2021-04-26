package main

import (
	"context"
	"fmt"
	"log"
	"net"

	aggregator "hyperledger_project/BWAggregator/aggregator"
	protos "hyperledger_project/BWAggregator/protos"

	"google.golang.org/grpc"
)

const portNumber = "9000"

type BWAggregatorServer struct {
	protos.AggregatorServer
	aggregator.Aggregator
}

// ProcessProposal returns BWTransactionResponse
func (aggregator *BWAggregatorServer) ReceiveBWTransaction(ctx context.Context, req *protos.BWTransaction) (*protos.BWTransactionResponse, error) {
	aggregator.GetBWTxSendChannel() <- req
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
	msg := <-aggregator.GetBWTxReceiveChannel()
	fmt.Println("msg =>", msg)
	var message string
	message = "success"
	b := []byte(message)
	return &protos.BWTransactionResponse{
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
	var bWAggregatorServer BWAggregatorServer
	bWAggregatorServer.Aggregator = *aggregator.Init()
	protos.RegisterAggregatorServer(grpcServer, &bWAggregatorServer)

	log.Printf("start gRPC server on !!%s port", portNumber)
	if err := grpcServer.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %s", err)
	}
}
