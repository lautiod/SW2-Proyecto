import { defineConfig } from 'vite'
import react from '@vitejs/plugin-react'

// https://vitejs.dev/config/
export default defineConfig({
  plugins: [react()],
  server: {
    host: '0.0.0.0',  // Esto expone el servidor en todas las interfaces de red
    port: 5173,        // Puerto en el que está escuchando
  },
})
