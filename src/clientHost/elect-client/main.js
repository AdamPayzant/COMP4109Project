
/*##################################*\
  Basic Electron Startup
\*##################################*/

const { app, BrowserWindow } = require('electron')

let UIView = null 

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

//For inter-application communication
let fileName = null;
const {dialog, ipcMain, webContents} = require('electron');
const {sanatizeText} = require('./modules/utilityFunctions.js');

const {chatHistory} = require('./modules/chatHistoryClass.js')

const fs = require('fs')

/*##################################*\
  Internal data
\*##################################*/

let connectionData = {"network": "127.0.0.1", "port": "9090", "username": "user"}
let chat = new chatHistory(117013);

/*##################################*\
  Client Responce Event Listeners
\*##################################*/


ipcMain.on('chatSent', (event, chatText)=>{
    console.log("Client Sent: " + JSON.parse(chatText).payload)
    //event.reply('inBoundChat', JSON.parse(chatText).payload)

    let msgData = {speaker:-1, order:chatHistory.newID, messageText:sanatizeText(JSON.parse(chatText).payload)}
    chat.addMSG(msgData)
    UIView.webContents.send('inBoundChat', msgData)
})

ipcMain.on('deleteMSG', (event, msgID)=>{
  chat.removeMSG(msgID)
})

ipcMain.on('loginAttempt',(event, loginText)=>{

    console.log("Login Sent: " + JSON.parse(loginText).payload)
    let attemptData = JSON.parse(JSON.parse(loginText).payload)

    //Update local connection data
    connectionData.network =  attemptData.network || connectionData.network
    connectionData.port =     attemptData.port ||  connectionData.port
    connectionData.username = attemptData.username ||  connectionData.username

    let loginAttempt = (attemptData.password == "12345678")

    if (loginAttempt){  
      changeView('chat')
      return
    }

    UIView.webContents.send('loginStatus', loginAttempt)

    //event.returnValue = (JSON.parse(JSON.parse(loginText).payload).password == "12345678") 
})

ipcMain.on('moveTo', (event, destination)=>{

  changeView(destination)

})

ipcMain.on('fetch', (event, obj)=>{

  event.returnValue = fetchData(obj)

})


/*##################################*\
  Inter Application processes
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

  }

}

function fetchData(requestOBJ){

  switch(requestOBJ){

    case 'connectionData':
      return connectionData

    case 'activeChat':
      //console.log(JSON.stringify(chat))
      return chat

    case 'connection':
      return 2

  }

}

/*##################################*\
  Copy of the Functions Found in host.go
\*##################################*/

function SendText() {
	return null, null
}
function DeleteMessage() {
	return null, null
}
function RecieveText() {
	return null, null
}


function InitializeConvo() {
	return null, null
}

function ConfirmConvo() {
	return null, null
}

function GetConversation() {
	return null, null
}
