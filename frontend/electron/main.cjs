const { app, BrowserWindow, ipcMain, Menu } = require('electron')
const path = require('path')

let mainWindow = null
let previousBounds = null

function createWindow() {
  mainWindow = new BrowserWindow({
    width: 1280,
    height: 840,
    minWidth: 960,
    minHeight: 600,
    autoHideMenuBar: true,
    alwaysOnTop: false,
    webPreferences: {
      preload: path.join(__dirname, 'preload.cjs'),
      contextIsolation: true,
      nodeIntegration: false,
    },
  })

  // 移除应用菜单并隐藏菜单栏
  Menu.setApplicationMenu(null)
  mainWindow.setMenuBarVisibility(false)
  mainWindow.setAlwaysOnTop(false, 'normal')

  // 记录初始窗口位置与尺寸
  previousBounds = mainWindow.getBounds()
  mainWindow.on('move', () => {
    try { previousBounds = mainWindow.getBounds() } catch {}
  })
  mainWindow.on('resize', () => {
    try { previousBounds = mainWindow.getBounds() } catch {}
  })

  const devUrl = process.env.VITE_DEV_SERVER_URL || 'http://localhost:5173'
  if (!app.isPackaged) {
    mainWindow.loadURL(devUrl)
  } else {
    mainWindow.loadFile(path.join(__dirname, '../dist/index.html'))
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

// 渲染层请求：最小化窗口并记录当前位置
ipcMain.handle('window:minimize', () => {
  if (!mainWindow) return false
  try {
    previousBounds = mainWindow.getBounds()
    mainWindow.minimize()
    return true
  } catch {
    return false
  }
})

// 渲染层请求：还原窗口到之前的位置与大小
ipcMain.handle('window:restore', () => {
  if (!mainWindow) return false
  try {
    if (mainWindow.isMinimized()) mainWindow.restore()
    if (previousBounds) mainWindow.setBounds(previousBounds, true)
    mainWindow.show()
    mainWindow.focus()
    return true
  } catch {
    return false
  }
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