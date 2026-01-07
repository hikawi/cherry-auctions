import { fileURLToPath, URL } from "node:url";

import { defineConfig } from "vite";
import vue from "@vitejs/plugin-vue";
import sitemap from "vite-plugin-sitemap";
import tailwindcss from "@tailwindcss/vite";
import routes from "./src/router";

// https://vite.dev/config/
export default defineConfig({
  plugins: [
    vue(),
    tailwindcss(),
    sitemap({
      hostname: "https://cherry-auctions.luny.dev",
      dynamicRoutes: routes.getRoutes().map((r) => r.path),
    }),
  ],
  resolve: {
    alias: {
      "@": fileURLToPath(new URL("./src", import.meta.url)),
    },
  },
  server: {
    host: "0.0.0.0",
  },
});
