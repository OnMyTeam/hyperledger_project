export PATH=$PATH:/usr/local/go/bin:$GOPATH/bin:~/go/src/hyperledger_project/bin
export FABRIC_CFG_PATH=/home/mypaper/go/src/hyperledger_project/config/
export CORE_PEER_TLS_ENABLED=true
export CORE_PEER_ADDRESS=localhost:7051
export CORE_PEER_LOCALMSPID="Org1MSP"
export CORE_PEER_TLS_ROOTCERT_FILE=/home/mypaper/go/src/hyperledger_project/test-network/organizations/peerOrganizations/org1.example.com/peers/peer0.org1.example.com/tls/ca.crt
export CORE_PEER_MSPCONFIGPATH=/home/mypaper/go/src/hyperledger_project/test-network/organizations/peerOrganizations/org1.example.com/users/Admin@org1.example.com/msp
index=1
number=0

type=$1

if [ $type == "InitLedger" ];
then

    peer chaincode invoke -o localhost:7050 --ordererTLSHostnameOverride orderer.example.com --tls --cafile /home/mypaper/go/src/hyperledger_project/test-network/organizations/ordererOrganizations/example.com/orderers/orderer.example.com/msp/tlscacerts/tlsca.example.com-cert.pem -C mychannel -n fabcar --peerAddresses localhost:7051 --tlsRootCertFiles /home/mypaper/go/src/hyperledger_project/test-network/organizations/peerOrganizations/org1.example.com/peers/peer0.org1.example.com/tls/ca.crt --peerAddresses localhost:9051 --tlsRootCertFiles /home/mypaper/go/src/hyperledger_project/test-network/organizations/peerOrganizations/org2.example.com/peers/peer0.org2.example.com/tls/ca.crt -c '{"Args":["InitLedger"]}'


elif [ $type == "BuyCar" ];
then

    peer chaincode invoke -o localhost:7050 --ordererTLSHostnameOverride orderer.example.com --tls --cafile /home/mypaper/go/src/hyperledger_project/test-network/organizations/ordererOrganizations/example.com/orderers/orderer.example.com/msp/tlscacerts/tlsca.example.com-cert.pem -C mychannel -n fabcar --peerAddresses localhost:7051 --tlsRootCertFiles /home/mypaper/go/src/hyperledger_project/test-network/organizations/peerOrganizations/org1.example.com/peers/peer0.org1.example.com/tls/ca.crt --peerAddresses localhost:9051 --tlsRootCertFiles /home/mypaper/go/src/hyperledger_project/test-network/organizations/peerOrganizations/org2.example.com/peers/peer0.org2.example.com/tls/ca.crt -c '{"function":"BuyCar","Args":["CAR0","999"]}'

elif [ $type == "QueryCarCouchDB" ];
then

    peer chaincode query -C mychannel -n fabcar -c '{"function":"QueryCarCouchDB","Args":["{\"selector\": {\"make\": \"Ford\"}}"]}'

elif [ $type == "QueryHistoryCar" ];
then

    peer chaincode query -C mychannel -n fabcar -c '{"function":"QueryHistoryCars","Args":["CAR0"]}'
elif [ $type == "GetBlock" ];
then

    peer channel fetch newest -c mychannel --ordererTLSHostnameOverride orderer.example.com
    configtxgen -inspectBlock mychannel_newest.block > mychannel_1.block.JSON

else
    echo "please select 'initLedger' OR 'BuyCar' OR 'QueryCar' OR 'QueryHistoryCar' OR 'GetBlock'"
fi    