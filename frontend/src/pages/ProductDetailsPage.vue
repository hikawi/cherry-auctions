<script setup lang="ts">
import ProductCard from "@/components/index/ProductCard.vue";
import LoadingSpinner from "@/components/shared/LoadingSpinner.vue";
import NavigationBar from "@/components/shared/NavigationBar.vue";
import PlaceholderAvatar from "@/components/shared/PlaceholderAvatar.vue";
import WhiteContainer from "@/components/shared/WhiteContainer.vue";
import { endpoints } from "@/consts";
import { useAuthFetch } from "@/hooks/use-auth-fetch";
import type { Product } from "@/types";
import dayjs from "dayjs";
import { LucideChevronLeft, LucideCircle, LucideReply } from "lucide-vue-next";
import { computed, onMounted, ref } from "vue";
import { useI18n } from "vue-i18n";
import { useRoute } from "vue-router";

const route = useRoute();
const { locale } = useI18n();
const { authFetch } = useAuthFetch();

const loading = ref(false);
const data = ref<
  Product & { similar_products?: Product[]; categories: { id: number; name: string }[] }
>();
const expiresAtDisplay = computed(() => {
  return data.value != null
    ? dayjs(data.value.expired_at).locale(locale.value).format("lll")
    : "N/A";
});
const createdAtDisplay = computed(() => {
  return data.value != null
    ? dayjs(data.value.created_at).locale(locale.value).format("lll")
    : "N/A";
});
const sortedBids = computed(() => {
  if (!data.value) {
    return undefined;
  }

  return [...data.value.bids].sort((a, b) => b.price - a.price);
});

function createAbsoluteTime(time: string): string {
  return dayjs(time).locale(locale.value).format("lll");
}

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

// Yes, I know I should truncate on the backend side.
function truncate(s: string): string {
  const prefix = s.slice(0, Math.min(s.length / 2 + 1, 4));
  return prefix.padEnd(8, "*");
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
      <a
        href="/"
        class="bg-claret-600 hover:bg-claret-700 w-fit rounded-full px-4 py-2 font-semibold text-white"
        >{{ $t("general.back_home") }}</a
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

      <!-- The product -->
      <div class="grid w-full max-w-4xl grid-cols-1 gap-8 lg:grid-cols-2">
        <!-- The product card -->
        <div class="flex w-full flex-col items-center justify-start gap-4">
          <img
            :src="data.thumbnail_url"
            :alt="data.name"
            width="400"
            height="400"
            class="aspect-square h-auto w-full rounded-xl object-cover"
          />
        </div>

        <!-- The data card -->
        <div class="flex h-full w-full flex-col justify-between gap-4">
          <div class="flex w-full flex-col gap-4">
            <div class="flex w-full flex-col gap-1">
              <h1 class="text-2xl font-bold">{{ data.name }}</h1>
              <p>{{ data.description }}</p>
            </div>

            <div
              v-if="data.categories && data.categories.length > 0"
              class="flex w-full flex-row flex-wrap items-center gap-2 text-sm"
            >
              <div
                v-for="category in data.categories"
                :key="category.id"
                class="bg-claret-100 rounded-lg px-2 py-1"
              >
                {{ category.name }}
              </div>
            </div>

            <div class="flex w-full flex-row items-center gap-2">
              <PlaceholderAvatar :name="data.seller.name" class="size-8" />
              <span class="text-lg"
                >{{ data.seller.name }} ({{
                  data.seller.email ? data.seller.email : $t("products.deleted_email")
                }})</span
              >
            </div>

            <div class="flex w-full flex-col gap-1">
              <span class="text-sm">{{ $t("products.created_at", { at: createdAtDisplay }) }}</span>
              <span class="text-sm">{{ $t("products.expires_at", { at: expiresAtDisplay }) }}</span>
            </div>
          </div>

          <div class="grid grid-cols-1 gap-2 md:grid-cols-2">
            <button
              class="bg-claret-600 hover:text-claret-600 border-claret-600 flex cursor-pointer flex-col items-center justify-center gap-0 rounded-xl border-2 px-4 py-2 text-white duration-200 hover:bg-white"
            >
              <span class="text-lg font-semibold">{{ $t("products.bid") }}</span>
              <span class="text-sm"
                >${{ data.current_highest_bid?.price || data.starting_bid }}</span
              >
            </button>

            <button
              class="flex cursor-pointer flex-col items-center justify-center gap-0 rounded-xl border-2 border-black bg-black px-4 py-2 text-white duration-200 hover:bg-white hover:text-black"
              v-if="data.bin_price"
            >
              <span class="text-lg font-semibold">{{ $t("products.bin") }}</span>
              <span class="text-sm">${{ data.bin_price }}</span>
            </button>
          </div>
        </div>
      </div>

      <!-- Bids Section -->
      <section class="flex w-full flex-col gap-4">
        <h2 class="text-xl font-bold">{{ $t("products.bids") }}</h2>

        <div class="flex w-full flex-col gap-2" v-if="data.bids && data.bids.length > 0">
          <template v-for="(bid, idx) in sortedBids" :key="idx">
            <div class="flex w-full flex-row items-center justify-center gap-2">
              <LucideCircle
                class="fill-claret-600 relative size-4 translate-y-0.5 text-transparent"
                :class="{
                  'fill-claret-600': idx == 0,
                  'fill-watermelon-200': idx != 0,
                }"
              />

              <span class="w-full">
                {{
                  $t("products.bid_list", {
                    name: truncate(bid.bidder.name),
                    price: $n(bid.price, "currency"),
                    at: createAbsoluteTime(bid.created_at),
                  })
                }}
              </span>
            </div>
          </template>
        </div>
        <p v-else class="w-full px-4 py-6 text-center">
          {{ $t("products.no_bids") }}
        </p>
      </section>

      <!-- Questions Section -->
      <section class="flex w-full flex-col gap-4">
        <h2 class="text-xl font-bold">{{ $t("products.questions") }}</h2>

        <div class="flex w-full flex-col gap-2" v-if="data.questions && data.questions.length > 0">
          <div
            v-for="(question, idx) in data.questions"
            :key="idx"
            class="flex w-full flex-col gap-2 rounded-2xl border border-zinc-500 p-4"
          >
            <span class="flex flex-row items-center gap-2 text-sm font-semibold">
              <PlaceholderAvatar :name="question.user.name" class="size-6!" />
              {{ $t("products.asked_by", { name: question.user.name }) }}
            </span>

            <p>{{ question.content }}</p>

            <div v-if="question.answer" class="flex flex-col gap-1 px-4">
              <span class="flex flex-row items-center gap-2 text-sm text-black/50">
                <LucideReply class="size-4" />

                {{ $t("products.replied_by_seller") }}
              </span>

              <p>{{ question.answer }}</p>
            </div>
          </div>
        </div>
        <p v-else class="w-full px-4 py-6 text-center">
          {{ $t("products.no_questions") }}
        </p>
      </section>

      <!-- Similar products Section -->
      <section class="flex w-full flex-col gap-4" v-if="data.similar_products">
        <h2 class="text-xl font-bold">{{ $t("products.similar_products") }}</h2>

        <div class="grid grid-cols-1 gap-4 sm:grid-cols-2 lg:grid-cols-3">
          <ProductCard v-for="product in data.similar_products" :key="product.id" :product />
        </div>
      </section>
    </div>
  </WhiteContainer>
</template>
