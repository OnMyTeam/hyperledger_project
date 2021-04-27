package aggregator

import (
	"fmt"
	protos "hyperledger_project/BWAggregator/protos"
	"sync"
	"time"
)

type Aggregator struct {
	BWTxChan    chan *protos.BWTransaction
	BWTxSetChan chan []*protos.BWTransaction
	BWTxSet     []*protos.BWTransaction
}

func Init() *Aggregator {
	aggregator := &Aggregator{

		BWTxChan:    make(chan *protos.BWTransaction, 5),
		BWTxSetChan: make(chan []*protos.BWTransaction, 5),
		BWTxSet:     make([]*protos.BWTransaction, 5),
	}
	go aggregator.MakeBWTxset()
	go aggregator.Aggregate()
	return aggregator

}
func (aggregator *Aggregator) MakeBWTxset() {
	var mutex = &sync.Mutex{}

	ticker := time.NewTicker(time.Millisecond * 1500)
	for {
		select {
		case BWTx := <-aggregator.GetBWTxReceiveChannel():
			fmt.Println("=========== MakeBWTxset ===========")
			fmt.Println(BWTx)
			aggregator.BWTxSet = append(aggregator.BWTxSet, BWTx)

		case <-ticker.C:
			copyBWTxset := make([]*protos.BWTransaction, len(aggregator.BWTxSet))
			mutex.Lock()
			copy(copyBWTxset, aggregator.BWTxSet)
			aggregator.BWTxSet = make([]*protos.BWTransaction, 5)
			aggregator.GetBWTxSetSendChannel() <- copyBWTxset
			mutex.Unlock()
			ticker = time.NewTicker(time.Millisecond * 1500)
		}
		fmt.Println("=========== MakeBWTxsetEnd ===========")

	}

}

// Aggregate BWTxset을 활용하여 연산 후 각 키에 Write할 value생성
func (aggregator *Aggregator) Aggregate() {
	for {
		select {
		case BWTxset := <-aggregator.GetBWTxSetReceiveChannel():
			fmt.Println("=========== Aggregator ===========")
			for i, s := range BWTxset {
				fmt.Println(i, s)
			}
		}
		fmt.Println("=========== AggregatorEnd ===========")

	}
}
func (aggregator *Aggregator) MakeWriteValue() {

}
func (aggregator *Aggregator) SendTxProposals() {

}
func (aggregator *Aggregator) GetBWTxReceiveChannel() <-chan *protos.BWTransaction {
	return aggregator.BWTxChan
}
func (aggregator *Aggregator) GetBWTxSendChannel() chan<- *protos.BWTransaction {
	return aggregator.BWTxChan
}
func (aggregator *Aggregator) GetBWTxSetReceiveChannel() <-chan []*protos.BWTransaction {
	return aggregator.BWTxSetChan
}
func (aggregator *Aggregator) GetBWTxSetSendChannel() chan<- []*protos.BWTransaction {
	return aggregator.BWTxSetChan
}
