import { defineConfig } from 'vite';
import { svelte } from '@sveltejs/vite-plugin-svelte';

export default defineConfig({
  plugins: [svelte()],
  server: {
    port: 5173,
    allowedHosts: 'all',
    proxy: {
      '/api': { target: process.env.API_TARGET ?? 'http://localhost:3000' },
    },
  },
  build: {
    outDir: 'dist',
  },
});
