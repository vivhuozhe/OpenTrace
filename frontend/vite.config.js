import { defineConfig } from 'vite'
import react from '@vitejs/plugin-react'

export default defineConfig({
  plugins: [react()],
  server: {
    port: 3000,
    proxy: {
      '/dashboard': {
        target: 'http://localhost:8080',
        changeOrigin: true
      },
      '/levels': {
        target: 'http://localhost:8080',
        changeOrigin: true
      },
      '/nodes': {
        target: 'http://localhost:8080',
        changeOrigin: true
      },
      '/edges': {
        target: 'http://localhost:8080',
        changeOrigin: true
      },
      '/route': {
        target: 'http://localhost:8080',
        changeOrigin: true
      },
      '/map': {
        target: 'http://localhost:8080',
        changeOrigin: true
      }
    }
  }
})
