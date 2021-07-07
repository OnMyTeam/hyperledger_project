package main

import (
	"context"
	"log"
	"net"

	aggregator "hyperledger_project/TxAggregator/aggregator"
	protos "hyperledger_project/TxAggregator/protos"
	sender "hyperledger_project/TxAggregator/sender"

	"google.golang.org/grpc"
)

const portNumber = "9000"

type TxAggregatorServer struct {
	protos.AggregatorServer
	aggregator.Aggregator
}

// ProcessProposal returns TaggedTransactionResponse
func (aggregator *TxAggregatorServer) ReceiveTaggedTransaction(ctx context.Context, req *protos.TaggedTransaction) (*protos.TaggedTransactionResponse, error) {
	// TaggedTxset생성을 위한 메세지 전달
	aggregator.GetTaggedTxSendChannel() <- req

	TaggedTxResponse := <-aggregator.GetTaggedTxResponseReceiveChannel()
	log.Println("TaggedTxResponse :", TaggedTxResponse)
	return TaggedTxResponse, nil
}

func main() {
	lis, err := net.Listen("tcp", ":"+portNumber)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	grpcServer := grpc.NewServer()
	var txAggregatorServer TxAggregatorServer
	contract := sender.InitSender()
	txAggregatorServer.Aggregator = *aggregator.Init(contract)
	protos.RegisterAggregatorServer(grpcServer, &txAggregatorServer)

	log.Printf("Aggregate gRPC server Start!! %s port", portNumber)
	if err := grpcServer.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %s", err)
	}
}
