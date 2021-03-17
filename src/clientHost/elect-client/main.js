/*##################################*\
  Required Libraries
\*##################################*/
const {ipcMain, app, BrowserWindow } = require('electron');
const {sanatizeText, passwordCheckDebug} = require('./modules/utilityFunctions.js');
const {chatHistory} = require('./modules/chatHistoryClass.js')
const fs = require('fs');


/*##################################*\
  Internal data
\*##################################*/

let UIView = null 
let connectionData = {"network": "127.0.0.1", "port": "9090", "username": "user"}
let chat = new chatHistory(101013731);


/*##################################*\
  Basic Electron Startup
\*##################################*/

function createWindow () {
  const win = new BrowserWindow({
    minWidth: 1280,
    minHeight: 720,
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


/*##################################*\
  Internal Functions
\*##################################*/

ipcMain.on('moveTo', (event, destination)=>{
  changeView(destination)
})

ipcMain.on('fetch', (event, obj)=>{
  event.returnValue = fetchData(obj)
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
      return connectionData

    case 'activeChat':
      return chat

    case 'connection':
      return 2
}}


/*##################################*\
  Copy of the Functions Found in host.go
\*##################################*/

/* Events for Send Texts */
function SendText(chatText) {

    console.log("Client Sent: " + JSON.parse(chatText).payload)
    
    if(chat != null){
      let msgData = chat.addMSGuser(sanatizeText(JSON.parse(chatText).payload), null)
      UIView.webContents.send('inBoundChat', msgData)
    }

}
function DeleteMessage(msgID) {

  if(chat != null)
	  chat.removeMSG(msgID)

}
function RecieveText() {
	return null, null
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
function connectToHost(loginText){
  
  let attemptData = JSON.parse(JSON.parse(loginText).payload)
  //console.log("Login Sent: " + JSON.stringify(attemptData))
  
  //Update local connection data
  connectionData.network =  attemptData.network || connectionData.network
  connectionData.port =     attemptData.port ||  connectionData.port
  connectionData.username = attemptData.username ||  connectionData.username

  //This is for debugging purposes (Switch to actual login proceedure)
  if (passwordCheckDebug(attemptData.password)){
    return 0
  }

  return 1

  /**
   * Two stages of this process (impliment latter):
   * 
   *  1. Make sure that the client can connect
   * 
   *  2. Attempt provide credentials
   * 
   *  Inform the User of the success of these two steps
   * 
   */

}
function dissconnectFromHost(){

}
