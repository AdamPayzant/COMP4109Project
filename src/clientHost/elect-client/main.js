/*##################################*\
  Required Libraries
\*##################################*/
const {ipcMain, app, BrowserWindow, clipboard} = require('electron');
const {sanatizeText, passwordCheckDebug, fragmentStreamlined} = require('./modules/utilityFunctions.js');
const {chatHistory} = require('./modules/chatHistoryClass.js')
const {clientCommunication} = require('./modules/clientCommuniction.js');
const {networkInformation} = require('./modules/networkInformation.js');


/*##################################*\
  Internal data
\*##################################*/
let userData = {key:"12345678", name:"Bob"}
let hostConnectionData = new networkInformation("127.0.0.1", "9090", "user")
let otherUser = {key:"12345678", name:"Alice", IP:"",status:""}

let addressBook = [
  {name:"Alice", indentifier:"101013731", publicKey:"b6eeae78488755fe12bf9ea1028882fd", IP: "127.0.0.1"},
  {name:"Alice", indentifier:"101013731", publicKey:"b6eeae78488755fe12bf9ea1028882fd", IP: "127.0.0.1"},
  {name:"Alice", indentifier:"101013731", publicKey:"b6eeae78488755fe12bf9ea1028882fd", IP: "127.0.0.1"},
  {name:"Alice", indentifier:"101013731", publicKey:"b6eeae78488755fe12bf9ea1028882fd", IP: "127.0.0.1"},
  {name:"Alice", indentifier:"101013731", publicKey:"b6eeae78488755fe12bf9ea1028882fd", IP: "127.0.0.1"},
  {name:"Alice", indentifier:"101013731", publicKey:"b6eeae78488755fe12bf9ea1028882fd", IP: "127.0.0.1"},
  {name:"Alice", indentifier:"101013731", publicKey:"b6eeae78488755fe12bf9ea1028882fd", IP: "127.0.0.1"},
  {name:"Alice", indentifier:"101013731", publicKey:"b6eeae78488755fe12bf9ea1028882fd", IP: "127.0.0.1"},
  {name:"Alice", indentifier:"101013731", publicKey:"b6eeae78488755fe12bf9ea1028882fd", IP: "127.0.0.1"},
  {name:"Alice", indentifier:"101013731", publicKey:"b6eeae78488755fe12bf9ea1028882fd", IP: "127.0.0.1"},
  {name:"Alice", indentifier:"101013731", publicKey:"b6eeae78488755fe12bf9ea1028882fd", IP: "127.0.0.1"},
  {name:"Alice", indentifier:"101013731", publicKey:"b6eeae78488755fe12bf9ea1028882fd", IP: "127.0.0.1"}
]



let UIView = null
let outbound = new clientCommunication();
let chat = new chatHistory(101013731);
outbound.establishConnection(hostConnectionData)

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
  changeView('chat')
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
})

ipcMain.on('deleteMSG', (event, msgID)=>{
  DeleteMessage(msgID)
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

ipcMain.on('exportCoversationToJSON', ()=>{

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

  let innerHTML = ""

  let contacts = []

  if(searchType == 0){
    contacts = addressBook.filter((value)=>{return value.name.indexOf(searchValue) != -1})
  } else {
    contacts = addressBook.filter((value)=>{return value.publicKey.indexOf(searchValue.lowercase()) == 0})
  }

  for (c of contacts){

    let tempValue = {name:c.name, indentifier:c.indentifier, publicKey:c.publicKey, IP:c.IP, status:"Online"}
    innerHTML += fragmentStreamlined('contactListing.html', tempValue)

  }

  event.returnValue = innerHTML

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

    case 'otherClientName':
      return otherUser.name || "U. N. Owen"

}}
function fragmentRouter(fragmentName){

  switch(fragmentName){

    case 'userData':
      return fragmentStreamlined("userData.html", userData)

    case 'otherUserData':
      return fragmentStreamlined("otherUserData.html", otherUser)

    case 'hostData':
      return fragmentStreamlined("hostData.html", {IP:"Bob", status:"Alive"})

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
      let msgData = chat.addMSGuser(sanatizeText(chatText.payload), null)
      UIView.webContents.send('inBoundChat', msgData)
    }

}
function DeleteMessage(msgID) {

  if(chat != null)
	  chat.removeMSG(msgID)

}
function RecieveText(text, idenitfier) {

  console.log("Other Sent: " + chatText.payload)

  if(chat != null){
    let msgData = chat.addMSGOther(idenitfier, sanatizeText(text), null)
    UIView.webContents.send('inBoundChat', msgData)
  }
}

/* Events for Conversations*/
function InitializeConvo() {
	return null, null
}
function ConfirmConvo() {
	return null, null
}
function GetConversation() {
	return null, null
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

}
