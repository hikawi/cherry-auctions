<script setup lang="ts">
import { useScript } from "@/hooks/use-script";
import { ref } from "vue";
import { useI18n } from "vue-i18n";
import { useRouter } from "vue-router";

const { t } = useI18n({ useScope: "global" });

const name = ref("");
const email = ref("");
const password = ref("");
const confirmPassword = ref("");
const acceptPolicy = ref(false);

const loading = ref(false);
const error = ref("");

const router = useRouter();

useScript({
  src: `https://www.google.com/recaptcha/api.js?render=${import.meta.env.VITE_SITE_KEY}`,
  defer: true,
  id: "recaptcha-script",
});

async function register() {
  loading.value = true;
  error.value = "";

  if (!acceptPolicy.value) {
    loading.value = false;
    error.value = "register.accept_privacy_required";
    return;
  }

  // Setup recaptcha
  const captchaPromise = await grecaptcha.execute(import.meta.env.VITE_SITE_KEY, {
    action: "submit",
  });

  if (password.value != confirmPassword.value) {
    loading.value = false;
    error.value = "register.passwords_dont_match";
    return;
  }

  try {
    const res = await fetch(`${import.meta.env.VITE_API}/v1/auth/register`, {
      method: "POST",
      credentials: "include",
      headers: {
        "content-type": "application/json",
      },
      body: JSON.stringify({
        name: name.value,
        email: email.value,
        password: password.value,
        captcha_token: captchaPromise,
      }),
    });

    switch (res.status) {
      case 400:
        error.value = "register.invalid_request";
        break;
      case 409:
        error.value = "register.conflict";
        break;
      case 500:
        error.value = "register.internal_error";
        break;
      case 201:
        router.push({ path: "/login" });
        break;
    }
  } catch {
    error.value = "register.internet_error";
  }

  loading.value = false;
}
</script>

<template>
  <div class="flex w-full max-w-lg flex-col items-center gap-4 rounded-2xl p-6 shadow-xl">
    <h1 class="text-2xl font-bold">{{ t("register.title") }}</h1>

    <div class="flex w-full flex-col gap-2">
      <label class="flex w-full flex-col gap-1">
        {{ t("register.name") }}

        <input
          type="text"
          required
          v-model="name"
          :placeholder="t('register.name_placeholder')"
          class="hover:ring-claret-200 focus:ring-claret-600 w-full rounded-lg border border-zinc-300 px-4 py-2 duration-200 outline-none placeholder:text-black/50 hover:ring-2 focus:ring-2"
        />
      </label>

      <label class="flex w-full flex-col gap-1">
        {{ t("register.email") }}

        <input
          type="email"
          required
          v-model="email"
          :placeholder="t('register.email_placeholder')"
          class="hover:ring-claret-200 focus:ring-claret-600 w-full rounded-lg border border-zinc-300 px-4 py-2 duration-200 outline-none placeholder:text-black/50 hover:ring-2 focus:ring-2"
        />
      </label>

      <label class="flex w-full flex-col gap-1">
        {{ t("register.password") }}

        <input
          type="password"
          required
          v-model="password"
          class="hover:ring-claret-200 focus:ring-claret-600 w-full rounded-lg border border-zinc-300 px-4 py-2 duration-200 outline-none hover:ring-2 focus:ring-2"
        />
      </label>

      <label class="flex w-full flex-col gap-1">
        {{ t("register.confirm_password") }}

        <input
          type="password"
          required
          v-model="confirmPassword"
          class="hover:ring-claret-200 focus:ring-claret-600 w-full rounded-lg border border-zinc-300 px-4 py-2 duration-200 outline-none hover:ring-2 focus:ring-2"
        />
      </label>

      <label class="flex items-center gap-2">
        <input type="checkbox" v-model="acceptPolicy" class="h-4 w-4 rounded border-zinc-300" />
        <span class="text-sm">
          {{ t("register.accept_privacy") }}
          <a href="/privacy" class="ml-1 underline">{{ t("register.privacy_policy") }}</a>
        </span>
      </label>
    </div>

    <p
      v-if="error"
      class="bg-claret-100 border-claret-500 text-claret-700 w-full rounded-xl border-2 px-4 py-2"
    >
      {{ t(error) }}
    </p>

    <hr class="my-2 h-px w-full rounded-full border border-zinc-300" />

    <button
      @click="register"
      :disabled="loading || !acceptPolicy"
      class="bg-claret-600 disabled:bg-claret-700 border-claret-600 enabled:hover:text-claret-600 disabled:border-claret-700 w-full cursor-pointer rounded-xl border-2 p-2 py-3 text-white transition-all duration-200 hover:shadow-md enabled:hover:bg-transparent disabled:cursor-progress disabled:opacity-50"
    >
      {{ loading ? t("register.loading") : t("register.action") }}
    </button>
  </div>
</template>
