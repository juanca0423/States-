// frontend/vite.config.js
import { defineConfig } from "vite";
import react from "@vitejs/plugin-react";
import { resolve } from "path";

export default defineConfig({
  plugins: [react()],
  define: {
    "process.env.NODE_ENV": JSON.stringify(
      process.env.NODE_ENV || "production",
    ),
    "process.env": {},
  },
  build: {
    // Vite compilará el código y lo mandará directo a tus estáticos de Go
    outDir: resolve(__dirname, "../estaticos/js"),
    emptyOutDir: false,
    lib: {
      entry: resolve(__dirname, "src/main.jsx"),
      formats: ["iife"],
      name: "FormularioReact",
      fileName: () => "bundle-formulario.js",
    },
  },
});
