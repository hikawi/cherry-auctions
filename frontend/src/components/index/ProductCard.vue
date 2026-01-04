<script setup lang="ts">
import type { Product } from "@/types";
import { computed } from "vue";
import { useI18n } from "vue-i18n";
import dayjs from "dayjs";
import { useTimestamp } from "@vueuse/core";
import { LucideHeart, LucideSparkle } from "lucide-vue-next";
import { useAuthFetch } from "@/hooks/use-auth-fetch";
import { endpoints } from "@/consts";
import { useProfileStore } from "@/stores/profile";
import AvatarCircle from "../shared/AvatarCircle.vue";

const props = defineProps<{
  product: Product;
}>();
const emits = defineEmits<{
  favoriteToggle: [id: number];
}>();

const { locale } = useI18n();
const now = useTimestamp({ interval: 1000 });
const { authFetch } = useAuthFetch();
const profile = useProfileStore();

const productLink = computed(() => `/products/${props.product.id}`);
const expiresDisplay = computed(() => {
  return dayjs(props.product.expired_at).locale(locale.value).from(now.value);
});
const expiresAtDisplay = computed(() => {
  return dayjs(props.product.expired_at).locale(locale.value).format("ll");
});
const createdAtDisplay = computed(() => {
  return dayjs(props.product.created_at).locale(locale.value).format("lll");
});
const shouldBeRelative = computed(() => {
  const expiration = dayjs(props.product.expired_at);
  const diffMs = expiration.diff(now.value);
  const THREE_DAYS_MS = 3 * 24 * 60 * 60 * 1000;
  return diffMs > 0 && diffMs <= THREE_DAYS_MS;
});
const isNew = computed(() => {
  const diffMs = dayjs(props.product.created_at).diff(now.value);
  const ONE_HOUR_MS = 60 * 60 * 1000;
  return Math.abs(diffMs) <= ONE_HOUR_MS;
});

async function toggleFavorite() {
  const res = await authFetch(endpoints.products.favorite, {
    method: "POST",
    body: JSON.stringify({ id: props.product.id }),
  });
  if (res.ok) {
    emits("favoriteToggle", props.product.id);
  }
}
</script>

<template>
  <div
    class="relative flex flex-col gap-4 rounded-lg border border-zinc-300 bg-white p-4 pt-8 duration-200 hover:border-zinc-500"
  >
    <button
      class="absolute top-0 right-0 z-10 translate-x-1/2 -translate-y-1/2 cursor-pointer rounded-full border border-zinc-300 bg-white p-2 shadow-md hover:border-zinc-500"
      @click="toggleFavorite"
      v-if="profile.hasProfile"
    >
      <LucideHeart
        class="duration-200"
        :class="{
          'text-claret-600 hover:fill-claret-600/50': !product.is_favorite,
          'fill-claret-600 text-claret-600': product.is_favorite,
        }"
      />
    </button>

    <span
      class="from-watermelon-500 bg-linear-to-r via-violet-500 to-pink-500 bg-clip-text text-lg font-black text-transparent uppercase"
      v-if="isNew"
    >
      <LucideSparkle
        class="fill-watermelon-500 text-watermelon-500 inline-block size-4 -translate-y-px"
      />
      {{ $t("products.new") }}
    </span>

    <a :href="productLink" class="hover:text-claret-600 text-lg font-semibold duration-200">{{
      product.name
    }}</a>
    <img
      :src="product.thumbnail_url"
      :alt="product.name"
      class="aspect-video h-auto w-full rounded-lg object-cover object-center"
    />

    <div class="flex w-full flex-row items-center justify-start gap-2">
      <AvatarCircle :name="product.seller.name" :avatar_url="product.seller.avatar_url" />
      <span>{{ product.seller.name }}</span>
    </div>

    <span class="text-sm">{{ $t("products.created_at", { at: createdAtDisplay }) }}</span>

    <div class="flex w-full flex-col gap-0">
      <div class="flex flex-row items-center justify-between" v-if="product.current_highest_bid">
        <span>{{ $t("products.current_bid") }}</span>
        <span class="text-claret-600 text-xl font-semibold">{{
          $n(product.current_highest_bid.price / 100, "currency")
        }}</span>
      </div>
      <div class="flex flex-row items-center justify-between" v-else>
        <span>{{ $t("products.starting_bid") }}</span>
        <span class="text-claret-600 text-xl font-semibold">{{
          $n(product.starting_bid / 100, "currency")
        }}</span>
      </div>
      <div class="flex flex-row items-center justify-between" v-if="product.bin_price">
        <span>{{ $t("products.bin_price") }}</span>
        <span class="text-claret-600 text-xl font-semibold">{{
          $n(product.bin_price / 100, "currency")
        }}</span>
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
  </div>
</template>
