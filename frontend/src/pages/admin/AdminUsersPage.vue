<script setup lang="ts">
import AvatarCircle from "@/components/shared/AvatarCircle.vue";
import LoadingSpinner from "@/components/shared/LoadingSpinner.vue";
import { endpoints } from "@/consts";
import { useAuthFetch } from "@/hooks/use-auth-fetch";
import dayjs from "dayjs";
import { onMounted, ref } from "vue";
import { useI18n } from "vue-i18n";

const { authFetch } = useAuthFetch({ json: true });
const { locale } = useI18n();

const data = ref();
const loading = ref(true);
const page = ref(1);
const maxPages = ref(1);
const perPage = 20;

function buildUsersURL(): URL {
  const url = new URL(endpoints.users.all);
  url.searchParams.append("page", page.value.toString());
  url.searchParams.append("per_page", perPage.toString());
  return url;
}

async function loadUsers() {
  loading.value = true;
  try {
    const res = await authFetch(buildUsersURL());
    if (res.ok) {
      const json = await res.json();
      maxPages.value = json.total_pages;
      page.value = json.page;
      data.value = json.data;
    }
  } finally {
    loading.value = false;
  }
}

async function approve(id: number) {
  const res = await authFetch(endpoints.users.approve, {
    method: "POST",
    body: JSON.stringify({ id }),
  });
  if (res.status == 204) {
    loadUsers();
  }
}

function parseSusbcriptionExpires(subscription: Record<string, string>) {
  if (subscription.expired_at) {
    return dayjs(subscription.expired_at).locale(locale.value).fromNow();
  }
  return "n/a";
}

onMounted(loadUsers);
</script>

<template>
  <h1 class="text-2xl font-bold">{{ $t("admin.users.title") }}</h1>

  <div class="w-full py-4" v-if="loading">
    <LoadingSpinner />
  </div>
  <div class="w-full py-4 text-xl font-semibold" v-else-if="!data">
    <p>{{ $t("admin.users.cant_load") }}</p>
  </div>
  <div class="flex w-full max-w-4xl flex-col gap-8" v-else>
    <div class="flex w-full flex-row items-center justify-between">
      <span>{{ $t("admin.users.page", { page: page, max_pages: maxPages }) }}</span>
    </div>

    <div class="flex w-full flex-col gap-4">
      <template v-for="user in data" :key="user.id">
        <div
          class="flex flex-col gap-2 rounded-xl border border-zinc-300 p-4 duration-200 hover:border-zinc-500"
        >
          <div class="flex flex-row items-center gap-2 text-lg font-semibold">
            <AvatarCircle :name="user.name" :avatar_url="user.avatar_url" />
            {{ user.name }}
          </div>

          <ul class="list-inside list-disc">
            <li>
              {{ $t("admin.users.email", { email: user.email }) }}
            </li>
            <li>
              {{ $t("admin.users.verified", { verified: user.verified }) }}
            </li>
            <li>
              {{ $t("admin.users.average_rating", { rating: user.average_rating }) }}
            </li>
            <li>
              {{ $t("admin.users.roles", { roles: user.roles }) }}
            </li>
            <li v-if="user.subscription">
              {{
                $t("admin.users.subscription_expires_in", {
                  in: parseSusbcriptionExpires(user.subscription),
                })
              }}
            </li>
          </ul>

          <button
            v-if="user.waiting_approval"
            @click="() => approve(user.id)"
            class="bg-claret-600 hover:bg-claret-700 flex cursor-pointer flex-row items-center-safe justify-center gap-1 self-start rounded-full px-2 py-2 font-semibold text-white"
          >
            {{ $t("admin.users.approve") }}
          </button>
        </div>
      </template>
    </div>
  </div>
</template>
