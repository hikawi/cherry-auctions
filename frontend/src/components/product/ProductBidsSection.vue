<script setup lang="ts">
import { endpoints } from "@/consts";
import { useAuthFetch } from "@/hooks/use-auth-fetch";
import { useProfileStore } from "@/stores/profile";
import type { Product } from "@/types";
import dayjs from "dayjs";
import { LucideCircle, LucideCrown } from "lucide-vue-next";
import { computed, ref } from "vue";
import { useI18n } from "vue-i18n";
import AvatarCircle from "../shared/AvatarCircle.vue";
import ErrorDialog from "../shared/ErrorDialog.vue";
import OverlayScreen from "../shared/OverlayScreen.vue";

const { locale } = useI18n();
const { authFetch } = useAuthFetch();
const profile = useProfileStore();

const props = defineProps<{
  data: Product & { similar_products?: Product[]; categories: { id: number; name: string }[] };
}>();

const emit = defineEmits<{
  reload: [];
}>();

const error = ref<string>();
const denyingBidder = ref(false);
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

async function denyBidder(id: number) {
  denyingBidder.value = true;
  error.value = "";
  try {
    const res = await authFetch(endpoints.products.denials(props.data.id), {
      method: "POST",
      body: JSON.stringify({ user_id: id }),
    });
    if (!res.ok) {
      error.value = "products.error_cant_deny";
    } else {
      emit("reload");
    }
  } finally {
    denyingBidder.value = false;
  }
}
</script>

<template>
  <OverlayScreen shown v-if="error">
    <ErrorDialog :title="$t('products.error')" :description="$t(error)" @close="error = ''" />
  </OverlayScreen>

  <section class="flex w-full flex-col gap-4">
    <h2 class="text-xl font-bold">{{ $t("products.bids") }}</h2>

    <div
      v-if="sortedBids && sortedBids.length > 0"
      class="flex w-full flex-col items-center gap-4 rounded-xl border border-zinc-300 p-4 md:w-fit md:flex-row md:p-6"
    >
      <div class="flex flex-row items-center gap-2">
        <LucideCrown class="size-6 fill-amber-600 text-amber-600" />
        <h3 class="text-lg font-semibold">{{ $t("products.current_highest_bidder") }}</h3>
      </div>

      <div class="flex flex-col items-center gap-1 md:ml-8">
        <AvatarCircle
          :name="sortedBids[0].bidder.name"
          :avatar_url="sortedBids[0].bidder.avatar_url"
          class="size-12"
        />
        <span
          >{{ sortedBids[0].bidder.name }} ({{
            $t("products.rating", { rating: sortedBids[0].bidder.average_rating })
          }})</span
        >
        <span class="text-sm">{{ $n(sortedBids[0].price / 100, "currency") }}</span>
        <button
          class="bg-claret-600 hover:bg-claret-700 cursor-pointer rounded-lg px-2 py-1 font-semibold text-white duration-200"
          @click="() => denyBidder(sortedBids![0].bidder.id)"
          v-if="profile.profile && profile.profile.id == data.seller.id"
        >
          {{ $t("products.deny_bidder") }}
        </button>
      </div>
    </div>

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
                price: $n(bid.price / 100, "currency"),
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
