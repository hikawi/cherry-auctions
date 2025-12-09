import { describe, expect, test } from "vitest";
import { mount } from "@vue/test-utils";
import LoginForm from "../../components/login/LoginForm.vue";
import { createI18n } from "vue-i18n";

import enUS from "../../i18n/en-US.json";
import jaJP from "../../i18n/ja-JP.json";

describe("LoginForm", () => {
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
  });

  test("mounts correctly", async () => {
    const comp = mount(LoginForm, { global: { plugins: [i18n] } });
    expect(comp.text()).toContain("Login");
  });
});
