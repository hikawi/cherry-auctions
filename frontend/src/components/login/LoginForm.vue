<script setup lang="ts">
import { onMounted, onUnmounted, ref } from "vue";
import { RotateCcwKey } from "lucide-vue-next";
import { useI18n } from "vue-i18n";
import { useRouter } from "vue-router";
import { useTokenStore } from "@/stores/token";

const { t } = useI18n({ useScope: "global" });

const email = ref("");
const password = ref("");
const loading = ref(false);
const error = ref("");

const hoveringForgotPassword = ref(false);
const smallWindow = ref(false);

const router = useRouter();
const token = useTokenStore();

function setSmallWindow() {
  smallWindow.value = window.innerWidth < 640;
}

onMounted(() => {
  setSmallWindow();
  window.addEventListener("resize", setSmallWindow);
});
onUnmounted(() => window.removeEventListener("resize", setSmallWindow));

async function forgotPassword() {
  router.push("/forgot");
}

async function login() {
  loading.value = true;
  error.value = "";

  try {
    const res = await fetch(`${import.meta.env.VITE_API}/v1/auth/login`, {
      method: "POST",
      credentials: "include",
      headers: {
        "content-type": "application/json",
      },
      body: JSON.stringify({ email: email.value, password: password.value }),
    });

    switch (res.status) {
      case 404:
        error.value = t("login.account_does_not_exist");
        break;
      case 400:
        error.value = t("login.invalid_request");
        break;
      case 401:
        error.value = t("login.wrong_password");
        break;
      case 421:
        error.value = t("login.wrong_method");
        break;
      case 500:
        error.value = t("login.internal_error");
        break;
      case 200:
        const json = await res.json();
        token.setToken(json.access_token);
        router.push("/");
        break;
    }
  } catch {
    error.value = t("login.internet_error");
  }

  loading.value = false;
}
</script>

<template>
  <div class="flex w-full max-w-lg flex-col items-center gap-4 rounded-2xl p-6 shadow-xl">
    <h1 class="text-2xl font-bold">{{ t("login.title") }}</h1>

    <div class="flex w-full flex-col gap-2">
      <label class="flex w-full flex-col gap-1">
        {{ t("login.email") }}

        <input
          type="email"
          required
          v-model="email"
          :placeholder="t('login.email_placeholder')"
          class="hover:ring-claret-200 focus:ring-claret-600 w-full rounded-lg border border-zinc-300 px-4 py-2 duration-200 outline-none placeholder:text-black/50 hover:ring-2 focus:ring-2"
        />
      </label>

      <label class="flex w-full flex-col gap-1">
        {{ t("login.password") }}

        <input
          type="password"
          required
          v-model="password"
          class="hover:ring-claret-200 focus:ring-claret-600 w-full rounded-lg border border-zinc-300 px-4 py-2 duration-200 outline-none hover:ring-2 focus:ring-2"
        />
      </label>
    </div>

    <p
      v-if="error"
      class="bg-claret-100 border-claret-500 text-claret-700 w-full rounded-xl border-2 px-4 py-2"
    >
      {{ error }}
    </p>

    <hr class="my-2 h-px w-full rounded-full border border-zinc-300" />

    <div class="flex w-full flex-col gap-2 font-semibold ease-out sm:flex-row">
      <button
        @click="forgotPassword"
        @mouseenter="hoveringForgotPassword = true"
        @mouseleave="hoveringForgotPassword = false"
        class="peer flex min-w-fit cursor-pointer flex-row items-center justify-center gap-2 overflow-x-hidden rounded-xl border-2 border-zinc-300 p-2 py-3 text-black transition-all duration-200 hover:border-zinc-600 hover:bg-zinc-300 hover:shadow-md sm:flex-1 sm:grow hover:sm:grow-3"
      >
        <RotateCcwKey class="size-6" />

        {{ hoveringForgotPassword || smallWindow ? t("login.forgot_password") : "" }}
      </button>

      <button
        @click="login"
        :disabled="loading"
        class="bg-claret-600 disabled:bg-claret-700 border-claret-600 enabled:hover:text-claret-600 disabled:border-claret-700 cursor-pointer rounded-xl border-2 p-2 py-3 text-white transition-all duration-200 hover:shadow-md enabled:hover:bg-transparent disabled:cursor-progress disabled:opacity-50 sm:flex-1 sm:grow-3 peer-hover:sm:grow"
      >
        {{ loading ? t("login.loading") : t("login.action") }}
      </button>
    </div>
  </div>
</template>
