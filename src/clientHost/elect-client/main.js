/*##################################*\
  Required Libraries
\*##################################*/
const {ipcMain, app, BrowserWindow, clipboard, dialog} = require('electron');
const {sanatizeText, passwordCheckDebug, fragmentStreamlined} = require('./modules/utilityFunctions.js');
const {chatHistory} = require('./modules/chatHistoryClass.js')
const {clientCommunication} = require('./modules/clientCommuniction.js');
const {networkInformation} = require('./modules/networkInformation.js');
const grpcLibrary = require('@grpc/grpc-js');
const protoLoader = require('@grpc/proto-loader');


/*##################################*\
  Internal data
\*##################################*/
let userData = null
let hostConnectionData = null
let otherUser = null
let addressBook = []
let chatHistories = {}
let UIView = null
let hostConnection = null;
let chat = null
let communicationToken = "123asdasdasdasdas45678"


/*##################################*\
  Basic Electron Startup
\*##################################*/

function createWindow () {
  const win = new BrowserWindow({
    minWidth: 1280,
    minHeight: 720,
    width: 1280,
    height: 720,
    webPreferences: {
      nodeIntegration: true,
      devTools: true,
      contextIsolation: false
    }   
  })
  UIView = win;
  //win.setMenu(null)
  changeView('menu')
}

app.whenReady().then(createWindow)

app.on('window-all-closed', () => {
  if (process.platform !== 'darwin') {
    app.quit()
  }
})

app.on('activate', () => {
  if (BrowserWindow.getAllWindows().length === 0) {
    createWindow()
  }
})


/*##################################*\
  Client Responce Event Listeners
\*##################################*/

ipcMain.on('chatSent', (event, chatText)=>{
  SendText(chatText)
  //Send to Host Here
  if(outbound != null)
  outbound.DeleteMessage({targetUser:otherUser.name, message:[sanatizeText(chatText)], token:communicationToken});

})

ipcMain.on('deleteMSG', (event, msgID)=>{
  deleteMessage(msgID)
  //Send to Host Here

  if(outbound != null)
  outbound.DeleteMessage({user:userData.name, messageID:int(msgID), token:communicationToken});
})

ipcMain.on('loginAttempt',(event, loginData)=>{

  let loginStatus = connectToHost(loginData)

  switch(loginStatus){
    case 0:
      changeView('chat')
      break;

    case 1:
      UIView.webContents.send('loginStatus', false)
      break;
      
  }
  
})

ipcMain.on('copyMessageText', (event, msgID)=>{
  clipboard.writeText(chat.messages.filter((e)=>{return e.order == msgID})[0].messageText)
})

ipcMain.on('copyMessageJSON', (event, msgID)=>{
  clipboard.writeText(JSON.stringify(chat.messages.filter((e)=>{return e.order == msgID})[0]))
})

ipcMain.on('exportCoversationToJSON', (event)=>{

  dialog.showMessageBox(UIView, {
    message:"Test",
    title:"Hello World",
    type:"question"
  });

})

//Menu
ipcMain.on('requestChat', (event, id)=>{
  requestChat(id);
})



ipcMain.on('updateUser', (event, name)=>{
  userData.name = name;
})

ipcMain.on('connectToHost', (event, networkObject)=>{

  hostConnectionData.ip = networkObject.ip
  hostConnectionData.port = networkObject.port
  createNewHostConnection(null);

})


//Leave
ipcMain.on('endChat', (event)=>{
  chat = null;
  otherUser = null;

  //Disconnect Action Here()
  changeView('menu');
})

ipcMain.on('endHostConnection', (event)=>{

  //Disconnect Action Here()
  dissconnectFromHost();
  UIView.webContents.send('hostDisconnect');


})

/*##################################*\
  Internal Functions
\*##################################*/

ipcMain.on('moveTo', (event, destination)=>{
  changeView(destination)
})

ipcMain.on('fetch', (event, obj)=>{
  event.returnValue = fetchData(obj)
})

ipcMain.on('renderFrament', (event, fragmentName)=>{
  event.returnValue = fragmentRouter(fragmentName)
})

ipcMain.on('provideContactList', (event, searchValue, searchType)=>{
  event.returnValue = populateAddressBook(searchValue, searchType);
})


/*##################################*\
  Router Functions for the Application
\*##################################*/

