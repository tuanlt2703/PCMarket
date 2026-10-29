import { fileURLToPath, URL } from 'node:url'
import { defineConfig } from 'vite'
import vue from '@vitejs/plugin-vue'

export default defineConfig({
  plugins: [vue()],
  resolve: {
    alias: {
      '@': fileURLToPath(new URL('./src', import.meta.url)),
      '@adminlte': fileURLToPath(new URL('./admin-lte-v4.10.1/dist', import.meta.url)),
    },
  },
  server: {
    port: 5173,
  },
})
