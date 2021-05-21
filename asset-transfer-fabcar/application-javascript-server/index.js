const path = require('path');
const grpc = require('grpc');
protoLoader = require('@grpc/proto-loader');

const PROTO_PATH = path.resolve(__dirname , '../../BWAggregator/protos/bwaggregator.proto');
   
const packageDefinition = protoLoader.loadSync(
    PROTO_PATH,
    {keepCase: true,
     longs: String,
     enums: String,
     defaults: true,
     oneofs: true
    });

const BWAggregatorProto = grpc.loadPackageDefinition(packageDefinition).protos;

function main(key, fieldname, amount) {
  var client = new BWAggregatorProto.Aggregator('localhost:9000',
                                       grpc.credentials.createInsecure());

  var data = {
    functionname: "BuyCarAfter",
    key: key,
    fieldname: fieldname,
    operator: 0,
    operand: 1,
    precondition: 0,
    postcondition: amount
  }
  client.ReceiveBWTransaction(data, function(err, response) {
  
    console.log('return:', response.response, response.payload.toString());
  });
}

for( var i=0; i <=process.argv[2]; i++){
  main("CAR0","amount", 1000);
}



