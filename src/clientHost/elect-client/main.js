/*##################################*\
  Required Libraries
\*##################################*/
const {ipcMain, app, BrowserWindow, clipboard, dialog} = require('electron');
const {sanatizeText, stringToByteArry, ByteArryTostring, fragmentStreamlined} = require('./modules/utilityFunctions.js');
const {chatHistory} = require('./modules/chatHistoryClass.js')
const fs = require('fs')
const grpcLibrary = require('@grpc/grpc-js');
const protoLoader = require('@grpc/proto-loader');
const rsaLib = require('node-rsa');


/*##################################*\
  Internal data
\*##################################*/
let config = null;
let otherUser = null
let addressBook = []
let chatHistories = {}
let UIView = null
let outbound = null;
let cConnection = null;

let chat = null;
let timer = null;

let sessionToken = null;
let rsaKey = null;
let rsaHostKey = null;
let configFilePath = "."


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
  start();

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
  hostfunc_SendText(chatText);
})

ipcMain.on('deleteMSG', (event, msgID)=>{
  deleteMessage(msgID)

  //Send to Host Here
  hostfunc_DeleteMessage(msgID);

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
  //requestChat(id);
})

ipcMain.on('requestChatUName', (event, id)=>{
  requestChat(id);
})

ipcMain.on('updateUser', (event, name)=>{
  config.name = name;
  fs.writeFile(configFilePath+"/userData.json", JSON.stringify(config),(err)=>{})
})

