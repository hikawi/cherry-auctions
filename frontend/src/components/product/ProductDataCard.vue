<script setup lang="ts">
import { useProfileStore } from "@/stores/profile";
import type { Product } from "@/types";
import dayjs from "dayjs";
import { computed } from "vue";
import { useI18n } from "vue-i18n";
import AvatarCircle from "../shared/AvatarCircle.vue";

const profile = useProfileStore();
const { locale } = useI18n();

const props = defineProps<{
  data: Product & { similar_products?: Product[]; categories: { id: number; name: string }[] };
}>();

const expiresAtDisplay = computed(() => {
  return props.data != null
    ? dayjs(props.data.expired_at).locale(locale.value).format("lll")
    : "N/A";
});
const createdAtDisplay = computed(() => {
  return props.data != null
    ? dayjs(props.data.created_at).locale(locale.value).format("lll")
    : "N/A";
});
</script>

<template>
  <div class="flex h-full w-full flex-col justify-between gap-4">
    <div class="flex w-full flex-col gap-4">
      <h1 class="text-2xl font-bold">{{ data.name }}</h1>

      <div
        v-if="data.categories && data.categories.length > 0"
        class="flex w-full flex-row flex-wrap items-center gap-2 text-sm"
      >
        <div
          v-for="category in data.categories"
          :key="category.id"
          class="bg-claret-100 rounded-lg px-2 py-1"
        >
          {{ category.name }}
        </div>
      </div>

      <div class="flex w-full flex-row items-center gap-2">
        <AvatarCircle
          :name="data.seller.name"
          :avatar_url="data.seller.avatar_url"
          class="size-8"
        />
        <span class="text-lg"
          >{{ data.seller.name }} ({{
            data.seller.email ? data.seller.email : $t("products.deleted_email")
          }})</span
        >
      </div>

      <div class="flex w-full flex-col gap-1">
        <span class="text-sm">{{ $t("products.created_at", { at: createdAtDisplay }) }}</span>
        <span class="text-sm">{{ $t("products.expires_at", { at: expiresAtDisplay }) }}</span>
      </div>
    </div>

    <div class="grid grid-cols-1 gap-2 md:grid-cols-2">
      <span
        class="text-claret-600 w-full text-center font-semibold md:col-span-2"
        v-if="
          !profile.profile || profile.profile.average_rating < 0.8 || !data.allows_unrated_buyers
        "
      >
        {{ $t("products.cant_bid") }}
      </span>

      <button
        class="bg-claret-600 enabled:hover:text-claret-600 border-claret-600 flex cursor-pointer flex-col items-center justify-center gap-0 rounded-xl border-2 px-4 py-2 text-white duration-200 enabled:hover:bg-white disabled:cursor-not-allowed disabled:opacity-50"
        :disabled="
          !profile.profile || profile.profile.average_rating < 0.8 || !data.allows_unrated_buyers
        "
      >
        <span class="text-lg font-semibold">{{ $t("products.bid") }}</span>
        <span class="text-sm">${{ data.current_highest_bid?.price || data.starting_bid }}</span>
      </button>

      <button
        class="flex cursor-pointer flex-col items-center justify-center gap-0 rounded-xl border-2 border-black bg-black px-4 py-2 text-white duration-200 hover:bg-white hover:text-black"
        v-if="data.bin_price"
      >
        <span class="text-lg font-semibold">{{ $t("products.bin") }}</span>
        <span class="text-sm">${{ data.bin_price }}</span>
      </button>
    </div>
  </div>
</template>
