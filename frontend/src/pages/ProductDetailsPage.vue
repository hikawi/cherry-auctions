<script setup lang="ts">
import ProductCard from "@/components/index/ProductCard.vue";
import ProductBidsSection from "@/components/product/ProductBidsSection.vue";
import ProductDataCard from "@/components/product/ProductDataCard.vue";
import ProductDescriptionSection from "@/components/product/ProductDescriptionSection.vue";
import ProductImageCard from "@/components/product/ProductImageCard.vue";
import ProductQuestionsSection from "@/components/product/ProductQuestionsSection.vue";
import LoadingSpinner from "@/components/shared/LoadingSpinner.vue";
import NavigationBar from "@/components/shared/NavigationBar.vue";
import WhiteContainer from "@/components/shared/WhiteContainer.vue";
import { endpoints } from "@/consts";
import { useAuthFetch } from "@/hooks/use-auth-fetch";
import type { Product } from "@/types";
import { LucideChevronLeft } from "lucide-vue-next";
import { computed, onMounted, ref, watch } from "vue";
import { useRoute } from "vue-router";

const route = useRoute();
const { authFetch } = useAuthFetch();

const loading = ref(false);
const data = ref<
  Product & { similar_products?: Product[]; categories: { id: number; name: string }[] }
>();
const isExpired = computed(
  () => data.value?.product_state == "expired" || data.value?.product_state == "ended",
);

watch(route, fetchProduct);

async function fetchProduct() {
  loading.value = true;

  const id = route.params.id;
  if (!id) {
    loading.value = false;
    return;
  }

  try {
    const res = await authFetch(endpoints.products.details(id));
    if (res.ok) {
      data.value = await res.json();
    } else {
      data.value = undefined;
    }
  } finally {
    loading.value = false;
  }
}

function toggleSimilarFavorite(id: number) {
  if (!data.value) {
    return;
  }

  data.value = {
    ...data.value,
    similar_products: data.value.similar_products?.map((val) =>
      val.id != id ? val : { ...val, is_favorite: !val.is_favorite },
    ),
  };
}

onMounted(() => {
  fetchProduct();
});
</script>

<template>
  <WhiteContainer class="justify-start gap-8">
    <NavigationBar />

    <LoadingSpinner v-if="loading" />

    <div class="flex h-full w-full max-w-4xl flex-col items-center gap-4 py-16" v-else-if="!data">
      <button
        class="hover:text-claret-600 flex w-fit cursor-pointer items-center gap-2 self-start font-semibold duration-200"
        @click="$router.back()"
      >
        <LucideChevronLeft class="size-4" />
        {{ $t("general.back") }}
      </button>

      <div class="flex flex-col gap-2 text-center">
        <h1 class="text-watermelon-400 text-4xl font-bold">{{ $t("others.404.title") }}</h1>
        <p class="text-center text-balance">{{ $t("others.404.description") }}</p>
      </div>
      <router-link
        :to="{ name: 'home' }"
        class="bg-claret-600 hover:bg-claret-700 w-fit rounded-full px-4 py-2 font-semibold text-white"
        >{{ $t("general.back_home") }}</router-link
      >
    </div>

    <div class="flex w-full max-w-4xl flex-col gap-8" v-else>
      <button
        class="hover:text-claret-600 flex w-fit cursor-pointer items-center gap-2 self-start font-semibold duration-200"
        @click="$router.back()"
      >
        <LucideChevronLeft class="size-4" />
        {{ $t("general.back") }}
      </button>

      <div
        v-if="isExpired"
        class="rounded-xl border-2 border-amber-600 bg-amber-200/50 px-4 py-2 font-semibold text-amber-600"
      >
        {{ $t("products.expired") }}
      </div>

      <!-- The product -->
      <div class="grid w-full max-w-4xl grid-cols-1 gap-8 lg:grid-cols-2">
        <!-- The product card -->
        <ProductImageCard
          :data
          @toggleFavorite="data = { ...data, is_favorite: !data.is_favorite }"
        />

        <!-- The data card -->
        <ProductDataCard :data @reload="fetchProduct" />
      </div>

      <!-- Product Description -->
      <ProductDescriptionSection :data @reload="fetchProduct" />

      <!-- Bids Section -->
      <ProductBidsSection :data @reload="fetchProduct" />

      <!-- Questions Section -->
      <ProductQuestionsSection :data @refresh="fetchProduct" />

      <!-- Similar products Section -->
      <section class="flex w-full flex-col gap-4" v-if="data.similar_products">
        <h2 class="text-xl font-bold">{{ $t("products.similar_products") }}</h2>

        <div class="grid grid-cols-1 gap-4 sm:grid-cols-2 lg:grid-cols-3">
          <ProductCard
            v-for="product in data.similar_products"
            :key="product.id"
            :product
            @favoriteToggle="() => toggleSimilarFavorite(product.id)"
          />
        </div>
      </section>
    </div>
  </WhiteContainer>
</template>