ipcMain.on('connectToHost', (event, networkObject)=>{

  config.ip = networkObject.ip
  //config.port = networkObject.port
  fs.writeFile(configFilePath+"/userData.json", JSON.stringify(config),(err)=>{})
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
      return config

    case 'activeChat':
      return chat

    case 'connection':
      return 2    
      
    case 'userData':
      return config

    case 'otherClientName':
      return otherUser.name || "U. N. Owen"

}}
function fragmentRouter(fragmentName){

  switch(fragmentName){

    case 'userData':
      return fragmentStreamlined("chat/userData.html", config)

    case 'otherUserData':
      return fragmentStreamlined("chat/otherUserData.html", otherUser)

    case 'hostData':
      return fragmentStreamlined("chat/hostData.html", config)

    case 'hostConnectionMenu':
      return fragmentStreamlined("menu/hostConnection.html", config)

    case 'hostConnectionForm':
      return fragmentStreamlined("menu/hostConnectForm.html", config)
  
    case 'userInfoMenu':
      return fragmentStreamlined("menu/userInfoSection.html", config)

    case 'chatRequestForm':
        return fragmentStreamlined("menu/chatRequestForm.html", config)
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
    
}
function deleteMessage(msgID) {

  if(chat != null){
	  console.log(msgID);
    chat.removeMSG(msgID)
    chatHistories[otherUser.identifier] = chat
  }

  UIView.webContents.send('deleteMessage', msgID);

}
function RecieveText(text) {

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

  if (typeof(id) == "string"){

    loadChatData(id)

  } else {return;}


  //Start Connection
  //createNewHostConnection();

  if(chat == null){
    UIView.webContents.send('userConnection', {status:"failed", issue:"User Not found"})
    return;
  }  

  timer = setInterval(getNewMessages(id), 1000);

  //Load user data
  //loadChatData(id);

  //Move to other screen
  changeView('chat')

}




function loadChatData(id){

  getUser(id);

  let tempHist = hostfunc_GetConversation();
  
  //Replace user getting data from other source

  console.log(tempHist);

  let processedList = [];

  for (s of tempHist.convo){

    //Decrypt 
    processedList.push(s)

  }

  if(tempHist.convo != null){
    chat = new chatHistory(0, null, processedList, processedList.sort((a,b)=>{ return a < b})[0].order + 1)

  for (s of chatHistory.messages.length()){
    chatHistory.messages[s] = rsaHostKey.decrypt(chatHistory.messages[s]);
  }
    otherUser = {name:id, indentifier:"5672", publicKey:"", ip: "127.0.0.1", port:"9090", status: "Online"};

  } else {
    return
  }

  chatHistories[id] = chat

}

function getNewMessages(id){

  let tempHist = hostfunc_RecieveText()

  if(!tempHist){
    return;
  }

  for (x of tempHist.message){
    //Decrypt
    RecieveText(x);
  }
  
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

  configFilePath = process.argv[2] || configFilePath;

  if(!fs.statSync(configFilePath).isDirectory()){
    configFilePath = ".";
  }

  try {
    config = JSON.parse(fs.readFileSync(configFilePath+"/userData.json"));
  } catch (err){

  }

  config = {key: config.key ||"101",name: config.name ||"Default",ip: config.ip ||"localhost",port: config.port ||"9091", centralserverIP: config.centralserverIP || "localhost:8080"}
  
  fs.writeFileSync(configFilePath+"/userData.json", JSON.stringify(config),(err)=>{})

  rsaKey = new rsaLib({signingAlgorithm:'sha512'});

  fs.readFile(configFilePath+"/rsaKeyPrivate.pem",(err, data)=>{

    if(err){
      console.log(err)
      generateRSAKeys()

    } else {
      rsaKey.importKey(data, "pkcs8-private-pem");
      
    }
  });


  fs.readFile(configFilePath + "/client_public.pem", (err, data)=>{
    if(err){
      console.error(err);
      return;
    }
    rsaHostKey = new rsaLib(data, "pkcs8-public-pem", {signingAlgorithm:'sha512'});

  });

    //app.exit();
  


}

function generateRSAKeys(){

  rsaKey.generateKeyPair(2048, 65537);
  fs.writeFile(configFilePath+"/rsaKeyPrivate.pem", rsaKey.exportKey("pkcs8-private-pem"), (err)=>{});
  fs.writeFile(configFilePath+"/rsaKeyPublic.pem", rsaKey.exportKey("pkcs8-public-pem"), (err)=>{});

}


/*##################################*\
  grpc Functions
\*##################################*/

function createNewHostConnection(credentials){

  if(outbound != null)
    dissconnectFromHost();

  let networkAddr = "" + config.ip + ":" + config.port || "127.0.0.1:9091";

  try {
    let packageDefinition = protoLoader.loadSync('./modules/proto/host.proto', {
      keepCase: true,
      longs: String,
      enums: String,
      defaults: true,
      oneofs: true
    });

    /* Found an example of a chat application https://techblog.fexcofts.com/2018/07/20/grpc-nodejs-chat-example/*/
    const clientConstructor = grpcLibrary.loadPackageDefinition(packageDefinition).smvs.clientHost;
    outbound = new clientConstructor(networkAddr, grpcLibrary.credentials.createInsecure());

    packageDefinition = protoLoader.loadSync('./modules/proto/server.proto', {
      keepCase: true,
      longs: String,
      enums: String,
      defaults: true,
      oneofs: true
    });

    let cServer = new grpcLibrary.loadPackageDefinition(packageDefinition).smvs.Server;
    cConnection = new cServer(config.centralserverIP, grpcLibrary.credentials.createInsecure());
    Register();
    
    UIView.webContents.send('hostConnect');

  } catch (error) {
    console.log(error);
    UIView.webContents.send('userConnection', {status:"Failed", message:"host Login failed"});
    return;
  }

}

function Register(){

  if(!cConnection)
    return

  cConnection.Register({username: config.name, key:stringToByteArry(rsaKey.exportKey("pkcs8-public-pem")), ip:"localhost" },(err, responce)=>{
    console.log(err);
    console.log(responce);
  });

}
function getToken(){

  if(!cConnection)
    return

  cConnection.getToken({username:config.name},(err, responce)=>{
    if(err){
      console.log(err);
      return;
    }
    sessionToken = responce.authKey;

    console.log(responce);
  });
}

function UpdateIP(){

  if(!cConnection)
    return

  cConnection.UpdateIP({username:config.name, authKey:sessionToken, newIP:networkAddr},(err, responce)=>{
    if(err){
      console.log(err);
      return;
    }
    console.log(responce);
  });

}

function UpdateKey(){

  if(!cConnection)
    return
  
  //generateRSAKeys()
  cConnection.UpdateKey({username:config.name, authKey:sessionToken, newKey:stringToByteArry(rsaKey.exportKey("pkcs8-public-pem"))},(err, responce)=>{
    if(err){
      console.log(err);
      return;
    }
    console.log(responce);
  });

}

function getUser(uname){

  if(!cConnection)
    return

  cConnection.getUser({username:uname}, (err, responce)=>{

    if(err){
      console.log(err);
      return;
    }
    console.log(responce);

    otherUser.ip = responce.IP;
    otherUser.key = responce.publicKey;

  });

}

function hostfunc_LogIn(token, username, ip) {

  if(outbound == null)
    return;

  outbound.LogIn({token:[], username:config.name, ip:config.ip},(err, responce)=>{
    if(err){
      console.log(err)
      dissconnectFromHost()
    }
    console.log(responce);
  }) 

}; 
function hostfunc_UpdateKey(token, key) {

  outbound.UpdateKey({token:[], key:[]},(err, responce)=>{
    if(err){
      console.log(err)
      dissconnectFromHost()
    }
    console.log(responce);
  }) 


}; // Requests client's host to change keys
function hostfunc_PingUser(token, username) {

  if(outbound == null)
    return;

  outbound.PingUser({token:[], username:config.name, ip:ip},(err, responce)=>{
    if(err){
      console.log(err)
      dissconnectFromHost()
    }
    console.log(responce);
  }) 


};

// Messaging calls
function hostfunc_DeleteMessage(user, messageID, token) {

  if(outbound == null)
    return;
    
  outbound.DeleteMessage({user:otherUser.name, messageID:int(messageID), token:config.key}, function(err, responce){
    if(err){
      console.log(err)
      dissconnectFromHost()
    }
    console.log(responce);
  })
  
};

function hostfunc_SendText(targetUser, message, token) {

  if(outbound == null || rsaKey == null){
    return;
  }

  rsaKey.option
  key.setOptions({})
  rsaKey.sign(sanatizeText(chatText));

  outbound.SendText({username:config.name, text:rsaKey.encryptPrivate(sanatizeText(chatText)), token:rsaKey.sign(sanatizeText(chatText))}, function(err, responce){
    if(err){
      console.log(err)
      dissconnectFromHost()
    }
    console.log(responce);
  })

};

function hostfunc_RecieveText(listOfMessages, user, secret) {

  if(outbound == null || rsaKey == null){
    return;
  }
  
  outbound.RecieveText({token:config.key, username:id}, function(err, responce){
    console.log(responce);
    if(err){
      console.log(err)
      dissconnectFromHost()
    }

    for(m of responce.message){
      RecieveText(rsaHostKey.decryptPublic(m));
    }

  })

};

function hostfunc_GetConversation(token, username) {

  if(outbound == null){
    return
  }

  outbound.GetConversation({token:config.key, username:id}, function(err, responce){
    console.log(responce);
      
    if(err){
      console.log(err)
      dissconnectFromHost()
    }

  })
};

function hostfunc_Ping() {

  if(outbound == null){
    return
  }

  outbound.Ping({}, function(err, responce){
    console.log(responce);
      
    if(err){
      console.log(err)
      dissconnectFromHost()
    }

  })

};

function hostfunc_LogOut(ClientInfo) {

  if(outbound == null){
    return
  }

  outbound.LogOut({token:config.key, username:id}, function(err, responce){
    console.log(responce);
      
    if(err){
      console.log(err)
      dissconnectFromHost()
    }

  })
  
}; 