<script setup lang="ts">
import { endpoints } from "@/consts";
import { useAuthFetch } from "@/hooks/use-auth-fetch";
import { useProfileStore } from "@/stores/profile";
import type { Product } from "@/types";
import dayjs from "dayjs";
import { computed, ref } from "vue";
import { useI18n } from "vue-i18n";
import AvatarCircle from "../shared/AvatarCircle.vue";
import ErrorDialog from "../shared/ErrorDialog.vue";
import OverlayScreen from "../shared/OverlayScreen.vue";
import ConfirmBidDialog from "./ConfirmBidDialog.vue";
import { LucideThumbsUp } from "lucide-vue-next";
import { useRouter } from "vue-router";

const profile = useProfileStore();
const { locale } = useI18n();
const { authFetch } = useAuthFetch();
const router = useRouter();

const props = defineProps<{
  data: Product & { similar_products?: Product[]; categories: { id: number; name: string }[] };
}>();

const emits = defineEmits<{
  reload: [];
  rate: [productId: number, revieweeId: number];
}>();

const confirmBidDialog = ref(false);
const loading = ref(false);
const error = ref<string>();
const isExpired = computed(() => dayjs(props.data.expired_at).isBefore(dayjs()));
const cantBidReason = computed(() => {
  if (!props.data) {
    return undefined;
  }

  if (!profile.profile) {
    return "products.cant_bid_logged_out";
  }

  if (dayjs(props.data.expired_at).isBefore(dayjs())) {
    return "products.cant_bid_ended";
  }

  if (props.data.seller.id == profile.profile.id) {
    return "products.cant_bid_self";
  }

  if (!props.data.allows_unrated_buyers && profile.profile.average_rating < 0.8) {
    return "products.cant_bid_no_rating";
  }

  if (props.data.current_highest_bid?.bidder.id == profile.profile.id) {
    return "products.cant_bid_outbid_self";
  }

  if (
    props.data.denied_bidders &&
    props.data.denied_bidders.some((bidder) => bidder.email == profile.profile?.email)
  ) {
    return "products.cant_bid_denied";
  }

  return undefined;
});
const expiresAtDisplay = computed(() => {
  return props.data != null
    ? dayjs(props.data.expired_at).locale(locale.value).format("lll")
    : "N/A";
});
const expiresInDisplay = computed(() => {
  return props.data != null && !isExpired.value
    ? dayjs(props.data.expired_at).locale(locale.value).fromNow(false)
    : "N/A";
});
const createdAtDisplay = computed(() => {
  return props.data != null
    ? dayjs(props.data.created_at).locale(locale.value).format("lll")
    : "N/A";
});
const nextBidValue = computed(() => {
  if (!props.data) {
    return 0;
  }

  if (!props.data.current_highest_bid) {
    return props.data.starting_bid;
  }

  const highestBid = props.data.current_highest_bid.price;
  return highestBid + props.data.step_bid_value;
});

// Send a request to add a bid.
async function bid(cents: number, auto: boolean) {
  loading.value = true;
  error.value = "";

  try {
    const url = auto
      ? endpoints.products.autobids(props.data.id)
      : endpoints.products.bids(props.data.id);
    const res = await authFetch(url, {
      method: "POST",
      body: JSON.stringify({ bid: cents }),
    });
    confirmBidDialog.value = false;

    if (res.ok) {
      emits("reload");
    } else {
      error.value = "products.failed_to_bid";
    }
  } finally {
    loading.value = false;
  }
}

// Creates a new conversation to finalize product
async function finalizeProduct() {
  loading.value = true;
  error.value = "";

  try {
    const res = await authFetch(endpoints.chat.index, {
      method: "POST",
      body: JSON.stringify({ product_id: props.data.id }),
    });
    if (res.status == 409 || res.ok) {
      router.push({ name: "messages" });
      return;
    }

    error.value = "products.failed_to_create_chat";
  } finally {
    loading.value = false;
  }
}
</script>

