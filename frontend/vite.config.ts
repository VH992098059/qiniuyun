import { defineConfig } from 'vite'
import vue from '@vitejs/plugin-vue'

// https://vite.dev/config/
export default defineConfig({
  plugins: [vue()],
  server:{
    proxy:{
      '/ws': {
        target: 'ws://localhost:8000', // 需要跨域代理的本地路径
        ws: true,
        rewrite: (path) => path.replace(/^\/*/, ''),
      },
    }
  }
})
