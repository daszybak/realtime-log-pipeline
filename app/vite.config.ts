import { tanstackRouter } from "@tanstack/router-plugin/vite";
import react from "@vitejs/plugin-react";
import { defineConfig } from "vite";
import tsConfigPaths from "vite-tsconfig-paths";
import { visualizer } from 'rollup-plugin-visualizer';
import { fileURLToPath } from 'node:url';
import { dirname, resolve as resolvePath } from 'node:path';

const __filename = fileURLToPath(import.meta.url);
const __dirname = dirname(__filename);

export default defineConfig({
  appType: "spa",
  server: {
    proxy: {
      '/api': {
        target: 'http://localhost:8081',
        changeOrigin: true,
      },
    },
  },
  plugins: [
    tsConfigPaths(),
    tanstackRouter({
      autoCodeSplitting: true,
      enableRouteGeneration: false,
    }),
    react(),
    visualizer({
      open: false,
      filename: 'dist/stats.html',
      gzipSize: true,
      brotliSize: true
    })
  ],
  build: {
    target: 'esnext',
    minify: 'esbuild',
    sourcemap: true,
    commonjsOptions: {
      include: [/node_modules/],
      transformMixedEsModules: true
    },
  },
  optimizeDeps: {
    force: true,
    esbuildOptions: {
      target: 'esnext'
    },
  },
  resolve: {
    alias: {
      'react-query-swagger': resolvePath(__dirname, 'node_modules/react-query-swagger'),
    },
    preserveSymlinks: true
  },
  define: {
    global: 'globalThis',
  }
});