<template>
  <OverlayScreen shown v-if="confirmBidDialog" class="p-6">
    <ConfirmBidDialog
      :name="data.name"
      :nextBidValue
      :loading
      @cancel="confirmBidDialog = false"
      @confirm="bid"
    />

    <div class="flex flex-col"></div>
  </OverlayScreen>

  <OverlayScreen shown v-if="error" class="p-6">
    <ErrorDialog :title="$t('general.error')" :description="$t(error)" @close="error = undefined" />
  </OverlayScreen>

  <div class="flex h-full w-full flex-col justify-between gap-4">
    <div class="flex w-full flex-col gap-4">
      <h1 class="text-2xl font-bold">{{ data.name }}</h1>

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
        <AvatarCircle
          :name="data.seller.name"
          :avatar_url="data.seller.avatar_url"
          class="size-8"
        />
        <span
          >{{ data.seller.name }} ({{
            data.seller.email ? data.seller.email : $t("products.deleted_email")
          }})</span
        >
        <span>
          {{ $t("products.rating", { rating: $n(data.seller.average_rating, "decimal") }) }}
        </span>
        <button
          v-if="
            data.product_state == 'ended' &&
            data.current_highest_bid?.bidder.id == profile.profile?.id
          "
          @click="$emit('rate', data.id, data.seller.id)"
        >
          <LucideThumbsUp class="size-4 cursor-pointer text-zinc-500 hover:text-black" />
        </button>
      </div>

      <div class="flex w-full flex-col gap-1">
        <span class="text-sm">{{ $t("products.created_at", { at: createdAtDisplay }) }}</span>
        <span class="text-sm">{{ $t("products.expires_at", { at: expiresAtDisplay }) }}</span>
        <span class="text-sm">{{ $t("products.expires_in", { in: expiresInDisplay }) }}</span>
      </div>
    </div>

    <div class="grid grid-cols-1 gap-2 md:grid-cols-2">
      <span
        class="text-claret-600 w-full text-center font-semibold md:col-span-2"
        v-if="cantBidReason"
      >
        {{ $t(cantBidReason) }}
      </span>

      <button
        class="bg-claret-600 enabled:hover:text-claret-600 border-claret-600 flex cursor-pointer flex-col items-center justify-center gap-0 rounded-xl border-2 px-4 py-2 text-white duration-200 enabled:hover:bg-white disabled:cursor-not-allowed disabled:opacity-50"
        @click="confirmBidDialog = true"
        :disabled="cantBidReason != undefined"
      >
        <span class="text-lg font-semibold">{{ $t("products.bid") }}</span>
        <span class="text-sm">{{ $n(nextBidValue / 100, "currency") }}</span>
      </button>

      <button
        class="flex cursor-pointer flex-col items-center justify-center gap-0 rounded-xl border-2 border-black bg-black px-4 py-2 text-white duration-200 enabled:hover:bg-white enabled:hover:text-black disabled:cursor-not-allowed disabled:opacity-50"
        :disabled="cantBidReason != undefined"
        v-if="data.bin_price"
      >
        <span class="text-lg font-semibold">{{ $t("products.bin") }}</span>
        <span class="text-sm">{{ $n(data.bin_price / 100, "currency") }}</span>
      </button>

      <button
        class="w-full cursor-pointer rounded-xl border-2 border-emerald-600 bg-emerald-600 p-4 font-semibold text-white duration-200 enabled:hover:bg-white enabled:hover:text-emerald-600 disabled:cursor-not-allowed disabled:opacity-50 md:col-span-2"
        v-if="
          data.product_state == 'ended' &&
          (data.current_highest_bid?.bidder.id == profile.profile?.id ||
            data.seller.id == profile.profile?.id)
        "
        :disabled="data.finalized_at != null"
        @click="finalizeProduct"
      >
        {{ data.finalized_at ? $t("products.already_checkout") : $t("products.checkout") }}
      </button>
    </div>
  </div>
</template>
