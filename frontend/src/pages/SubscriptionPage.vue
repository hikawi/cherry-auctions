<script setup lang="ts">
import NavigationBar from "@/components/shared/NavigationBar.vue";
import WhiteContainer from "@/components/shared/WhiteContainer.vue";
import { endpoints } from "@/consts";
import { useAuthFetch } from "@/hooks/use-auth-fetch";
import { useProfileStore } from "@/stores/profile";
import dayjs from "dayjs";
import { LucideUnlock } from "lucide-vue-next";
import { computed, ref } from "vue";
import { useI18n } from "vue-i18n";

const profile = useProfileStore();
const { authFetch } = useAuthFetch();
const { locale } = useI18n();

const loading = ref(false);
const expiresIn = computed(() => {
  const time = profile.profile?.subscription?.expires_at;
  if (time) {
    return dayjs(time).locale(locale.value).fromNow();
  }
  return undefined;
});

async function requestPrivileges() {
  loading.value = true;

  try {
    const res = await authFetch(endpoints.users.request, {
      method: "POST",
    });
    if (res.status == 204) {
      profile.fetchProfile();
    } else {
      console.log(res.status);
    }
  } finally {
    loading.value = false;
  }
}
</script>

<template>
  <WhiteContainer class="justify-start gap-8">
    <NavigationBar />

    <section class="flex w-full max-w-4xl flex-col gap-4">
      <h2 class="text-2xl font-bold">{{ $t("subscriptions.title") }}</h2>

      <div
        class="border-watermelon-500 bg-watermelon-100 text-watermelon-500 flex w-full flex-col rounded-xl border-2 p-4"
        v-if="!profile.profile?.subscription"
      >
        {{ $t("subscriptions.no_subscriptions") }}
      </div>
      <div
        class="flex w-full flex-col rounded-xl border-2 border-emerald-500 bg-emerald-100 p-4 text-emerald-500"
        v-else
      >
        {{ $t("subscriptions.expires_in", { in: expiresIn }) }}
      </div>

      <button
        class="bg-claret-600 hover:bg-claret-700 flex w-fit cursor-pointer flex-row items-center justify-center gap-2 self-end rounded-full px-4 py-1 font-semibold text-white duration-200 disabled:cursor-not-allowed disabled:opacity-50"
        :disabled="profile.profile?.waiting_approval"
        @click="requestPrivileges"
      >
        <LucideUnlock class="size-4" />
        {{
          loading
            ? $t("subscriptions.requesting")
            : profile.profile?.waiting_approval
              ? $t("subscriptions.already_requested")
              : $t("subscriptions.request")
        }}
      </button>
    </section>
  </WhiteContainer>
</template>
