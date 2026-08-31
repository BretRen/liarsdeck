import { defineConfig } from 'vite';
import vue from '@vitejs/plugin-vue';
import tailwindcss from '@tailwindcss/vite';
import { resolve } from 'path';
import { execSync } from 'child_process';

let appVersion = 'v2.3.0';
try {
  appVersion = execSync('git describe --tags --always').toString().trim();
} catch (_) {}

export default defineConfig({
  plugins: [vue(), tailwindcss()],
  define: {
    __APP_VERSION__: JSON.stringify(appVersion),
  },
  build: {
    outDir: resolve(__dirname, '../public'),
    emptyOutDir: true,
  },
  server: {
    port: 5173,
    proxy: {
      '/ws': {
        target: 'ws://localhost:8095',
        ws: true,
        changeOrigin: true,
      },
    },
  },
});
