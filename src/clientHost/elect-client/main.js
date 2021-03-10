const { app, BrowserWindow } = require('electron')

let UIView = null 

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
    //win.setMenu(null)
   win.loadFile('./UI/HTML/display.html')

   UIView = win;
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

ipcMain.on('chatSent', (event, chatText)=>{
    console.log("Client Sent: " + JSON.parse(chatText).payload)
    //event.reply('inBoundChat', JSON.parse(chatText).payload)
    UIView.webContents.send('inBoundChat', sanatizeText(JSON.parse(chatText).payload))
})

ipcMain.on('loginAttempt',(event, loginText)=>{
    console.log("Login Sent: " + JSON.parse(loginText).payload)


    UIView.webContents.send('loginStatus', (JSON.parse(JSON.parse(loginText).payload).password == "12345678"))


    
    //event.returnValue = (JSON.parse(JSON.parse(loginText).payload).password == "12345678") 

})




