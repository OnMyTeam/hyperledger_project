package aggregator

import (
	protos "hyperledger_project/BWAggregator/protos"
	sender "hyperledger_project/BWAggregator/sender"
	"log"
	"sync"
	"time"

	gateway "github.com/hyperledger/fabric-sdk-go/pkg/gateway"
)

type Aggregator struct {
	BWTxChan          chan *protos.BWTransaction
	BWTxRsponseChan   chan *protos.BWTransactionResponse
	BWTxSetChan       chan []*protos.BWTransaction
	WriteValueSetChan chan map[string]*WriteValue
	BWTxSet           []*protos.BWTransaction
	WriteValueSet     map[string]*WriteValue
	gatewayContract   *gateway.Contract
}
type WriteValue struct {
	FunctionName string
	Key          string
	Value        string
	WriteColumn  string
	WriteValue   int
}

func Init(contract *gateway.Contract) *Aggregator {
	aggregator := &Aggregator{

		BWTxChan:          make(chan *protos.BWTransaction),
		BWTxRsponseChan:   make(chan *protos.BWTransactionResponse),
		BWTxSetChan:       make(chan []*protos.BWTransaction),
		WriteValueSetChan: make(chan map[string]*WriteValue),
		WriteValueSet:     make(map[string]*WriteValue),
		gatewayContract:   contract,
	}
	go aggregator.MakeBWTxset()
	go aggregator.Aggregate()
	go aggregator.SendTxProposals(contract)
	return aggregator

}
func (aggregator *Aggregator) MakeBWTxset() {
	var mutex = &sync.Mutex{}

	ticker := time.NewTicker(time.Millisecond * 1500)
	for {
		select {
		case BWTx := <-aggregator.GetBWTxReceiveChannel():
			log.Println("=========== ReceiveBWTxset ===========")
			log.Println(BWTx)
			aggregator.BWTxSet = append(aggregator.BWTxSet, BWTx)
			log.Println("=========== EndReceiveBWTxset ===========")

		case <-ticker.C:
			log.Println("=========== MakeBWTxset ===========")
			copyBWTxset := make([]*protos.BWTransaction, len(aggregator.BWTxSet))
			mutex.Lock()
			copy(copyBWTxset, aggregator.BWTxSet)
			aggregator.BWTxSet = nil
			log.Println("=========== EndMakeBWTxset ===========")
			aggregator.GetBWTxSetSendChannel() <- copyBWTxset
			mutex.Unlock()
			ticker = time.NewTicker(time.Millisecond * 1500)

		}

	}

}

// Aggregate BWTxset을 활용하여 연산 후 각 키에 Write할 value생성
func (aggregator *Aggregator) Aggregate() {
	var bytes []byte
	Response := 0
	for {
		select {
		case BWTxset := <-aggregator.GetBWTxSetReceiveChannel():
			log.Println("=========== Aggregate ===========")
			for _, BWTx := range BWTxset {
				key := BWTx.Key
				result := aggregator.WriteValueSet[key]

				//empty struct check
				if result == nil {
					resultValue, _ := sender.ReadChaincode(BWTx.Functionname, BWTx.Key)
					result = &WriteValue{
						FunctionName: BWTx.Functionname,
						Key:          BWTx.Key,
						Value:        resultValue,
						WriteColumn:  BWTx.Fieldname,
						WriteValue:   int(BWTx.Operand),
					}
					aggregator.WriteValueSet[key] = result
					Response = 200
					bytes = []byte(BWTx.Key + " SUCCESS")
				} else {
					// 사전 사후 검사
					if result.WriteValue < int(BWTx.Precondition) || result.WriteValue > int(BWTx.Postcondition) {

						bytes = []byte(BWTx.Key + " REJECT")

					} else { // Operator 별 연산
						tempWriteValue := result.WriteValue
						if BWTx.Operator == int32(ADD) {
							tempWriteValue += int(BWTx.Operand)
						}

						if tempWriteValue < int(BWTx.Precondition) || tempWriteValue > int(BWTx.Postcondition) {
							Response = 500
							bytes = []byte(BWTx.Key + " " + BWTx.Fieldname + " REJECT")

						} else {
							result.WriteValue = tempWriteValue
							aggregator.WriteValueSet[key] = result
							Response = 200
							bytes = []byte(BWTx.Key + " SUCCESS")
						}

					}
				}
				aggregator.GetBWTxesponseSendChannel() <- &protos.BWTransactionResponse{
					Response: int32(Response),
					Payload:  bytes,
				}

			}
			log.Println("=========== EndAggregate ===========")
			aggregator.GetWriteValueSetSendChannel() <- aggregator.WriteValueSet
			// aggregator.WriteValueSet = make(map[string]*WriteValue)

		}

	}
}

func (aggregator *Aggregator) SendTxProposals(contract *gateway.Contract) {

	for {
		select {
		case WriteValueSet := <-aggregator.GetWriteValueSetReceiveChannel():
			log.Println("=========== SendTxProposals ===========")
			for key, result := range WriteValueSet {
				log.Println(key, result)
				writeValue := 1000 - result.WriteValue
				sender.WriteChaincode(contract, result.FunctionName, result.Key, result.Value, result.WriteColumn, writeValue)

			}
		}

	}
}

// GetBWTxReceiveChanneln Receive BWTx Chan
func (aggregator *Aggregator) GetBWTxReceiveChannel() <-chan *protos.BWTransaction {
	return aggregator.BWTxChan
}
func (aggregator *Aggregator) GetBWTxSendChannel() chan<- *protos.BWTransaction {
	return aggregator.BWTxChan
}
func (aggregator *Aggregator) GetBWTxResponseReceiveChannel() <-chan *protos.BWTransactionResponse {
	return aggregator.BWTxRsponseChan
}
func (aggregator *Aggregator) GetBWTxesponseSendChannel() chan<- *protos.BWTransactionResponse {
	return aggregator.BWTxRsponseChan
}
func (aggregator *Aggregator) GetBWTxSetReceiveChannel() <-chan []*protos.BWTransaction {
	return aggregator.BWTxSetChan
}
func (aggregator *Aggregator) GetBWTxSetSendChannel() chan<- []*protos.BWTransaction {
	return aggregator.BWTxSetChan
}
func (aggregator *Aggregator) GetWriteValueSetReceiveChannel() <-chan map[string]*WriteValue {
	return aggregator.WriteValueSetChan
}
func (aggregator *Aggregator) GetWriteValueSetSendChannel() chan<- map[string]*WriteValue {
	return aggregator.WriteValueSetChan
}
