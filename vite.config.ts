import { defineConfig } from "vite";

import react from "@vitejs/plugin-react";

export default defineConfig({
  plugins: [react()],
  build: {
    outDir: "static/scripts/islands",
    rollupOptions: {
      input: "./react/main.tsx",
      output: {
        entryFileNames: "islands.js",
        format: "iife",
      },
    },
  },
});
