package aggregator

import (
	protos "hyperledger_project/TxAggregator/protos"
	sender "hyperledger_project/TxAggregator/sender"
	"log"
	"sync"
	"time"

	gateway "github.com/hyperledger/fabric-sdk-go/pkg/gateway"
)

type Aggregator struct {
	TaggedTxChan        chan *protos.TaggedTransaction
	TaggedTxRsponseChan chan *protos.TaggedTransactionResponse
	TaggedTxSetChan     chan []*protos.TaggedTransaction
	WriteValueSetChan   chan map[string]*WriteValue
	TaggedTxSet         []*protos.TaggedTransaction
	WriteValueSet       map[string]*WriteValue
	gatewayContract     *gateway.Contract
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

		TaggedTxChan:        make(chan *protos.TaggedTransaction),
		TaggedTxRsponseChan: make(chan *protos.TaggedTransactionResponse),
		TaggedTxSetChan:     make(chan []*protos.TaggedTransaction),
		WriteValueSetChan:   make(chan map[string]*WriteValue),
		WriteValueSet:       make(map[string]*WriteValue),
		gatewayContract:     contract,
	}
	go aggregator.MakeTaggedTxset()
	go aggregator.Aggregate()
	go aggregator.SendTxProposals(contract)
	return aggregator

}
func (aggregator *Aggregator) MakeTaggedTxset() {
	var mutex = &sync.Mutex{}

	ticker := time.NewTicker(time.Millisecond * 1500)
	for {
		select {
		case TaggedTx := <-aggregator.GetTaggedTxReceiveChannel():
			log.Println("=========== ReceiveTaggedTxset ===========")
			log.Println(TaggedTx)
			aggregator.TaggedTxSet = append(aggregator.TaggedTxSet, TaggedTx)
			log.Println("=========== EndReceiveTaggedTxset ===========")

		case <-ticker.C:
			log.Println("=========== MakeTaggedTxset ===========")
			copyTaggedTxset := make([]*protos.TaggedTransaction, len(aggregator.TaggedTxSet))
			mutex.Lock()
			copy(copyTaggedTxset, aggregator.TaggedTxSet)
			aggregator.TaggedTxSet = nil
			log.Println("=========== EndMakeTaggedTxset ===========")
			aggregator.GetTaggedTxSetSendChannel() <- copyTaggedTxset
			mutex.Unlock()
			ticker = time.NewTicker(time.Millisecond * 1500)

		}

	}

}

// Aggregate TaggedTxset을 활용하여 연산 후 각 키에 Write할 value생성
func (aggregator *Aggregator) Aggregate() {
	var bytes []byte
	Response := 0
	var tempWriteValueSet = make(map[string]*WriteValue)
	for {
		select {

		case TaggedTxset := <-aggregator.GetTaggedTxSetReceiveChannel():
			log.Println("=========== StartAggregate ===========")
			for _, TaggedTx := range TaggedTxset {
				key := TaggedTx.Key
				result := aggregator.WriteValueSet[key]

				//empty struct check
				if result == nil {
					resultValue, _ := sender.ReadChaincode(TaggedTx.Functionname, TaggedTx.Key)
					result = &WriteValue{
						FunctionName: TaggedTx.Functionname,
						Key:          TaggedTx.Key,
						Value:        resultValue,
						WriteColumn:  TaggedTx.Fieldname,
						WriteValue:   int(TaggedTx.Operand),
					}
					aggregator.WriteValueSet[key] = result
					tempWriteValueSet = aggregator.WriteValueSet

				} else {

					// 사전 사후 검사
					if result.WriteValue < int(TaggedTx.Precondition) || result.WriteValue > int(TaggedTx.Postcondition) {

						bytes = []byte(TaggedTx.Key + " REJECT")
						aggregator.GetTaggedTxesponseSendChannel() <- &protos.TaggedTransactionResponse{
							Response: int32(Response),
							Payload:  bytes,
						}

					} else { // Operator 별 연산
						tempWriteValue := result.WriteValue
						if TaggedTx.Operator == int32(ADD) {
							tempWriteValue += int(TaggedTx.Operand)
						}

						if tempWriteValue < int(TaggedTx.Precondition) || tempWriteValue > int(TaggedTx.Postcondition) {
							Response = 500
							bytes = []byte(TaggedTx.Key + " " + TaggedTx.Fieldname + " REJECT")

						} else {
							result.WriteValue = tempWriteValue

							aggregator.WriteValueSet[key] = result
							tempWriteValueSet = aggregator.WriteValueSet
							Response = 200
							bytes = []byte(TaggedTx.Key + " SUCCESS")
						}

					}
				}
				aggregator.GetTaggedTxesponseSendChannel() <- &protos.TaggedTransactionResponse{
					Response: int32(Response),
					Payload:  bytes,
				}

			}
			log.Println("=========== EndAggregate ===========")
			aggregator.GetWriteValueSetSendChannel() <- tempWriteValueSet
			tempWriteValueSet = make(map[string]*WriteValue)
			// aggregator.WriteValueSet = make(map[string]*WriteValue)

		}

	}
}

func (aggregator *Aggregator) SendTxProposals(contract *gateway.Contract) {
	var bytes []byte
	for {
		select {
		case WriteValueSet := <-aggregator.GetWriteValueSetReceiveChannel():
			log.Println("=========== SendTxProposals ===========")
			for _, result := range WriteValueSet {

				err := sender.WriteChaincode(contract, result.FunctionName, result.Key, result.Value, result.WriteColumn, result.WriteValue)

				if err != nil {

					bytes = []byte(" MVCC CONFLICT")
					aggregator.GetTaggedTxesponseSendChannel() <- &protos.TaggedTransactionResponse{
						Response: int32(500),
						Payload:  bytes,
					}
				} else {
					bytes = []byte(" VALID")
					aggregator.GetTaggedTxesponseSendChannel() <- &protos.TaggedTransactionResponse{
						Response: int32(200),
						Payload:  bytes,
					}
				}

			}
		}

	}
}

// GetTaggedTxReceiveChanneln Receive TaggedTx Chan
func (aggregator *Aggregator) GetTaggedTxReceiveChannel() <-chan *protos.TaggedTransaction {
	return aggregator.TaggedTxChan
}
func (aggregator *Aggregator) GetTaggedTxSendChannel() chan<- *protos.TaggedTransaction {
	return aggregator.TaggedTxChan
}
func (aggregator *Aggregator) GetTaggedTxResponseReceiveChannel() <-chan *protos.TaggedTransactionResponse {
	return aggregator.TaggedTxRsponseChan
}
func (aggregator *Aggregator) GetTaggedTxesponseSendChannel() chan<- *protos.TaggedTransactionResponse {
	return aggregator.TaggedTxRsponseChan
}
func (aggregator *Aggregator) GetTaggedTxSetReceiveChannel() <-chan []*protos.TaggedTransaction {
	return aggregator.TaggedTxSetChan
}
func (aggregator *Aggregator) GetTaggedTxSetSendChannel() chan<- []*protos.TaggedTransaction {
	return aggregator.TaggedTxSetChan
}
func (aggregator *Aggregator) GetWriteValueSetReceiveChannel() <-chan map[string]*WriteValue {
	return aggregator.WriteValueSetChan
}
func (aggregator *Aggregator) GetWriteValueSetSendChannel() chan<- map[string]*WriteValue {
	return aggregator.WriteValueSetChan
}
