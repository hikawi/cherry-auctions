<script setup lang="ts">
import WhiteContainer from "@/components/shared/WhiteContainer.vue";
import { computed, onMounted, ref, watch, watchEffect } from "vue";
import {
  LucideChevronLeft,
  LucideChevronRight,
  LucideEllipsis,
  LucideSearch,
} from "lucide-vue-next";
import { useAuthFetch } from "@/hooks/use-auth-fetch";
import { endpoints } from "@/consts";
import type { Category, Product } from "@/types";
import ProductCard from "@/components/index/ProductCard.vue";
import NavigationBar from "@/components/shared/NavigationBar.vue";

const { authFetch } = useAuthFetch();

const loading = ref(true);
const search = ref<string>();
const products = ref<Product[]>();
const categories = ref<Category[]>();
const sortType = ref<string>("id_asc");
const selectedCategories = ref<number[]>([]);
const page = ref(1);
const total = ref(0);
const maxPages = ref(1);

const pageRange = computed(() => {
  return {
    min: Math.max(page.value - 2, 1),
    max: Math.min(page.value + 2, maxPages.value),
  };
});

watchEffect(() => {
  page.value = Math.min(page.value, maxPages.value);
});

watch(page, fetchProducts);
watch(selectedCategories, fetchProducts);
watch(sortType, fetchProducts);

function toggleCategory(cat: Category) {
  if (selectedCategories.value.includes(cat.id)) {
    selectedCategories.value = selectedCategories.value.filter((id) => id != cat.id);
  } else {
    selectedCategories.value = [...selectedCategories.value, cat.id];
  }
}

async function fetchCategories() {
  try {
    const res = await authFetch(endpoints.categories.get);
    if (res.ok) {
      categories.value = await res.json();
    }
  } catch (e) {
    console.log(e);
  }
}

async function fetchProducts() {
  loading.value = true;

  // First, we build the URL
  const url = new URL(endpoints.products.get);
  url.searchParams.append("page", Math.max(page.value, 1).toString());
  url.searchParams.append("per_page", "12"); // just an arbitrary number because it fits well.
  selectedCategories.value.forEach((id) => url.searchParams.append("category", id.toString()));

  switch (sortType.value) {
    case "id_desc":
      url.searchParams.append("sort", "id");
      url.searchParams.append("asc", "false");
      break;
    case "price_asc":
      url.searchParams.append("sort", "price");
      url.searchParams.append("asc", "true");
      break;
    case "price_desc":
      url.searchParams.append("sort", "price");
      url.searchParams.append("asc", "false");
      break;
    case "time_asc":
      url.searchParams.append("sort", "time");
      url.searchParams.append("asc", "true");
      break;
    case "time_desc":
      url.searchParams.append("sort", "time");
      url.searchParams.append("asc", "false");
      break;
    default:
      url.searchParams.append("sort", "id");
      url.searchParams.append("asc", "true");
      break;
  }

  if (search.value) {
    url.searchParams.append("query", search.value);
  }

  console.log(url.toString());

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
  await Promise.all([fetchProducts(), fetchCategories()]);
});
</script>

<template>
  <WhiteContainer class="justify-start gap-4">
    <NavigationBar />

    <!-- Classic 12-column layout -->
    <div class="grid h-fit w-full grid-cols-1 gap-4 sm:flex-row sm:gap-8 md:grid-cols-4 lg:px-6">
      <!-- NavigationBar -->
      <div class="flex flex-col items-center gap-2 md:col-span-4 md:flex-row">
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

        <!-- <button -->
        <!--   class="bg-claret-200 hover:bg-claret-300 flex w-full shrink-0 cursor-pointer flex-row items-center justify-center gap-2 rounded-lg px-4 py-2 font-semibold duration-200 md:w-fit" -->
        <!--   @click="cycleSort" -->
        <!-- > -->
        <!--   <LucideArrowUpDown class="size-6 text-black" /> -->
        <!---->
        <!--   {{ $t(`search.sort_${sortType}`) }} -->
        <!-- </button> -->
        <select
          class="flex h-full w-full min-w-fit shrink-0 cursor-pointer flex-row items-center gap-2 rounded-lg border border-zinc-300 px-4 py-2 duration-200 hover:bg-zinc-200 md:max-w-fit"
          v-model="sortType"
        >
          <option value="id_asc">{{ $t("search.sort_id_asc") }}</option>
          <option value="id_desc">{{ $t("search.sort_id_desc") }}</option>
          <option value="price_asc">{{ $t("search.sort_price_asc") }}</option>
          <option value="price_desc">{{ $t("search.sort_price_desc") }}</option>
          <option value="time_asc">{{ $t("search.sort_time_asc") }}</option>
          <option value="time_desc">{{ $t("search.sort_time_desc") }}</option>
        </select>
      </div>

      <!-- Sidebar -->
      <aside
        class="flex h-full snap-x snap-mandatory snap-start snap-always flex-row items-center gap-2 overflow-x-scroll rounded-2xl border border-zinc-300 bg-white p-2 shadow-md md:flex-col md:p-6"
      >
        <h2 class="mb-2 hidden w-full text-left text-xl font-bold md:block">
          {{ $t("search.all_categories") }}
        </h2>

        <template v-for="category in categories" :key="category.id">
          <div
            class="bg-claret-50 min-w-fit cursor-pointer rounded-lg px-4 py-2 md:w-full"
            @click="toggleCategory(category)"
            :class="{ 'bg-claret-200': selectedCategories.includes(category.id) }"
          >
            {{ category.name }}
          </div>

          <template v-if="selectedCategories.includes(category.id)">
            <div
              class="bg-claret-50 min-w-fit cursor-pointer rounded-lg px-4 py-2 md:w-full"
              v-for="cat in category.subcategories"
              :key="cat.id"
              @click="toggleCategory(cat)"
              :class="{ 'bg-claret-200': selectedCategories.includes(cat.id) }"
            >
              {{ cat.name }}
            </div>
          </template>
        </template>
      </aside>

      <!-- Content -->
      <div class="flex w-full flex-col gap-4 md:col-span-3">
        <h2 class="text-xl font-bold">{{ $t("search.all_products") }}</h2>
        <p class="opacity-50" v-if="products && total > 0">
          {{ $t("search.products_count", { count: total }) }}
        </p>
        <div
          class="grid size-full grid-cols-1 gap-4 overflow-visible rounded-2xl md:grid-cols-2 xl:grid-cols-3"
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
