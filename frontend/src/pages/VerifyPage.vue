<script setup lang="ts">
import ErrorDialog from "@/components/shared/ErrorDialog.vue";
import NavigationBar from "@/components/shared/NavigationBar.vue";
import OverlayScreen from "@/components/shared/OverlayScreen.vue";
import WhiteContainer from "@/components/shared/WhiteContainer.vue";
import { endpoints } from "@/consts";
import { useAuthFetch } from "@/hooks/use-auth-fetch";
import { useProfileStore } from "@/stores/profile";
import { LucideCheckCircle2 } from "lucide-vue-next";
import { onMounted, ref } from "vue";
import { useRouter } from "vue-router";

const profile = useProfileStore();
const router = useRouter();
const { authFetch } = useAuthFetch();

const code = ref<string>("");
const sent = ref(false);
const loading = ref(false);
const error = ref<string>();

onMounted(() => {
  if (!profile.hasProfile) {
    router.push({ path: "/login" });
  }
});

async function send() {
  loading.value = true;
  error.value = "";
  try {
    const res = await authFetch(endpoints.auth.verify, {
      method: "POST",
    });
    if (!res.ok) {
      error.value = "verify.error_cant_send";
    } else {
      sent.value = true;
    }
  } finally {
    loading.value = false;
  }
}

async function verify() {
  loading.value = true;
  error.value = "";

  try {
    const res = await authFetch(endpoints.auth.verifyCheck, {
      method: "POST",
      body: JSON.stringify({ code: code.value.toString() }),
    });

    if (!res.ok) {
      error.value = "verify.error_cant_verify";
    } else {
      await profile.fetchProfile();
    }
  } finally {
    loading.value = false;
  }
}
</script>

<template>
  <OverlayScreen shown v-if="error">
    <ErrorDialog :title="$t('verify.error')" :description="$t(error)" @close="error = undefined" />
  </OverlayScreen>

  <WhiteContainer class="justify-start gap-6">
    <NavigationBar />

    <div
      class="flex h-full w-full flex-col items-center justify-center gap-4"
      v-if="profile.profile && profile.profile.verified"
    >
      <LucideCheckCircle2 class="size-16 text-emerald-600" />

      <h1 class="text-2xl font-bold">{{ $t("verify.already_verified") }}</h1>
      <p>{{ $t("verify.already_verified_description") }}</p>
    </div>
    <div class="flex h-full w-full flex-col items-center justify-center gap-4" v-else>
      <!-- The dialog box to verify -->
      <div class="flex w-full max-w-xl flex-col gap-4 rounded-2xl bg-white p-6 shadow-md">
        <h1 class="text-center text-xl font-bold">{{ $t("verify.not_verified") }}</h1>

        <label class="flex flex-col gap-2">
          {{ $t("verify.code") }}

          <input
            type="number"
            class="focus:ring-claret-600 rounded-lg border border-zinc-300 px-4 py-2 text-center duration-200 outline-none hover:border-zinc-500 focus:ring-2"
            placeholder="123456"
            v-model.trim="code"
          />
        </label>

        <p class="text-claret-600 text-center font-semibold" v-if="sent">
          {{ $t("verify.email_sent") }}
        </p>

        <div class="flex w-full flex-col gap-2">
          <button
            class="bg-claret-600 border-claret-600 enabled:hover:text-claret-600 w-full cursor-pointer rounded-lg border-2 px-4 py-2 font-semibold text-white duration-200 enabled:hover:bg-white disabled:cursor-not-allowed disabled:opacity-50"
            :disabled="loading"
            @click="send"
          >
            {{ !loading ? $t("verify.send") : $t("verify.sending") }}
          </button>
          <button
            class="w-full cursor-pointer rounded-lg border-2 border-black bg-black px-4 py-2 font-semibold text-white duration-200 enabled:hover:bg-white enabled:hover:text-black disabled:cursor-not-allowed disabled:opacity-50"
            :disabled="loading"
            @click="verify"
            v-if="sent"
          >
            {{ !loading ? $t("verify.verify") : $t("verify.verifying") }}
          </button>
        </div>
      </div>
    </div>
  </WhiteContainer>
</template>

<style lang="css" scoped>
input[type="number"] {
  appearance: textfield;
}

input[type="number"]::-webkit-inner-spin-button,
input[type="number"]::-webkit-outer-spin-button {
  display: none;
}
</style>
