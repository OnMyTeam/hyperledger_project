package aggregator

import (
	"fmt"
	protos "hyperledger_project/BWAggregator/protos"
	"sync"
	"time"
)

type Aggregator struct {
	BWTxChan chan *protos.BWTransaction
	BWTxSet  []*protos.BWTransaction
}

func Init() *Aggregator {
	aggregator := &Aggregator{

		BWTxChan: make(chan *protos.BWTransaction, 5),
		BWTxSet:  make([]*protos.BWTransaction, 5),
	}
	go aggregator.MakeBWTxset()
	return aggregator

}
func (aggregator *Aggregator) MakeBWTxset() {
	var mutex = &sync.Mutex{}
	count := 0
	// timer := time.NewTimer(time.Second * 3)
	ticker := time.NewTicker(time.Millisecond * 1500)
	for {
		select {
		case BWTx := <-aggregator.GetBWTxReceiveChannel():
			fmt.Println(BWTx)
			aggregator.BWTxSet = append(aggregator.BWTxSet, BWTx)
			count++
		case <-ticker.C:
			mutex.Lock()
			for i, s := range aggregator.BWTxSet {
				fmt.Println(i, s)
			}
			fmt.Println("=================")
			mutex.Unlock()
			ticker = time.NewTicker(time.Millisecond * 1500)
		}

	}

}
func (aggregator *Aggregator) Aggregate() {

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
