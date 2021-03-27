exports.chatHistory = class chatHistory{

  constructor(id, lastID, users, messages){
    
    this.messages = messages || [ //{order:0, speaker:0, messageText:"", metadata:{}}  [
      {order:0, speaker:-1, messageText:"Hello", metadata:{}}, 
      {order:1, speaker:0, messageText:"Bonjour", metadata:{}},
      {order:2, speaker:-1, messageText:"...", metadata:{}}
    ]
    this.newID = lastID || this.messages.length || 0
    this.id = id || parseInt(Math.floor( Math.random() * Math.pow(2,42)))
    this.users = users || [] //{id:0, publickey:"", name:""}
    
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

  createNewMessageEntry(id, text, metadata){
    this.messages.push({order:this.newID++, speaker:id, messageText:text, metadata:metadata})
  }

  addMSGOther(identifier, text, matadata){
    this.createNewMessageEntry(id, text, metadata)
    return this.getLastMessage()
  }
  addMSGuser(msg, metadata){
    this.createNewMessageEntry(-1, msg, metadata)
    return this.getLastMessage()
  }

  getmsgList(){return this.messages}

  getLastMessage(){return this.messages[this.messages.length - 1]}

  removeMSG(id){
    let list = this.messages
    //list = []
    this.messages = list.filter((e)=>{return e.order != id})
  }

}