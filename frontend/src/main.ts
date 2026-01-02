import { createApp } from "vue";
import { createPinia } from "pinia";

import App from "./App.vue";
import "./styles.css";
import router from "./router";
import { createI18n } from "vue-i18n";
import piniaPluginPersistedState from "pinia-plugin-persistedstate";

// I18n stuff
import enUS from "@/i18n/en-US.json";
import jaJP from "@/i18n/ja-JP.json";

import dayjs from "dayjs";
import relativeTime from "dayjs/plugin/relativeTime";
import localizedFormat from "dayjs/plugin/localizedFormat";
import "@vueup/vue-quill/dist/vue-quill.snow.css";

// Import core locales your app starts with
import "dayjs/locale/en";
import "dayjs/locale/ja";

// Extend Day.js once for the entire app life-cycle
dayjs.extend(relativeTime);
dayjs.extend(localizedFormat);

const app = createApp(App);

type I18nSchema = typeof enUS;
const i18n = createI18n<[I18nSchema], "en-US" | "ja-JP">({
  legacy: false,
  availableLocales: ["en-US", "ja-JP"],
  fallbackLocale: "en-US",
  formatFallbackMessages: true,
  messages: {
    "en-US": enUS,
    "ja-JP": jaJP,
  },
  numberFormats: {
    "en-US": {
      currency: { style: "currency", currency: "USD" },
      decimal: { style: "decimal", minimumFractionDigits: 2 },
    },
    "ja-JP": {
      currency: { style: "currency", currency: "USD" },
      compact: { notation: "compact" }, // Turns 10,000 into 1万
    },
  },
});

const pinia = createPinia();
pinia.use(piniaPluginPersistedState);

app.use(pinia);
app.use(i18n);
app.use(router);

app.mount("#app");
