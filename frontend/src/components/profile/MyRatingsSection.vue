<script setup lang="ts">
import { endpoints } from "@/consts";
import { useAuthFetch } from "@/hooks/use-auth-fetch";
import type { Rating } from "@/types";
import dayjs from "dayjs";
import { onMounted, ref, watch } from "vue";
import { useI18n } from "vue-i18n";
import AvatarCircle from "../shared/AvatarCircle.vue";
import PagingSection from "../shared/PagingSection.vue";
import { LucideThumbsUp, LucideThumbsDown } from "lucide-vue-next";
import { useProfileStore } from "@/stores/profile";

const { locale } = useI18n();
const { authFetch } = useAuthFetch();
const profile = useProfileStore();

const page = ref(1);
const maxPages = ref(1);
const error = ref<string>("");
const ratings = ref<Rating[]>();

watch(page, fetchRatings);
onMounted(fetchRatings);

async function fetchRatings() {
  const url = new URL(endpoints.users.me.rated);
  url.searchParams.append("page", page.value.toString());
  url.searchParams.append("per_page", "8");

  try {
    const res = await authFetch(url);
    if (res.ok) {
      const json = await res.json();
      ratings.value = json.data;
      maxPages.value = Math.max(json.total_pages, 1);
    }
  } catch {
    error.value = "profile.error_cant_fetch_ratings";
  }
}

function renderAtTime(time: string): string {
  return dayjs(time).locale(locale.value).format("lll");
}
</script>

<template>
  <h2 class="text-2xl font-semibold">{{ $t("profile.others_ratings") }}</h2>

  <p>
    <span class="font-semibold">{{ $t("profile.rating") }}</span>
    {{ profile.profile?.average_rating }}
  </p>

  <div
    v-if="error"
    class="border-watermelon-600 bg-watermelon-200/50 text-watermelon-600 w-full rounded-xl border-2 px-4 py-2"
  >
    {{ $t(error) }}
  </div>
  <div class="flex flex-col gap-4" v-else-if="ratings && ratings.length > 0">
    <template v-for="(rating, idx) in ratings" :key="idx">
      <div
        class="flex w-full flex-col gap-2 rounded-xl bg-white p-4 shadow-md sm:p-6"
        :class="{
          'border-2 border-emerald-600': rating.rating == 1,
          'border-watermelon-600 border-2': rating.rating == 0,
        }"
      >
        <div class="flex w-full flex-row items-center justify-start gap-2">
          <AvatarCircle
            :name="rating.reviewer.name"
            :avatar_url="rating.reviewer.avatar_url"
            class="size-5!"
          />
          <span class="font-semibold">{{
            rating.reviewer.name || $t("general.deleted_user")
          }}</span>

          <LucideThumbsDown v-if="rating.rating == 0" class="fill-watermelon-600 size-5" />
          <LucideThumbsUp v-else class="size-5 fill-emerald-600" />
        </div>

        <p>
          {{ rating.feedback }}
        </p>

        <p class="text-sm text-black/50">
          {{ renderAtTime(rating.created_at) }}
        </p>
      </div>
    </template>

    <!-- Paging section -->
    <PagingSection :page :maxPages @setPage="(p) => (page = p)" />
  </div>
  <p class="w-full py-6 text-center text-xl font-semibold" v-else>
    {{ $t("profile.no_ratings") }}
  </p>
</template>
