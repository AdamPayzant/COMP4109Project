/*##################################*\
  Required Libraries
\*##################################*/
const {ipcMain, app, BrowserWindow, clipboard, dialog} = require('electron');
const {sanatizeText, passwordCheckDebug, fragmentStreamlined} = require('./modules/utilityFunctions.js');
const {chatHistory} = require('./modules/chatHistoryClass.js')
const fs = require('fs')
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
let outbound = null;
let chat = null
let communicationToken = "123asdasdasdasdas45678"
let timer = null;


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
    outbound.SendText({targetUser:otherUser.name, message:[sanatizeText(chatText)], token:communicationToken}, function(err, responce){
      if(err){
        console.log(err)
        dissconnectFromHost()
      }
      console.log(responce);
    });

})

ipcMain.on('deleteMSG', (event, msgID)=>{
  deleteMessage(msgID)
  //Send to Host Here

  if(outbound != null)
    outbound.DeleteMessage({user:userData.name, messageID:int(msgID), token:communicationToken}, function(err, responce){
      if(err){
        console.log(err)
        dissconnectFromHost()
      }
      console.log(responce);
    })
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

ipcMain.on('requestChatUName', (event, id)=>{
  requestChat(id);
})

ipcMain.on('updateUser', (event, name)=>{
  userData.name = name;
  fs.writeFileSync("./userData.json", JSON.stringify(userData))
})

ipcMain.on('connectToHost', (event, networkObject)=>{

  hostConnectionData.ip = networkObject.ip
  hostConnectionData.port = networkObject.port
  createNewHostConnection(null);

})


//Leave
ipcMain.on('endChat', (event)=>{
  timer = null;
  chat = null;
  otherUser = null;

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
      outbound.send({username:userData.name, text:sanatizeText(chatText)}, function(err, responce){
        if(err){
          console.log(err)
          dissconnectFromHost()
        }
        console.log(responce);
      })
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

  if(outbound === null){
    UIView.webContents.send('userConnection', {status:"failed", issue:"Host Not Connected"});
    return;
  }

  if(typeof(id) == "number"){
    
    if(id !== -1){
      otherUser = addressBook.filter((value)=>{return value.indentifier == id})[0] || null;
    }

    if(otherUser != null){
      loadChatData(otherUser.name)
    } 


  } else if (typeof(id) == "string"){

    loadChatData(id)

  } else {return;}


  //Start Connection
  //createNewHostConnection();

  if(chat == null){
    UIView.webContents.send('userConnection', {status:"failed", issue:"User Not found"})
    return;
  }  

  timer = setInterval(loadChatData(id), 1000);

  //Load user data
  //loadChatData(id);

  //Move to other screen
  changeView('chat')



}


function loadChatData(id){

  let tempHist = {}
  
  //Replace user getting data from other source
  if(outbound != null){
    tempHist = outbound.GetConversation({token:communicationToken, username:id}, function(err, responce){
      if(err){
        console.log(err)
        dissconnectFromHost()
      }
      console.log(responce);
    })
  }

  console.log(tempHist);

  if(tempHist.convo != null){
    chat = new chatHistory(0, null, tempHist.convo, tempHist.convo.sort((a,b)=>{ return a < b})[0].order + 1)
    otherUser = {name:id, indentifier:"5672", publicKey:"", ip: "127.0.0.1", port:"9090", status: "Online"};
  } else {
    return
  }

  chatHistories[id] = chat

}

/* Extra Fucntion related to the others */
function dissconnectFromHost(){
  if(outbound != null){
    timer = null;
    outbound.close();
    outbound = null
  }
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

  userData = JSON.parse(fs.readFileSync("./userData.json"))

  hostConnectionData = {
    ip : "localhost",
    port : "9090"
  };




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

    /* Found an example of a chat application https://techblog.fexcofts.com/2018/07/20/grpc-nodejs-chat-example/*/

    const clientConstructor = grpcLibrary.loadPackageDefinition(packageDefinition).smvs.clientHost;
    outbound = new clientConstructor(networkAddr, grpcLibrary.credentials.createInsecure());
    outbound.ReKey({token:communicationToken}, function(err, responce){
      if(err){
        console.log(err)
        dissconnectFromHost()
        throw err
      }
      console.log(responce);
    });

  } catch (error) {
    console.log(error);
    UIView.webContents.send('userConnection', {status:"Failed", message:"host Login failed"});
    return;
  }

  UIView.webContents.send('hostConnect');

}



