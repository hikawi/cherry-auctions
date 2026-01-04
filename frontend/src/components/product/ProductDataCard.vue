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

const cantBidReason = computed(() => {
  if (!props.data) {
    return undefined;
  }

  if (!profile.profile) {
    return "products.cant_bid_logged_out";
  }

  if (props.data.seller.id == profile.profile.id) {
    return "products.cant_bid_self";
  }

  if (!props.data.allows_unrated_buyers && profile.profile.average_rating < 0.8) {
    return "products.cant_bid_no_rating";
  }

  // TODO:
  // Add the flow to check if the user is denied.

  return undefined;
});
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
const nextBidValue = computed(() => {
  if (!props.data) {
    return 0;
  }

  if (!props.data.current_highest_bid) {
    return props.data.starting_bid;
  }

  const highestBid = props.data.current_highest_bid.price;
  return highestBid + props.data.step_bid_value;
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
        v-if="cantBidReason"
      >
        {{ $t(cantBidReason) }}
      </span>

      <button
        class="bg-claret-600 enabled:hover:text-claret-600 border-claret-600 flex cursor-pointer flex-col items-center justify-center gap-0 rounded-xl border-2 px-4 py-2 text-white duration-200 enabled:hover:bg-white disabled:cursor-not-allowed disabled:opacity-50"
        :disabled="cantBidReason != undefined"
      >
        <span class="text-lg font-semibold">{{ $t("products.bid") }}</span>
        <span class="text-sm">{{ $n(nextBidValue / 100, "currency") }}</span>
      </button>

      <button
        class="flex cursor-pointer flex-col items-center justify-center gap-0 rounded-xl border-2 border-black bg-black px-4 py-2 text-white duration-200 enabled:hover:bg-white enabled:hover:text-black disabled:cursor-not-allowed disabled:opacity-50"
        :disabled="cantBidReason != undefined"
        v-if="data.bin_price"
      >
        <span class="text-lg font-semibold">{{ $t("products.bin") }}</span>
        <span class="text-sm">{{ $n(data.bin_price / 100, "currency") }}</span>
      </button>
    </div>
  </div>
</template>