//Function to change what view is loaded
function changeView(destination){

  switch(destination){

    case 'login':
      UIView.loadFile('./UI/HTML/login.html')
      break

    case 'chat':
      UIView.loadFile('./UI/HTML/chat.html');
      break    
/* 
    case 'chatHistory':
      UIView.loadFile('./UI/HTML/chatHistory.html');
      break    
      
    case 'blocklist':
      UIView.loadFile('./UI/HTML/blocklist.html');
      break
*/
    case 'menu':
      UIView.loadFile('./UI/HTML/mainMenu.html');
      break

    default:
      changeView('login')
}}
function fetchData(requestOBJ){

  switch(requestOBJ){

    case 'connectionData':
      return hostConnectionData

    case 'activeChat':
      return chat

    case 'connection':
      return 2    
      
    case 'userData':
      return userData

    case 'otherClientName':
      return otherUser.name || "U. N. Owen"

}}
function fragmentRouter(fragmentName){

  switch(fragmentName){

    case 'userData':
      return fragmentStreamlined("chat/userData.html", userData)

    case 'otherUserData':
      return fragmentStreamlined("chat/otherUserData.html", otherUser)

    case 'hostData':
      return fragmentStreamlined("chat/hostData.html", hostConnectionData)

    case 'hostConnectionMenu':
      return fragmentStreamlined("menu/hostConnection.html", hostConnectionData)

    case 'hostConnectionForm':
      return fragmentStreamlined("menu/hostConnectForm.html", hostConnectionData)
  
    case 'userInfoMenu':
      return fragmentStreamlined("menu/userInfoSection.html", userData)

    case 'chatRequestForm':
        return fragmentStreamlined("menu/chatRequestForm.html", userData)
  }

  return ""
}

/*##################################*\
  Copy of the Functions Found in host.go
\*##################################*/

/* Events for Send Texts */
function SendText(chatText) {

    console.log("Client Sent: " + chatText.payload)
    
    if(chat != null){
      let msgData = chat.addMSGOther(-1, sanatizeText(chatText.payload), null)
      chatHistories[otherUser.identifier] = chat
      UIView.webContents.send('inBoundChat', msgData)
    }
    
    if(outbound != null){
      outbound.send({username:userData.name, text:sanatizeText(chatText)});
    }
    
}
function deleteMessage(msgID) {

  if(chat != null){
	  console.log(msgID);
    chat.removeMSG(msgID)
    chatHistories[otherUser.identifier] = chat
  }

  UIView.webContents.send('deleteMessage', msgID);

}
function RecieveText(text, identifier) {

  console.log("Other Sent: " + chatText.payload)

  if(chat != null){
    let msgData = chat.addMSGOther(1, sanatizeText(text), null)
    UIView.webContents.send('inBoundChat', msgData)
  }
}

function requestChat(id){

  if(id == null){
    return;
  }

  //Connect to the other user
  if(chat != null || otherUser != null){
    UIView.webContents.send('userConnection', {status:"failed", issue:"Chat in progress"})
    return;
  }

  if(outbound != null){
    return;
  }

  if(typeof(id) == "number"){
    
    if(id !== -1){
      otherUser = addressBook.filter((value)=>{return value.indentifier == id})[0] || null;
    }

    if(otherUser != null){
      loadChatData(otherUser.name)
    }


  }else if (typeof(id) == "string"){

    loadChatData(id)

  } else {return;}

  //Start Connection
  //createNewHostConnection();

  if(false){
    UIView.webContents.send('userConnection', {status:"failed", issue:"User Not found"})
    return;
  }  

  //Load user data
  //loadChatData(id);

  //Move to other screen
  changeView('chat')



}

function confirmChat(){

  let responce = 1;

  if(chat != null || otherUser != null){
    responce = dialog.showMessageBoxSync(UIView, {
        message:"XXXX would like to chat with you.",
        type:"info",
        buttons:["Accept", "Deny"],
        defaultID:1,
        title:"Chat Request",
        cancelId:1
      });
  }

  //Deny Selected
  if(responce){

    //

  //Accept Selected
  } else {


  }

  return;

}

function loadChatData(id){

  let tempHist = {}

  //Replace user getting data from other source
  if(outbound != null){
    tempHist = outbound.GetConversation({token:communicationToken, username:otherUser.name})

  } else {

    //Load Chat history (Fallback if Hisory is not found)
    //Create if does not exist
    tempHist = chatHistories[id]
  }

  if(tempHist != null){
    chat = new chatHistory(id, tempHist.speakers, tempHist.messages, tempHist.newID)
  } else {
    chat = new chatHistory(id)
  }

  chatHistories[id] = chat

}


/* Events for Conversations*/
function InitializeConvo() {

}

function ConfirmConvo() {

}

function GetConversation() {

}

/* Host connection events */

/*##################################*\
  Functions listed in diagram 4 in purposal
\*##################################*/

function registerWithHost(){

}

function updateIPaddress(){

}

function updateKey(){

}

function getUserIP(){

}

/* Extra Fucntion related to the others */
function dissconnectFromHost(){
  if(outbound != null){
    outbound.close();
    outbound = null
  }
}

function dissconnectFromOtherChat(){
}

