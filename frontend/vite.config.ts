import { defineConfig } from 'vite'
import { svelte } from '@sveltejs/vite-plugin-svelte'
import { fileURLToPath } from 'url'
import { dirname, resolve } from 'path'

const __dirname = dirname(fileURLToPath(import.meta.url))

export default defineConfig({
  plugins: [svelte()],
  resolve: {
    alias: {
      '@wailsjs': resolve(__dirname, './wailsjs'),
    },
  },
  build: {
    // Wails 需要将所有资源打包到单个目录
    outDir: 'dist',
    emptyOutDir: true,
    // 禁用代码分割，确保所有资源在一个文件中
    rollupOptions: {
      output: {
        manualChunks: undefined,
      },
    },
  },
  server: {
    // 开发服务器配置
    port: 5173,
    strictPort: true,
  },
})
