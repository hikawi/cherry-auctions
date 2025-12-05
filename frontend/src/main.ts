import { createApp } from "vue";
import { createPinia } from "pinia";

import App from "./App.vue";
import "./styles.css";
import router from "./router";
import { createI18n } from "vue-i18n";

const app = createApp(App);

const i18n = createI18n({
  locale: "en-US",
  availableLocales: ["en-US", "ja-JP"],
  fallbackLocale: "en-US",
  formatFallbackMessages: true,
  messages: {},
});

app.use(createPinia());
app.use(i18n);
app.use(router);

app.mount("#app");
