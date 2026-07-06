import { defineConfig } from "vite";

export default defineConfig({
  base: "/disk-tool/",
  build: {
    outDir: "dist",
    emptyOutDir: true,
  },
});
