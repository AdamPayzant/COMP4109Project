const { app, BrowserWindow } = require('electron')

function createWindow () {
  	const win = new BrowserWindow({
    	width: 1280,
    	height: 720,
        //minWidth: 1280,
        //minHeight: 720,
        //resizable: false,
	webPreferences: {
		nodeIntegration: true,
        devTools: true,
        contextIsolation: false
	}   

})
    win.setMenu(null)
   win.loadFile('./UI/display.html')
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
const {dialog, ipcMain} = require('electron');
const fs = require('fs')

ipcMain.on('chatSent', (event, chatText)=>{
    console.log("Client Sent: " + JSON.parse(chatText).payload)
    event.reply('inBoundChat', chatText)

})

ipcMain.on('loginAttempt',(event, loginText)=>{
    console.log("Login Sent: " + JSON.parse(loginText).payload)
    event.returnValue = (JSON.parse(loginText).payload == "12345678") 

})


