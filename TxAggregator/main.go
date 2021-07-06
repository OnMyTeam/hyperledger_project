package main

import (
	"context"
	"log"
	"net"

	aggregator "hyperledger_project/BWAggregator/aggregator"
	protos "hyperledger_project/BWAggregator/protos"
	sender "hyperledger_project/BWAggregator/sender"

	"google.golang.org/grpc"
)

const portNumber = "9000"

type BWAggregatorServer struct {
	protos.AggregatorServer
	aggregator.Aggregator
}

// ProcessProposal returns BWTransactionResponse
func (aggregator *BWAggregatorServer) ReceiveBWTransaction(ctx context.Context, req *protos.BWTransaction) (*protos.BWTransactionResponse, error) {
	// BWTxset생성을 위한 메세지 전달
	aggregator.GetBWTxSendChannel() <- req

	BWTxResponse := <-aggregator.GetBWTxResponseReceiveChannel()
	return BWTxResponse, nil
}

func main() {
	lis, err := net.Listen("tcp", ":"+portNumber)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	grpcServer := grpc.NewServer()
	var bWAggregatorServer BWAggregatorServer
	contract := sender.InitSender()
	bWAggregatorServer.Aggregator = *aggregator.Init(contract)
	protos.RegisterAggregatorServer(grpcServer, &bWAggregatorServer)

	log.Printf("Aggregate gRPC server Start!! %s port", portNumber)
	if err := grpcServer.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %s", err)
	}
}
