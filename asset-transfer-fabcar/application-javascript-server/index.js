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

async function main(amount) {
  var client = new BWAggregatorProto.Aggregator('localhost:9000',
                                       grpc.credentials.createInsecure());

  var data = {
    functionname: "BuyCar",
    key: "CAR0",
    fieldname: "Amount",
    operator: 1,
    operand: 1,
    precondition: 0,
    postcondition: amount
  }
  client.ReceiveBWTransaction(data, function(err, response) {
  
    console.log('return:', response.response);
  });
}
main(999);
main(998);
main(997);
main(996);
main(996);

