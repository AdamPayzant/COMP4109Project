const grpcLibrary = require('@grpc/grpc-js');
const protoLoader = require('@grpc/proto-loader');

exports.clientCommunication =  class clientCommunication{


    constructor(){
        this.connection = null
    }

    //establishConnection(data) {}

    establishConnection(data) {

        let networkAddr = "" + data.network + ":" + data.port || "127.0.0.1:9090";

        //Drop existing connection is it exists
        if (this.connection != null){
            this.connection.close();
        }   

        //Load Client Definition
        const packageDefinition = protoLoader.loadSync('./modules/proto/client.proto', {
            keepCase: true,
            longs: String,
            enums: String,
            defaults: true,
            oneofs: true
        });

        //Create the connection
        const clientConstructor = grpcLibrary.loadPackageDefinition(packageDefinition).smvs;
        this.connection = new clientConstructor.client(networkAddr, grpcLibrary.credentials.createInsecure())

        //Section to bind handlers/functions

    }

}
