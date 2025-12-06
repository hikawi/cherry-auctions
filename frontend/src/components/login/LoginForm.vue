<script setup lang="ts">
import { onMounted, onUnmounted, ref } from "vue";
import { RotateCcwKey } from "lucide-vue-next";
import { useI18n } from "vue-i18n";

const { t } = useI18n({ useScope: "global" });

const email = ref("");
const password = ref("");
const loading = ref(false);
const hoveringForgotPassword = ref(false);
const smallWindow = ref(false);

function setSmallWindow() {
  smallWindow.value = window.innerWidth < 640;
}

onMounted(() => {
  setSmallWindow();
  window.addEventListener("resize", setSmallWindow);
});
onUnmounted(() => window.removeEventListener("resize", setSmallWindow));

async function forgotPassword() {}

async function login() {
  loading.value = true;

  try {
    // TODO:
    // Logic
    await new Promise((resolve) => setTimeout(resolve, 2000));
  } catch {}

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
