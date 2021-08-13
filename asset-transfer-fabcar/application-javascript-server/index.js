const path = require('path');
const grpc = require('grpc');
protoLoader = require('@grpc/proto-loader');

const PROTO_PATH = path.resolve(__dirname , '../../TxAggregator/protos/txaggregator.proto');
   
const packageDefinition = protoLoader.loadSync(
    PROTO_PATH,
    {keepCase: true,
     longs: String,
     enums: String,
     defaults: true,
     oneofs: true
    });

const TxAggregatorProto = grpc.loadPackageDefinition(packageDefinition).protos;

function main() {
  var client = new TxAggregatorProto.Aggregator('localhost:9000',
                                       grpc.credentials.createInsecure());

  var data = {
    functionname: "BuyCarAfter",
    key: "CAR0",
    fieldname: "amount",
    operator: 0,
    operand: 1,
    precondition: 0,
    postcondition: 10000
  }
  var start = Date.now()
  console.log('/s/', start);
  client.ReceiveTaggedTransaction(data, function(err, response) {
    // console.log(num, '/e/', (Date.now() - start)/1000,'/',response.payload.toString());
    console.log('/e/', Date.now(),'/',response.payload.toString());
    // console.log('return:', process.argv[2], process.argv[3], response.response, response.payload.toString());
  });
}
for(var i=0; i<50; i++){
  main();
}




