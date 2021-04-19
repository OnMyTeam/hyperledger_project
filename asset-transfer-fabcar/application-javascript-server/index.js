const path = require('path');
const grpc = require('grpc');
protoLoader = require('@grpc/proto-loader');

const PROTO_PATH = path.resolve(__dirname , '../../BWAggregator/protos/user.proto');

const packageDefinition = protoLoader.loadSync(
    PROTO_PATH,
    {keepCase: true,
     longs: String,
     enums: String,
     defaults: true,
     oneofs: true
    });
const helloProto = grpc.loadPackageDefinition(packageDefinition).user;

function main() {
  var client = new helloProto.User('localhost:9000',
                                       grpc.credentials.createInsecure());
  var user;
  if (process.argv.length >= 3) {
    user = process.argv[2];
  } else {
    user = 'world';
  }
  client.getUser({user_id: "1", name: "sanggi", phone_number:"010", age:11}, function(err, response) {
    console.log('Greeting:', response);
  });
}

main();