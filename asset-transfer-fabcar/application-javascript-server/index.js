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

async function main(key, fieldname, amount) {
  var client = new BWAggregatorProto.Aggregator('localhost:9000',
                                       grpc.credentials.createInsecure());

  var data = {
    functionname: "BuyCar",
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
main("CAR0","AA",3);
main("CAR0","BB",3);
main("CAR0","CC",3);
main("CAR0","DD",3);
main("CAR0","EE",3);
main("CAR1","FF",3);
main("CAR1","GG",3);
main("CAR1","HH",3);
main("CAR1","II",3);
main("CAR1","JJ",3);

