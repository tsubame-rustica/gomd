import { defineConfig, loadEnv } from 'vite'
import react from '@vitejs/plugin-react'
import tailwindcss from '@tailwindcss/vite'

export default defineConfig(({ mode }) => {
  const env = loadEnv(mode, process.cwd(), '')

  return {
    plugins: [react(), tailwindcss()],
    server: {
      host: '0.0.0.0',
      port: 5173,
      proxy: {
        '/api': {
          target: env.VITE_API_TARGET || 'http://backend:8080',
          changeOrigin: true,
          secure: false,
        },
        // マークダウン内の画像（相対パスなど）のルーティングをバックエンドへ流す
        '^/.*\\.(png|jpe?g|gif|svg|webp|ico)$': {
          target: env.VITE_API_TARGET || 'http://backend:8080',
          changeOrigin: true,
          rewrite: (path) => `/api/contents${path}`,
        }
      },
    },
  }
})