<script setup lang="ts">
import type { ProductListing } from "@/types";
import { computed } from "vue";

const props = defineProps<{
  product: ProductListing;
}>();

const productLink = computed(() => `/products/${props.product.id}`);
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

    <div class="flex w-full flex-col gap-0">
      <div class="flex flex-row items-center justify-between">
        <span>{{ $t("products.current_bid") }}</span>
        <span class="text-claret-600 text-xl font-semibold"
          >${{ product.starting_bid.toLocaleString() }}</span
        >
      </div>
      <div class="flex flex-row items-center justify-between">
        <span>{{ $t("products.bin_price") }}</span>
        <span class="text-claret-600 text-xl font-semibold"
          >${{ product.bin_price.toLocaleString() }}</span
        >
      </div>
    </div>
  </a>
</template>
