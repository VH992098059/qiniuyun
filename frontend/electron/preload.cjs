const { contextBridge, ipcRenderer } = require('electron')

contextBridge.exposeInMainWorld('electronAPI', {
  executeCommand: (text) => ipcRenderer.invoke('assistant:execute', { text }),
})