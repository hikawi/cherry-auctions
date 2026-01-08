<script setup lang="ts">
import ProductCard from "@/components/index/ProductCard.vue";
import NavigationBar from "@/components/shared/NavigationBar.vue";
import PagingSection from "@/components/shared/PagingSection.vue";
import WhiteContainer from "@/components/shared/WhiteContainer.vue";
import { endpoints } from "@/consts";
import { useAuthFetch } from "@/hooks/use-auth-fetch";
import type { Product } from "@/types";
import { useHead } from "@unhead/vue";
import { LucideRefreshCcw } from "lucide-vue-next";
import { onMounted, ref, watch } from "vue";
import { useI18n } from "vue-i18n";

const { authFetch } = useAuthFetch();
const { t } = useI18n();

useHead({
  title: t("meta.my_bids"),
});

const data = ref<Product[]>();
const page = ref(1);
const maxPages = ref(1);
const loading = ref(false);

watch(maxPages, (val) => (page.value = Math.min(page.value, val)));
watch(page, fetchMyBids);

async function fetchMyBids() {
  const url = new URL(endpoints.users.me.bids);
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
    await fetchMyBids();
  } finally {
    loading.value = false;
  }
});
</script>

<template>
  <WhiteContainer class="justify-start gap-8">
    <NavigationBar />

    <!-- Section for seeing bids -->
    <section class="flex w-full max-w-4xl flex-col gap-4">
      <div class="flex flex-col gap-2 sm:flex-row sm:items-center sm:justify-between">
        <h2 class="text-2xl font-bold">{{ $t("navigation.my_bids") }}</h2>

        <button
          class="bg-claret-600 hover:bg-claret-700 mt-2 flex w-fit cursor-pointer flex-row items-center justify-center gap-2 self-end rounded-full px-4 py-1 font-semibold text-white duration-200"
          @click="fetchMyBids"
        >
          <LucideRefreshCcw class="size-4 text-white" />
          {{ $t("general.refresh") }}
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
        <PagingSection :page :maxPages @setPage="(p) => (page = p)" />
      </div>
      <p class="w-full py-6 text-center text-xl font-semibold" v-else>
        {{ $t("auctions.no_auctions") }}
      </p>
    </section>
  </WhiteContainer>
</template>
