const { contextBridge, ipcRenderer } = require('electron')

contextBridge.exposeInMainWorld('electronAPI', {
  executeCommand: (text) => ipcRenderer.invoke('assistant:execute', { text }),
  minimize: () => ipcRenderer.invoke('window:minimize'),
  restore: () => ipcRenderer.invoke('window:restore'),
})