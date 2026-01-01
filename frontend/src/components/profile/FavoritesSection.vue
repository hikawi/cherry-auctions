<script setup lang="ts">
import { endpoints } from "@/consts";
import { useAuthFetch } from "@/hooks/use-auth-fetch";
import { onMounted, ref, watch } from "vue";
import ProductCard from "../index/ProductCard.vue";

const { authFetch } = useAuthFetch();

const data = ref();
const loading = ref(true);
const page = ref(1);
const maxPages = ref(1);

watch(maxPages, (val) => (page.value = Math.min(val, page.value)));
watch(page, fetchFavorites);

async function fetchFavorites() {
  const url = new URL(endpoints.products.favorite);
  url.searchParams.append("page", page.value.toString());
  url.searchParams.append("per_page", "12");

  loading.value = true;
  try {
    const res = await authFetch(url, {});
    if (res.ok) {
      const json = await res.json();
      data.value = json.data;
      maxPages.value = json.total_pages;
    }
  } finally {
    loading.value = false;
  }
}

onMounted(() => {
  fetchFavorites();
});
</script>

<template>
  <h2 class="text-2xl font-semibold">{{ $t("profile.favorites") }}</h2>

  <div
    class="grid w-full grid-cols-1 gap-4 sm:grid-cols-2 md:grid-cols-3"
    v-if="data && data.length > 0"
  >
    <template v-for="product in data" :key="product.id">
      <ProductCard :product @favoriteToggle="fetchFavorites" />
    </template>
  </div>
  <p class="w-full py-6 text-center text-xl font-semibold" v-else>
    {{ $t("profile.no_favorites") }}
  </p>
</template>
