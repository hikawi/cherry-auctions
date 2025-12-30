<script setup lang="ts">
import WhiteContainer from "@/components/shared/WhiteContainer.vue";
import { useProfileStore } from "@/stores/profile";
import { computed, onMounted, ref, watch, watchEffect } from "vue";
import {
  LucideChevronLeft,
  LucideChevronRight,
  LucideEllipsis,
  LucideSearch,
} from "lucide-vue-next";
import { useAuthFetch } from "@/hooks/use-auth-fetch";
import { endpoints } from "@/consts";
import type { ProductListing } from "@/types";
import ProductCard from "@/components/index/ProductCard.vue";

const profile = useProfileStore();
const { authFetch } = useAuthFetch();

const loading = ref(true);
const search = ref<string>();
const products = ref<ProductListing[]>();
const page = ref(1);
const total = ref(0);
const maxPages = ref(1);
const urlEncodedName = computed(() => {
  return `https://ui-avatars.com/api/?name=${encodeURIComponent(profile.profile?.name || "")}`;
});

const pageRange = computed(() => {
  return {
    min: Math.max(page.value - 2, 1),
    max: Math.min(page.value + 2, maxPages.value),
  };
});

watchEffect(() => {
  page.value = Math.min(page.value, maxPages.value);
});

watch(page, (nc, cc) => {
  if (nc != cc) {
    fetchProducts();
  }
});

async function fetchProducts() {
  loading.value = true;

  // First, we build the URL
  const url = new URL(endpoints.products.get);
  url.searchParams.append("page", Math.max(page.value, 1).toString());
  url.searchParams.append("per_page", "18"); // just an arbitrary number because it fits well.

  if (search.value) {
    url.searchParams.append("query", search.value);
  }

  try {
    const res = await authFetch(url);
    if (res.ok) {
      const json = await res.json();
      products.value = json.data;
      maxPages.value = json.total_pages;
      total.value = json.total;
    }
  } finally {
    loading.value = false;
  }

  // Scroll up
  window.scroll({ top: 0, behavior: "smooth" });
}

onMounted(async () => {
  await fetchProducts();
});
</script>

<template>
  <WhiteContainer class="justify-start">
    <!-- Classic 12-column layout -->
    <div class="grid h-fit w-full grid-cols-1 gap-4 sm:flex-row sm:gap-8 md:grid-cols-4 lg:px-6">
      <a href="/" class="hover:text-claret-600 text-center text-3xl font-bold"> CherryAuctions </a>

      <!-- NavigationBar -->
      <div class="flex flex-row items-center gap-4 md:col-span-3">
        <label
          class="hover:ring-claret-200 focus-within:ring-claret-600 group flex w-full flex-row items-center gap-4 rounded-lg border border-zinc-300 px-4 py-2 duration-200 outline-none placeholder:text-black/50 focus-within:ring-2 hover:ring-2"
        >
          <LucideSearch class="size-4 text-black/50 duration-200 group-focus-within:text-black" />

          <input
            type="text"
            :placeholder="$t('general.search')"
            v-model="search"
            class="w-full outline-none placeholder:text-black/50"
          />

          <button
            v-if="search"
            @click="
              search = '';
              fetchProducts();
            "
            class="flex min-w-fit cursor-pointer flex-row items-center-safe justify-center self-end rounded-full font-semibold text-black/50 hover:text-black"
          >
            {{ $t("general.clear") }}
          </button>

          <button
            v-if="search"
            @click="fetchProducts"
            class="text-claret-600/50 hover:text-claret-600 flex min-w-fit cursor-pointer flex-row items-center-safe justify-center self-end rounded-full font-semibold"
          >
            {{ $t("general.search") }}
          </button>
        </label>

        <img
          v-if="profile.hasProfile"
          :src="urlEncodedName"
          class="aspect-square h-10 w-auto rounded-full"
        />
        <a
          href="/login"
          class="bg-claret-600 flex h-full w-fit min-w-fit items-center justify-center rounded-lg px-4 font-semibold text-white"
          v-else
          >{{ $t("general.login") }}</a
        >
      </div>

      <!-- Sidebar -->
      <aside class="bg-claret-100 flex h-full flex-col rounded-2xl p-6">Some stuff here</aside>

      <!-- Content -->
      <div class="flex w-full flex-col gap-4 md:col-span-3">
        <h2 class="text-xl font-bold">{{ $t("search.all_products") }}</h2>
        <p class="opacity-50" v-if="products && total > 0">
          {{ $t("search.products_count", { count: total }) }}
        </p>
        <div
          class="grid size-full grid-cols-1 gap-4 overflow-visible rounded-2xl md:grid-cols-2 lg:grid-cols-3"
        >
          <p
            class="flex size-full items-center justify-center text-lg md:col-span-2 lg:col-span-3"
            v-if="!products || products.length == 0"
          >
            {{ $t("search.no_products") }}
          </p>
          <template v-else v-for="product in products" :key="product.id">
            <ProductCard :product="product" />
          </template>
        </div>

        <!-- Paging section -->
        <div class="flex w-full flex-row items-center justify-between">
          <button
            class="cursor-pointer rounded-lg border border-zinc-300 p-2 hover:border-zinc-500 disabled:cursor-not-allowed disabled:opacity-50"
            @click="page = Math.max(page - 1, 1)"
            :disabled="page == 1"
          >
            <LucideChevronLeft class="size-4 text-black" />
          </button>

          <div class="flex w-fit flex-row items-center justify-center gap-1 font-semibold">
            <template v-if="pageRange.min >= 3">
              <button
                class="size-10 cursor-pointer rounded-lg border border-zinc-300 p-2 hover:border-zinc-500 disabled:cursor-not-allowed disabled:opacity-50"
                @click="page = 1"
              >
                1
              </button>

              <LucideEllipsis class="mx-2 size-4 text-black" />
            </template>

            <template v-for="num in pageRange.max - pageRange.min + 1" :key="num">
              <button
                class="size-10 cursor-pointer rounded-lg border border-zinc-300 p-2 hover:border-zinc-500 disabled:cursor-not-allowed disabled:opacity-50"
                v-if="page != num + pageRange.min - 1"
                @click="page = num + pageRange.min - 1"
              >
                {{ num + pageRange.min - 1 }}
              </button>
              <button
                class="bg-claret-600 size-10 cursor-pointer rounded-lg border border-zinc-300 p-2 text-white hover:border-zinc-500 disabled:cursor-not-allowed disabled:opacity-50"
                @click="page = num + pageRange.min - 1"
                v-else
              >
                {{ num + pageRange.min - 1 }}
              </button>
            </template>

            <template v-if="pageRange.max <= maxPages - 2">
              <LucideEllipsis class="mx-2 size-4 text-black" />

              <button
                class="size-10 cursor-pointer rounded-lg border border-zinc-300 p-2 hover:border-zinc-500 disabled:cursor-not-allowed disabled:opacity-50"
                @click="page = maxPages"
              >
                {{ maxPages }}
              </button>
            </template>
          </div>

          <button
            class="cursor-pointer rounded-lg border border-zinc-300 p-2 hover:border-zinc-500 disabled:cursor-not-allowed disabled:opacity-50"
            @click="page = Math.min(page + 1, maxPages)"
            :disabled="page == maxPages"
          >
            <LucideChevronRight class="size-4 text-black" />
          </button>
        </div>
      </div>
    </div>
  </WhiteContainer>
</template>
