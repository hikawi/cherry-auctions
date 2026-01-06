<script setup lang="ts">
import { endpoints } from "@/consts";
import { useAuthFetch } from "@/hooks/use-auth-fetch";
import { onMounted, ref, watch } from "vue";
import ProductCard from "../index/ProductCard.vue";
import { LucideChevronLeft, LucideChevronRight } from "lucide-vue-next";

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
  url.searchParams.append("sort", "time"); // Just default to expired by, not required in this form.
  url.searchParams.append("asc", "true");

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

    <!-- Paging section -->
    <div class="flex w-full flex-row items-center justify-between sm:col-span-2 md:col-span-3">
      <button
        class="cursor-pointer rounded-lg border border-zinc-300 p-2 hover:border-zinc-500 disabled:cursor-not-allowed disabled:opacity-50"
        @click="page = Math.max(page - 1, 1)"
        :disabled="page == 1"
      >
        <LucideChevronLeft class="size-4 text-black" />
      </button>

      <span>{{ page }} / {{ maxPages }}</span>

      <button
        class="cursor-pointer rounded-lg border border-zinc-300 p-2 hover:border-zinc-500 disabled:cursor-not-allowed disabled:opacity-50"
        @click="page = Math.min(page + 1, maxPages)"
        :disabled="page == maxPages"
      >
        <LucideChevronRight class="size-4 text-black" />
      </button>
    </div>
  </div>
  <p class="w-full py-6 text-center text-xl font-semibold" v-else>
    {{ $t("profile.no_favorites") }}
  </p>
</template>
