exports.networkInformation = class networkInformation{

    constructor(n, p, u){
        this.network = n || "127.0.0.1"; 
        this.port = p || "9090"; 
        this.username = u || "user";
    }

}