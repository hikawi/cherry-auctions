<script setup lang="ts">
import NavigationBar from "@/components/shared/NavigationBar.vue";
import WhiteContainer from "@/components/shared/WhiteContainer.vue";
import { endpoints } from "@/consts";
import { useAuthFetch } from "@/hooks/use-auth-fetch";
import { LucideRefreshCw } from "lucide-vue-next";
import { onMounted, ref } from "vue";
import type { Product } from "@/types";
import ProductCard from "@/components/index/ProductCard.vue";

const { authFetch } = useAuthFetch();

const loading = ref(true);
const error = ref<string>();
const data = ref<{ ending_soon: Product[]; highest_bids: Product[]; top_bids: Product[] }>();

onMounted(fetchTopProducts);

async function fetchTopProducts() {
  loading.value = true;
  try {
    const res = await authFetch(endpoints.products.top);
    data.value = await res.json();
  } catch {
    error.value = "home.cant_fetch_products";
  } finally {
    loading.value = false;
  }
}
</script>

<template>
  <WhiteContainer class="justify-start gap-8">
    <NavigationBar />

    <div class="flex w-full max-w-4xl flex-col items-center justify-start gap-8">
      <button
        class="bg-claret-600 hover:bg-claret-700 flex w-fit cursor-pointer flex-row items-center justify-center gap-2 self-end rounded-full px-4 py-1 font-semibold text-white duration-200 md:mx-4"
        @click="fetchTopProducts"
      >
        <LucideRefreshCw class="size-4" :class="{ 'animate-spin': loading }" />
        {{ loading ? $t("general.refreshing") : $t("general.refresh") }}
      </button>

      <section class="flex w-full flex-col gap-4">
        <h2 class="text-xl font-semibold">{{ $t("home.top_ending_soons") }}</h2>

        <p
          class="w-full py-8 text-center"
          v-if="!data?.ending_soon || data.ending_soon.length == 0"
        >
          {{ $t("home.no_products") }}
        </p>
        <div class="grid grid-cols-1 gap-4 sm:grid-cols-2 lg:grid-cols-3" v-else>
          <ProductCard v-for="product in data.ending_soon" :key="product.id" :product />
        </div>
      </section>

      <section class="flex w-full flex-col gap-4">
        <h2 class="text-xl font-semibold">{{ $t("home.top_bids") }}</h2>

        <p class="w-full py-8 text-center" v-if="!data?.top_bids || data.top_bids.length == 0">
          {{ $t("home.no_products") }}
        </p>
        <div class="grid grid-cols-1 gap-4 sm:grid-cols-2 lg:grid-cols-3" v-else>
          <ProductCard v-for="product in data.top_bids" :key="product.id" :product />
        </div>
      </section>

      <section class="flex w-full flex-col gap-4">
        <h2 class="text-xl font-semibold">{{ $t("home.top_highest_bids") }}</h2>

        <p
          class="w-full py-8 text-center"
          v-if="!data?.highest_bids || data.highest_bids.length == 0"
        >
          {{ $t("home.no_products") }}
        </p>
        <div class="grid grid-cols-1 gap-4 sm:grid-cols-2 lg:grid-cols-3" v-else>
          <ProductCard v-for="product in data.highest_bids" :key="product.id" :product />
        </div>
      </section>
    </div>
  </WhiteContainer>
</template>
