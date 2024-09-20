import path from "path";
import react from "@vitejs/plugin-react";
import svgr from "vite-plugin-svgr";
import { defineConfig } from "vite";
// https://vitejs.dev/config/
export default defineConfig({
  server: {
    port: 3000,
    fs: {
      strict: false,
    },
  },
  plugins: [svgr(), react()],
  resolve: {
    alias: {
      "@powerpipe": "/src",
    },
  },
});
