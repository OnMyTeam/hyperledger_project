package aggregator

import (
	"fmt"
	protos "hyperledger_project/BWAggregator/protos"
	sender "hyperledger_project/BWAggregator/sender"
	"sync"
	"time"
)

type Aggregator struct {
	BWTxChan        chan *protos.BWTransaction
	BWTxSetChan     chan []*protos.BWTransaction
	BWTxSet         []*protos.BWTransaction
	ResultAggregate map[string]Result
}
type Result struct {
	FunctionName string
	Key          string
	WriteValue   int
}

func Init() *Aggregator {
	aggregator := &Aggregator{

		BWTxChan:        make(chan *protos.BWTransaction),
		BWTxSetChan:     make(chan []*protos.BWTransaction),
		ResultAggregate: make(map[string]Result),
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
			fmt.Println("=========== ReceiveBWTxset ===========")
			fmt.Println(BWTx)
			aggregator.BWTxSet = append(aggregator.BWTxSet, BWTx)
			fmt.Println("=========== EndReceiveBWTxset ===========")

		case <-ticker.C:
			fmt.Println("=========== MakeBWTxset ===========")
			copyBWTxset := make([]*protos.BWTransaction, len(aggregator.BWTxSet))
			mutex.Lock()
			copy(copyBWTxset, aggregator.BWTxSet)
			aggregator.BWTxSet = nil
			aggregator.GetBWTxSetSendChannel() <- copyBWTxset
			mutex.Unlock()
			ticker = time.NewTicker(time.Millisecond * 1500)
			fmt.Println("=========== EndMakeBWTxset ===========")
		}

	}

}

// Aggregate BWTxset을 활용하여 연산 후 각 키에 Write할 value생성
func (aggregator *Aggregator) Aggregate() {
	for {
		select {
		case BWTxset := <-aggregator.GetBWTxSetReceiveChannel():
			fmt.Println("=========== Aggregate ===========")
			for _, BWTx := range BWTxset {
				key := BWTx.Key
				result := aggregator.ResultAggregate[key]

				//empty struct check
				if result.FunctionName == "" {
					result.FunctionName = BWTx.Functionname
					result.Key = BWTx.Key
					result.WriteValue = int(BWTx.Operand)
					aggregator.ResultAggregate[key] = result
				} else {
					// 사전 사후 검사
					if result.WriteValue < int(BWTx.Precondition) && result.WriteValue > int(BWTx.Postcondition) {
						fmt.Println("no!!!!!!!!")
					} else { // Operator 별 연산
						if BWTx.Operator == int32(ADD) {
							result.WriteValue += int(BWTx.Operand)
						}
						aggregator.ResultAggregate[key] = result
					}
				}
				fmt.Println(result.FunctionName)
				fmt.Println(result.Key)
				fmt.Println(result.WriteValue)

			}
			if aggregator.ResultAggregate["CAR0"].FunctionName != "" {
				sender.WriteChaincode(aggregator.ResultAggregate["CAR0"].FunctionName, aggregator.ResultAggregate["CAR0"].Key, aggregator.ResultAggregate["CAR0"].WriteValue)
			}

			fmt.Println("=========== EndAggregate ===========")
		}

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
