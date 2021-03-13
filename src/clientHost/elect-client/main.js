const { app, BrowserWindow } = require('electron')

let UIView = null 

function createWindow () {
  const win = new BrowserWindow({
  	//width: 1280,
  	//height: 720,
    minWidth: 1280,
    minHeight: 720,
      //resizable: false,
    webPreferences: {
      nodeIntegration: true,
          devTools: true,
          contextIsolation: false
  }   
  })

  UIView = win;

  //win.setMenu(null)
  changeView('login')
  
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
const {sanatizeText} = require('./modules/textSantization.js');
const fs = require('fs')

/*##################################*\
\*##################################*/

let connectionData = {"network": "127.0.0.1", "port": "9090", "username": "user"}

/*##################################*\
\*##################################*/


ipcMain.on('chatSent', (event, chatText)=>{
    console.log("Client Sent: " + JSON.parse(chatText).payload)
    //event.reply('inBoundChat', JSON.parse(chatText).payload)
    UIView.webContents.send('inBoundChat', sanatizeText(JSON.parse(chatText).payload))
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
\*##################################*/

/* Inter Application processes */

//Function to change what view is loaded
function changeView(destination){

  switch(destination){

    case 'login':
      UIView.loadFile('./UI/HTML/login.html')
      break

    case 'chat':
      UIView.loadFile('./UI/HTML/display.html');
      break

    default:
      changeView('login')

  }

}

function fetchData(requestOBJ){

  switch(requestOBJ){

    case 'pastLogin':
      return connectionData

    case 'connection':
      return 2

  }

}

/*##################################*\
Copy of the Functions Found in host.go
\*##################################*/


function DeleteMessage() {
	return nil, nil
}

function InitializeConvo() {
	return nil, nil
}

function ConfirmConvo() {
	return nil, nil
}

function SendText() {
	return nil, nil
}

function RecieveText() {
	return nil, nil
}

function GetConversation() {
	return nil, nil
}
