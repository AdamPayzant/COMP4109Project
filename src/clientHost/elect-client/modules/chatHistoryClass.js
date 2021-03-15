exports.chatHistory = class chatHistory{

  constructor(id){
    
    this.messages = [ //{order:0, speaker:0, messageText:"", metadata:{}}  [
      {order:0, speaker:-1, messageText:"Hello", metadata:{}}, 
      {order:1, speaker:0, messageText:"Bonjour", metadata:{}},
      {order:2, speaker:-1, messageText:"...", metadata:{}}
    ]
    this.newID = this.messages.length || 0
    this.id = id || 0
    this.users = [] //{id:0, publickey:"", name:""}
    
  }

  getID() {
    return this.id
  }

  getUserNames() {
    return this.participants
  }

  getUserName(num) {
    for (x of this.users){
      if(x.id == num)
        return x.name
    }
    return ""
  }

  addMSG(msg){
    this.messages.push(msg)
    this.newID++
  }

  getmsgList(){return this.messages}

  removeMSG(id){
    let list = this.messages
    //list = []
    this.messages = list.filter((e)=>{return e.order != id})
  }

}
