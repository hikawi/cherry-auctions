<script setup lang="ts">
import { ref } from "vue";
import { Mail } from "lucide-vue-next";
import { useI18n } from "vue-i18n";

const { t } = useI18n({ useScope: "global" });

const email = ref("");
const loading = ref(false);
const error = ref("");
const success = ref("");

function validEmail(e: string) {
  return /^[^\s@]+@[^\s@]+\.[^\s@]+$/.test(e);
}

async function submit() {
  loading.value = true;
  error.value = "";
  success.value = "";

  if (!validEmail(email.value)) {
    loading.value = false;
    error.value = "forgot.invalid_request";
    return;
  }

  try {
    const res = await fetch(`${import.meta.env.VITE_API}/v1/auth/forgot-password`, {
      method: "POST",
      credentials: "include",
      headers: { "content-type": "application/json" },
      body: JSON.stringify({ email: email.value }),
    });

    switch (res.status) {
      case 400:
        error.value = "forgot.invalid_request";
        break;
      case 500:
        error.value = "forgot.internal_error";
        break;
      case 200:
      case 204:
        success.value = "forgot.sent";
        break;
      default:
        error.value = "forgot.internal_error";
    }
  } catch {
    error.value = "forgot.internet_error";
  }

  loading.value = false;
}
</script>

<template>
  <div class="flex w-full max-w-lg flex-col items-center gap-4 rounded-2xl p-6 shadow-xl">
    <h1 class="text-2xl font-bold">{{ t("forgot.title") }}</h1>

    <div class="flex w-full flex-col gap-2">
      <label class="flex w-full flex-col gap-1">
        {{ t("forgot.email") }}

        <input
          type="email"
          required
          v-model="email"
          :placeholder="t('forgot.email_placeholder')"
          class="hover:ring-claret-200 focus:ring-claret-600 w-full rounded-lg border border-zinc-300 px-4 py-2 duration-200 outline-none placeholder:text-black/50 hover:ring-2 focus:ring-2"
        />
      </label>
    </div>

    <p
      v-if="error"
      class="bg-claret-100 border-claret-500 text-claret-700 w-full rounded-xl border-2 px-4 py-2"
    >
      {{ t(error) }}
    </p>

    <p
      v-if="success"
      class="w-full rounded-xl border-2 border-emerald-500 bg-emerald-100 px-4 py-2 text-emerald-700"
    >
      {{ t(success) }}
    </p>

    <hr class="my-2 h-px w-full rounded-full border border-zinc-300" />

    <button
      @click="submit"
      :disabled="loading"
      class="bg-claret-600 disabled:bg-claret-700 border-claret-600 enabled:hover:text-claret-600 disabled:border-claret-700 flex w-full cursor-pointer items-center justify-center gap-2 rounded-xl border-2 p-2 py-3 font-semibold text-white transition-all duration-200 hover:shadow-md enabled:hover:bg-transparent disabled:cursor-progress disabled:opacity-50"
    >
      <Mail class="size-6" :class="{ 'animate-spin': loading }" />

      {{ loading ? t("forgot.loading") : t("forgot.action") }}
    </button>
  </div>
</template>
