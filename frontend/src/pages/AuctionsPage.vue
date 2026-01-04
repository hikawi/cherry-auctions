<script setup lang="ts">
import ProductCard from "@/components/index/ProductCard.vue";
import CreateAuctionDialog from "@/components/product/CreateAuctionDialog.vue";
import NavigationBar from "@/components/shared/NavigationBar.vue";
import OverlayScreen from "@/components/shared/OverlayScreen.vue";
import WhiteContainer from "@/components/shared/WhiteContainer.vue";
import { endpoints } from "@/consts";
import { useAuthFetch } from "@/hooks/use-auth-fetch";
import type { Category, Product } from "@/types";
import { LucideChevronLeft, LucideChevronRight, LucidePackage } from "lucide-vue-next";
import { onMounted, ref, watch } from "vue";

const { authFetch } = useAuthFetch();

const createDialogShown = ref(false);
const data = ref<Product[]>();
const categories = ref<Category[]>();
const page = ref(1);
const maxPages = ref(1);
const loading = ref(false);

watch(maxPages, (val) => (page.value = Math.min(page.value, val)));
watch(page, fetchMyAuctions);

function onCreate(status: number) {
  createDialogShown.value = false;
  if (status != 201) {
    console.log("uh oh");
  }
}

async function fetchCategories() {
  const res = await authFetch(endpoints.categories.get);
  if (res.ok) {
    categories.value = await res.json();
  }
}

async function fetchMyAuctions() {
  const url = new URL(endpoints.products.me);
  url.searchParams.append("page", page.value.toString());
  url.searchParams.append("per_page", "12");

  try {
    const res = await authFetch(url);
    if (res.ok) {
      const json = await res.json();
      data.value = json.data;
      maxPages.value = Math.max(json.total_pages, 1);
    }
  } catch {
    // TODO: Add error flow
  }
}

onMounted(async () => {
  loading.value = true;
  try {
    await Promise.all([fetchCategories(), fetchMyAuctions()]);
  } finally {
    loading.value = false;
  }
});
</script>

<template>
  <WhiteContainer class="justify-start gap-8">
    <NavigationBar />

    <OverlayScreen :shown="createDialogShown">
      <CreateAuctionDialog @close="createDialogShown = false" @status="onCreate" :categories />
    </OverlayScreen>

    <section class="flex w-full max-w-4xl flex-col gap-4">
      <div class="flex flex-col gap-2 sm:flex-row sm:items-center sm:justify-between">
        <h2 class="text-2xl font-bold">{{ $t("auctions.title") }}</h2>

        <button
          class="bg-claret-600 hover:bg-claret-700 mt-2 flex w-fit cursor-pointer flex-row items-center justify-center gap-2 self-end rounded-full px-4 py-1 font-semibold text-white duration-200"
          @click="createDialogShown = true"
        >
          <LucidePackage class="size-4 text-white" />
          {{ $t("auctions.new") }}
        </button>
      </div>

      <div
        class="grid w-full grid-cols-1 gap-4 sm:grid-cols-2 md:grid-cols-3"
        v-if="data && data.length > 0"
      >
        <template v-for="product in data" :key="product.id">
          <ProductCard :product />
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
    </section>
  </WhiteContainer>
</template>
