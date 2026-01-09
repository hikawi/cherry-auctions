<script setup lang="ts">
import { endpoints } from "@/consts";
import { useAuthFetch } from "@/hooks/use-auth-fetch";
import { ref } from "vue";
import TextInput from "../shared/inputs/TextInput.vue";

const { authFetch } = useAuthFetch();

const currentPassword = ref("");
const newPassword = ref("");
const confirmNewPassword = ref("");
const loading = ref(false);
const error = ref("");
const success = ref(false);

async function changePassword() {
  loading.value = true;
  error.value = "";

  if (confirmNewPassword.value != newPassword.value) {
    error.value = "profile.error_passwords_dont_match";
    loading.value = false;
    return;
  }

  try {
    const res = await authFetch(endpoints.users.me.password, {
      method: "PUT",
      body: JSON.stringify({
        current_password: currentPassword.value,
        new_password: newPassword.value,
      }),
    });

    switch (res.status) {
      case 421:
        error.value = "profile.error_oauth";
        break;
      case 400:
        error.value = "profile.error_invalid_password";
        break;
      case 403:
        error.value = "profile.error_wrong_password";
        break;
      case 200:
        success.value = true;
        break;
    }
  } finally {
    loading.value = false;
  }
}
</script>

<template>
  <h2 class="text-2xl font-semibold">{{ $t("profile.change_password") }}</h2>

  <form class="flex w-full flex-col items-center gap-4" novalidate>
    <TextInput :label="$t('profile.current_password')" type="password" v-model="currentPassword" />
    <TextInput :label="$t('profile.new_password')" type="password" v-model="newPassword" />
    <TextInput
      :label="$t('profile.confirm_password')"
      type="password"
      v-model="confirmNewPassword"
    />

    <div
      v-if="success"
      class="w-full rounded-xl border-2 border-emerald-600 bg-emerald-200/50 px-4 py-2 text-emerald-600"
    >
      {{ $t("profile.changed_password") }}
    </div>

    <div
      v-if="error"
      class="border-watermelon-600 bg-watermelon-200/50 text-watermelon-600 w-full rounded-xl border-2 px-4 py-2"
    >
      {{ $t(error) }}
    </div>

    <button
      type="submit"
      @click.prevent="changePassword"
      class="bg-claret-600 hover:bg-claret-700 mt-2 flex w-fit cursor-pointer flex-row items-center justify-center gap-2 self-end rounded-full px-4 py-1 font-semibold text-white duration-200"
    >
      {{ loading ? $t("general.loading") : $t("profile.change_password") }}
    </button>
  </form>
</template>
