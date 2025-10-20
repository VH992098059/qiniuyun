const { app, BrowserWindow, ipcMain } = require('electron')
const path = require('path')

function createWindow() {
  const win = new BrowserWindow({
    width: 1100,
    height: 740,
    webPreferences: {
      preload: path.join(__dirname, 'preload.js'),
      contextIsolation: true,
      nodeIntegration: false,
    },
  })

  const devUrl = process.env.VITE_DEV_SERVER_URL || 'http://localhost:5173'
  if (!app.isPackaged) {
    win.loadURL(devUrl)
    win.webContents.openDevTools()
  } else {
    win.loadFile(path.join(__dirname, '../dist/index.html'))
  }
}

app.whenReady().then(() => {
  createWindow()

  app.on('activate', () => {
    if (BrowserWindow.getAllWindows().length === 0) createWindow()
  })
})

app.on('window-all-closed', () => {
  if (process.platform !== 'darwin') app.quit()
})

// 简单的能力路由模拟：根据文本指令返回回应
ipcMain.handle('assistant:execute', async (_event, payload) => {
  const { text } = payload || { text: '' }
  let response = ''
  if (/播放音乐|音乐|song|music/i.test(text)) {
    response = '好的，准备播放音乐（模拟）。'
  } else if (/写(一篇)?文章|写作|article|write/i.test(text)) {
    response = '好的，开始起草文章（模拟）。'
  } else if (/打开|启动|launch|open/i.test(text)) {
    response = '收到，尝试打开指定应用（模拟）。'
  } else if (/组合|workflow|自动化|自动/i.test(text)) {
    response = '已创建组合任务的计划（模拟）。'
  } else if (text && text.trim().length) {
    response = `已收到指令：“${text}”，将尝试执行（模拟）。`
  } else {
    response = '请说出你的指令，例如：播放音乐、写一篇文章等。'
  }
  return { ok: true, reply: response }
})