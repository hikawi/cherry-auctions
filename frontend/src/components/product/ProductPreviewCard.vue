<script setup lang="ts">
import type { Product } from "@/types";
import dayjs from "dayjs";
import { computed } from "vue";
import AvatarCircle from "../shared/AvatarCircle.vue";
import { useI18n } from "vue-i18n";
import { useProfileStore } from "@/stores/profile";
import { LucideHeart } from "lucide-vue-next";
import { useAuthFetch } from "@/hooks/use-auth-fetch";
import { endpoints } from "@/consts";

const { locale } = useI18n();
const profile = useProfileStore();
const { authFetch } = useAuthFetch();

type CardFeature = "datetime" | "seller" | "price" | "favorite";

const props = defineProps<{
  product: Product;
  enabledFeatures?: CardFeature[];
}>();

const emits = defineEmits<{
  favoriteToggle: [id: number];
}>();

const features = computed(() => ({
  datetime: props.enabledFeatures?.includes("datetime"),
  seller: props.enabledFeatures?.includes("seller"),
  price: props.enabledFeatures?.includes("price"),
  favorite: props.enabledFeatures?.includes("favorite"),
}));
const isNew = computed(() => {
  const diff = dayjs(props.product.created_at).diff();
  return Math.abs(diff) < 15 * 60 * 1000; // 15 minutes
});
const isExpiring = computed(() => {
  const diff = dayjs(props.product.expired_at).diff();
  return diff >= 0 && diff < 3 * 24 * 60 * 60 * 1000; // 3 days
});
const hasWon = computed(
  () =>
    props.product.current_highest_bid &&
    props.product.current_highest_bid.bidder.id == profile.profile?.id,
);
const hasBid = computed(
  () =>
    props.product.product_state == "active" &&
    props.product.bids?.some((bid) => bid.bidder.id == profile.profile?.id) == true,
);
const createdAt = computed(() =>
  dayjs(props.product.created_at).locale(locale.value).format("lll"),
);
const expiresAt = computed(() =>
  dayjs(props.product.expired_at).locale(locale.value).format("lll"),
);
const expiresIn = computed(() => dayjs(props.product.expired_at).locale(locale.value).fromNow());

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
    class="relative flex h-full w-full flex-col items-center justify-between gap-4 rounded-xl shadow-xl"
    :class="{
      'ring-2 ring-amber-600': isExpiring && !hasBid && !hasWon && !isNew,
      'ring-2 ring-yellow-500': hasBid && !hasWon,
      'ring-2 ring-emerald-600': hasWon,
      'ring-claret-600 ring-2': isNew && !hasWon && !hasBid && !isExpiring,
    }"
  >
    <div
      class="from-claret-600 absolute top-0 left-1/2 w-fit -translate-x-1/2 rounded-b-xl bg-linear-to-r to-pink-600 px-2 pt-px pb-1.5 text-sm font-bold text-white"
      v-if="isNew"
    >
      {{ $t("products.new") }}
    </div>
    <div
      class="absolute top-0 left-1/2 w-fit -translate-x-1/2 rounded-b-xl bg-linear-to-r from-emerald-600 to-green-600 px-2 pt-px pb-1.5 text-sm font-bold text-white"
      v-else-if="hasWon"
    >
      {{ $t("products.won") }}
    </div>
    <div
      class="absolute top-0 left-1/2 w-fit -translate-x-1/2 rounded-b-xl bg-linear-to-r from-yellow-600 to-amber-600 px-2 pt-px pb-1.5 text-sm font-bold text-white"
      v-else-if="hasBid"
    >
      {{ $t("products.bidding") }}
    </div>
    <div
      class="absolute top-0 left-1/2 w-fit -translate-x-1/2 rounded-b-xl bg-linear-to-r from-amber-600 to-orange-600 px-2 pt-px pb-1.5 text-sm font-bold text-white"
      v-else-if="isExpiring"
    >
      {{ $t("products.expiring_soon") }}
    </div>

    <!-- Favorite button -->
    <button
      class="absolute top-4 right-4 z-10 cursor-pointer rounded-full border border-zinc-300 bg-white p-2 shadow-md hover:border-zinc-500"
      @click="toggleFavorite"
      v-if="profile.hasProfile && features.favorite"
    >
      <LucideHeart
        class="size-4 duration-200"
        :class="{
          'text-claret-600 hover:fill-claret-600/50': !product.is_favorite,
          'fill-claret-600 text-claret-600': product.is_favorite,
        }"
      />
    </button>

    <div class="flex w-full flex-col gap-4">
      <img
        :src="product.thumbnail_url"
        :alt="product.name"
        class="aspect-video h-auto w-full rounded-t-xl object-cover object-center"
      />

      <div class="flex w-full flex-col px-4 text-left">
        <router-link
          :to="{ name: 'product-details', params: { id: product.id } }"
          class="hover:text-claret-600 text-lg font-semibold duration-200"
          >{{ product.name }}</router-link
        >
      </div>

      <!-- Date time section -->
      <div class="flex w-full flex-col gap-0 px-4 text-sm" v-if="features.datetime">
        <span class="text-sm">{{ $t("products.created_at", { at: createdAt }) }}</span>
        <span class="text-sm">{{
          isExpiring
            ? $t("products.expires_in", { in: expiresIn })
            : $t("products.expires_at", { at: expiresAt })
        }}</span>
      </div>

      <!-- Seller section -->
      <div
        class="flex w-full flex-row items-center justify-start gap-2 px-4 pb-4 text-sm"
        v-if="features.seller"
      >
        <AvatarCircle
          :name="product.seller.name || $t('general.deleted_user')"
          :avatar_url="product.seller.avatar_url"
          class="size-6! border border-zinc-300"
        />

        <span class="font-semibold">
          {{ product.seller.name || $t("general.deleted_user") }}
        </span>
      </div>
    </div>

    <!-- Pricing section -->
    <div class="flex w-full flex-col gap-0 px-4 pb-4 text-sm" v-if="features.price">
      <div class="flex flex-row items-center justify-between" v-if="product.current_highest_bid">
        <span>{{ $t("products.current_bid") }}</span>
        <span class="text-claret-600 text-xl font-bold">{{
          $n(product.current_highest_bid.price / 100, "currency")
        }}</span>
      </div>
      <div class="flex flex-row items-center justify-between" v-else>
        <span>{{ $t("products.starting_bid") }}</span>
        <span class="text-claret-600 text-xl font-bold">{{
          $n(product.starting_bid / 100, "currency")
        }}</span>
      </div>
      <div class="flex flex-row items-center justify-between" v-if="product.bin_price">
        <span>{{ $t("products.bin_price") }}</span>
        <span class="text-claret-600 text-xl font-bold">{{
          $n(product.bin_price / 100, "currency")
        }}</span>
      </div>
    </div>
  </div>
</template>
