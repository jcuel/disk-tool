import { defineConfig } from "vite";

export default defineConfig(({ mode }) => ({
  base: mode === "demo" ? "/disk-tool/demo/" : "./",
  build: {
    outDir: mode === "demo" ? "../site/dist/demo" : "dist",
    emptyOutDir: mode === "demo" ? false : true,
  },
  define:
    mode === "demo"
      ? { "import.meta.env.VITE_DEMO_MODE": JSON.stringify("true") }
      : {},
}));
