<script setup lang="ts">
import type { Product } from "@/types";
import dayjs from "dayjs";
import { computed } from "vue";
import { useI18n } from "vue-i18n";

const { locale } = useI18n();

const props = defineProps<{
  data: Product & { similar_products?: Product[]; categories: { id: number; name: string }[] };
}>();

const sortedBids = computed(() => {
  if (!props.data) {
    return undefined;
  }
  return [...props.data.bids].sort((a, b) => b.price - a.price);
});

// Yes, I know I should truncate on the backend side.
function truncate(s: string): string {
  const prefix = s.slice(0, Math.min(s.length / 2 + 1, 4));
  return prefix.padEnd(8, "*");
}

function createAbsoluteTime(time: string): string {
  return dayjs(time).locale(locale.value).format("lll");
}
</script>

<template>
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
                name: truncate(bid.bidder.name || ""),
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
</template>