function populateAddressBook(searchValue, searchType){

  let innerHTML = ""

  let contacts = []

  if(searchType == 0){
    contacts = addressBook.filter((value)=>{return value.name.toLowerCase().indexOf(searchValue.toLowerCase()) != -1})
  } else {
    contacts = addressBook.filter((value)=>{return value.publicKey.toLowerCase().indexOf(searchValue.toLowerCase()) == 0})
  }

  for (c of contacts){

    let tempValue = {name:c.name, indentifier:c.indentifier, publicKey:c.publicKey, IP:c.IP, status:"Online"}
    innerHTML += fragmentStreamlined('menu/contactListing.html', tempValue)

  }

  return innerHTML

}


/*##################################*\
  Start Function
\*##################################*/
function start(){

  userData = {key:"12345678", name:"Bob"}
  //hostConnectionData = new networkInformation("127.0.0.1", "9090", "user")
  //otherUser = {key:"12345678", name:"Alice", IP:"",status:""}

  //Note public keys here are hashes (Just here for debugging)
  /*
  addressBook = [
    {name:"Alice", indentifier:"5672", publicKey:"64489c85dc2fe0787b85cd87214b3810", ip: "127.0.0.1", port:"9090", status: "Online"},
    {name:"Bob", indentifier:"7036", publicKey:"2fc1c0beb992cd7096975cfebf9d5c3b", ip: "127.0.0.1", port:"9090", status: "Online"},
    {name:"Eve", indentifier:"4962", publicKey:"d3f791f59cbeff0ec06afb94bb23e772", ip: "127.0.0.1", port:"9090", status: "Online"},
    {name:"Marvin", indentifier:"4085", publicKey:"7db16a4ce881aecec2bfeb3e0c741888", ip: "127.0.0.1", port:"9090", status: "Online"},
    {name:"Oscar", indentifier:"8754", publicKey:"48a0572e6e7cfc81b428b18da87cf613", ip: "127.0.0.1", port:"9090", status: "Online"},
    {name:"Peggy", indentifier:"9914", publicKey:"469a32447498e6238dab042c08098b98", ip: "127.0.0.1", port:"9090", status: "Online"},
    {name:"Victor", indentifier:"337", publicKey:"82233bce59652cf3cc0eb7a03f3109d1", ip: "127.0.0.1", port:"9090", status: "Online"},
    {name:"Trent", indentifier:"6200", publicKey:"a52f4256f1abed061d9cceee75907248", ip: "127.0.0.1", port:"9090", status: "Online"}
  ]
  chatHistories = {
      5672: {newID: 3, messages: [
        {order:0, speaker:-1, messageText:"Hello", metadata:{}}, 
        {order:1, speaker:0, messageText:"Bonjour", metadata:{}},
        {order:2, speaker:-1, messageText:"...", metadata:{}}
      ], speakers: [{speakerID:0, identifier:5672}]
  }}
  */
  hostConnectionData = {
    ip : "localhost",
    port : "9090"
  };

  //outbound = new clientCommunication();
  //chat = new chatHistory(101013731);
  //outbound.establishConnection(hostConnectionData)


}start();

/*##################################*\
  grpc Functions
\*##################################*/

function createNewHostConnection(credentials){

  if(outbound != null)
    dissconnectFromHost();

  let networkAddr = "" + hostConnectionData.ip + ":" + hostConnectionData.port || "127.0.0.1:9090";

  const packageDefinition = protoLoader.loadSync('./modules/proto/host.proto', {
    keepCase: true,
    longs: String,
    enums: String,
    defaults: true,
    oneofs: true
  });

  try {

    /* Found an example of a chat application https://techblog.fexcofts.com/2018/07/20/grpc-nodejs-chat-example/  

    */

    const clientConstructor = grpcLibrary.loadPackageDefinition(packageDefinition).smvs.clientHost;
    outbound = new clientConstructor(networkAddr, grpcLibrary.credentials.createInsecure());
    //hostConnection = outbound.GetConversation({token:communicationToken, username:userData.name})


    //console.log(hostConnection);

    /*
    let hostConnection = outbound.join({user:userData.name, text:JSON.stringify({password:"12345678"})})

    let test = outbound.GetConversation

    hostConnection.on('data', (data)=>{

      if(data.token != loginStatus){
        return
      }
     
      switch(data.action){

        case 'reciveMessage':
          RecieveText(data.text, 1)
          break;

        case 'deleteMessage':
          deleteMessage(data.text)
          break;

      }

    })

    hostConnection.on('status', ()=>{
    })

    hostConnection.on('error', ()=>{
      hostConnection = null;
    })

    hostConnection.on('close', ()=>{
      dissconnectFromHost();
    })  
    */

    UIView.webContents.send('hostConnect');

  } catch (error) {
    console.log(error);
    UIView.webContents.send('loginStatus', -1);
  }

}



