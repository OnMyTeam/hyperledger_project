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
		TaggedTxRsponseChan: make(chan *protos.TaggedTransactionResponse, 3),
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
// 특정 시간동안 TaggedTxSet 수집
// TaggedTxSet 수집 후 키 별 통합을 위해 채널을 통해 Aggregate함수로 전송
func (aggregator *Aggregator) MakeTaggedTxset() {
	var mutex = &sync.Mutex{}
	// 임의타임 아웃 시간 설정
	ticker := time.NewTicker(time.Millisecond * 1500)
	for {
		select {
		case TaggedTx := <-aggregator.GetTaggedTxReceiveChannel():
			log.Println("=========== ReceiveTaggedTxset ===========")
			log.Println(TaggedTx)
			// TaggedTx 수집
			aggregator.TaggedTxSet = append(aggregator.TaggedTxSet, TaggedTx)
			log.Println("=========== EndReceiveTaggedTxset ===========")

		// 타임아웃
		case <-ticker.C:
			log.Println("=========== MakeTaggedTxset ===========")
			mutex.Lock()
			copyTaggedTxset := make([]*protos.TaggedTransaction, len(aggregator.TaggedTxSet))
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
					Response = 2001
					bytes = []byte(TaggedTx.Key + " SUCCESS")
					tempWriteValueSet = aggregator.WriteValueSet

				} else {

					// 각 키 별 통합 후 사전 사후 검사
					// ex) 0미만 1000초과시 reject 처리
					if result.WriteValue < int(TaggedTx.Precondition) || result.WriteValue > int(TaggedTx.Postcondition) {

						Response = 500
						bytes = []byte(TaggedTx.Key + " " + TaggedTx.Fieldname + " REJECT")

					} else { // Operator 별 연산
						tempWriteValue := result.WriteValue
						if TaggedTx.Operator == int32(ADD) {
							tempWriteValue += int(TaggedTx.Operand)
						}

						if tempWriteValue < int(TaggedTx.Precondition) || tempWriteValue > int(TaggedTx.Postcondition) {
							Response = 5001
							bytes = []byte(TaggedTx.Key + " " + TaggedTx.Fieldname + " REJECT")

						} else {
							result.WriteValue = tempWriteValue

							aggregator.WriteValueSet[key] = result
							tempWriteValueSet = aggregator.WriteValueSet
							Response = 2001
							bytes = []byte(TaggedTx.Key + " SUCCESS")
						}

					}
				}
				aggregator.GetTaggedTxReponseSendChannel() <- &protos.TaggedTransactionResponse{
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
// 각 키별 통합 후 패브릭네트워크로 TxProposal 전송
func (aggregator *Aggregator) SendTxProposals(contract *gateway.Contract) {
	// var bytes []byte
	for {
		select {
		case WriteValueSet := <-aggregator.GetWriteValueSetReceiveChannel():
			log.Println("=========== SendTxProposals ===========")
			for _, result := range WriteValueSet {

				sender.WriteChaincode(contract, result.FunctionName, result.Key, result.Value, result.WriteColumn, result.WriteValue)

				// if err != nil {

				// 	bytes = []byte(" MVCC CONFLICT")
				// 	aggregator.GetTaggedTxReponseSendChannel() <- &protos.TaggedTransactionResponse{
				// 		Response: int32(500),
				// 		Payload:  bytes,
				// 	}

				// } else {
				// 	bytes = []byte(" VALID")
				// 	aggregator.GetTaggedTxReponseSendChannel() <- &protos.TaggedTransactionResponse{
				// 		Response: int32(200),
				// 		Payload:  bytes,
				// 	}
				// }

			}
		}

	}
}

// GetTaggedTxReceiveChannel Receive TaggedTx 
func (aggregator *Aggregator) GetTaggedTxReceiveChannel() <-chan *protos.TaggedTransaction {
	return aggregator.TaggedTxChan
}
// GetTaggedTxSendChannel Send TaggedTx 
func (aggregator *Aggregator) GetTaggedTxSendChannel() chan<- *protos.TaggedTransaction {
	return aggregator.TaggedTxChan
}
// GetTaggedTxResponseReceiveChannel Receive TaggedTxResponse 
func (aggregator *Aggregator) GetTaggedTxResponseReceiveChannel() <-chan *protos.TaggedTransactionResponse {
	return aggregator.TaggedTxRsponseChan
}
// GetTaggedTxReponseSendChannel Send TaggedTxResponse 
func (aggregator *Aggregator) GetTaggedTxReponseSendChannel() chan<- *protos.TaggedTransactionResponse {
	return aggregator.TaggedTxRsponseChan
}
// GetTaggedTxSetReceiveChannel Receive TaggedTxSet
func (aggregator *Aggregator) GetTaggedTxSetReceiveChannel() <-chan []*protos.TaggedTransaction {
	return aggregator.TaggedTxSetChan
}
// GetTaggedTxSetSendChannel Send TaggedTxSet
func (aggregator *Aggregator) GetTaggedTxSetSendChannel() chan<- []*protos.TaggedTransaction {
	return aggregator.TaggedTxSetChan
}
// GetWriteValueSetReceiveChannel Receive WriteValueSet each for Key
func (aggregator *Aggregator) GetWriteValueSetReceiveChannel() <-chan map[string]*WriteValue {
	return aggregator.WriteValueSetChan
}
// GetWriteValueSetSendChannel Send WriteValueSet each for Key
func (aggregator *Aggregator) GetWriteValueSetSendChannel() chan<- map[string]*WriteValue {
	return aggregator.WriteValueSetChan
}
