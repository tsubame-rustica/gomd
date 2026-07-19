import { defineConfig } from 'vite'
import react from '@vitejs/plugin-react'
import tailwindcss from '@tailwindcss/vite' // これを追加

export default defineConfig({
  plugins: [
    react(),
    tailwindcss(), // これを追加
  ],
})