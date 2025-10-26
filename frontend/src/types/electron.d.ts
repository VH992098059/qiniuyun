export {}

declare global {
  interface Window {
    electronAPI?: {
      executeCommand: (text: string) => Promise<{ ok: boolean; reply: string }>
    }
  }
}