<script setup lang="ts">
import type { Product } from "@/types";
import { computed } from "vue";
import { useI18n } from "vue-i18n";
import dayjs from "dayjs";
import { useTimestamp } from "@vueuse/core";
import PlaceholderAvatar from "../shared/PlaceholderAvatar.vue";

const props = defineProps<{
  product: Product;
}>();
const { locale } = useI18n();
const now = useTimestamp({ interval: 1000 });

const productLink = computed(() => `/products/${props.product.id}`);
const expiresDisplay = computed(() => {
  return dayjs(props.product.expired_at).locale(locale.value).from(now.value);
});
const expiresAtDisplay = computed(() => {
  return dayjs(props.product.expired_at).locale(locale.value).format("ll");
});
const shouldBeRelative = computed(() => {
  const expiration = dayjs(props.product.expired_at);
  const diffMs = expiration.diff(now.value);
  const THREE_DAYS_MS = 3 * 24 * 60 * 60 * 1000;
  return diffMs > 0 && diffMs <= THREE_DAYS_MS;
});
</script>

<template>
  <a
    :href="productLink"
    class="flex flex-col gap-4 rounded-lg border border-zinc-300 bg-white p-4 duration-200 hover:border-zinc-500"
  >
    <span class="text-lg font-semibold">{{ product.name }}</span>
    <img
      :src="product.thumbnail_url"
      :alt="product.name"
      class="aspect-video h-auto w-full rounded-lg object-cover object-center"
    />

    <div class="flex w-full flex-row items-center justify-start gap-2">
      <PlaceholderAvatar :name="product.seller.name" />
      <span>{{ product.seller.name }}</span>
    </div>

    <div class="flex w-full flex-col gap-0">
      <div class="flex flex-row items-center justify-between" v-if="product.current_highest_bid">
        <span>{{ $t("products.current_bid") }}</span>
        <span class="text-claret-600 text-xl font-semibold"
          >${{ product.current_highest_bid.price.toLocaleString() }}</span
        >
      </div>
      <div class="flex flex-row items-center justify-between" v-else>
        <span>{{ $t("products.starting_bid") }}</span>
        <span class="text-claret-600 text-xl font-semibold"
          >${{ product.starting_bid.toLocaleString() }}</span
        >
      </div>
      <div class="flex flex-row items-center justify-between" v-if="product.bin_price">
        <span>{{ $t("products.bin_price") }}</span>
        <span class="text-claret-600 text-xl font-semibold"
          >${{ product.bin_price.toLocaleString() }}</span
        >
      </div>
    </div>

    <div class="flex w-full flex-row items-center justify-between text-sm text-black/75">
      <span v-if="product.bids_count != 1">{{
        $t("products.bids_count_plural", { count: product.bids_count })
      }}</span>
      <span v-else>{{ $t("products.bids_count_singular", { count: product.bids_count }) }}</span>

      <span v-if="shouldBeRelative">{{ $t("products.expires_in", { in: expiresDisplay }) }}</span>
      <span v-else>{{ $t("products.expires_at", { at: expiresAtDisplay }) }}</span>
    </div>
  </a>
</template>
