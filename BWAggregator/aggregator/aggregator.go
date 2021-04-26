package aggregator

import (
	protos "hyperledger_project/BWAggregator/protos"
)

type Aggregator struct {
	BWTxChan chan *protos.BWTransaction
}

func Init() *Aggregator {
	aggregator := &Aggregator{

		BWTxChan: make(chan *protos.BWTransaction, 5),
	}
	return aggregator

}

func (aggregator *Aggregator) GetBWTxReceiveChannel() <-chan *protos.BWTransaction {
	return aggregator.BWTxChan
}
func (aggregator *Aggregator) GetBWTxSendChannel() chan<- *protos.BWTransaction {
	return aggregator.BWTxChan
}
