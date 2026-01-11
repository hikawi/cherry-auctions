<script setup lang="ts">
import { endpoints } from "@/consts";
import { useAuthFetch } from "@/hooks/use-auth-fetch";
import { ref } from "vue";
import { LucideThumbsUp, LucideThumbsDown } from "lucide-vue-next";
import TextInput from "../shared/inputs/TextInput.vue";

const { authFetch } = useAuthFetch();

const props = defineProps<{
  revieweeId: number;
  productId: number;
}>();

const emit = defineEmits<{
  rate: [rating: number];
  fail: [];
  cancel: [];
}>();

const feedback = ref("");
const rating = ref<0 | 1>(1);
const loading = ref(false);

async function rate() {
  const ratingVal = rating.value;
  loading.value = true;

  const res = await authFetch(endpoints.ratings.index, {
    method: "POST",
    body: JSON.stringify({
      reviewee_id: props.revieweeId,
      product_id: props.productId,
      feedback: feedback.value,
      rating: rating.value,
    }),
  });
  loading.value = false;
  if (res.ok) {
    emit("rate", ratingVal);
  } else {
    emit("fail");
  }
}
</script>

<template>
  <div class="flex w-full max-w-xl flex-col items-center gap-4 rounded-2xl bg-white p-6 shadow-md">
    <h2 class="text-lg font-semibold">{{ $t("products.rate") }}</h2>
    <hr class="w-full rounded-full border-zinc-300" />
    <div class="flex flex-row items-center gap-4">
      <button
        class="cursor-pointer rounded-full border p-2"
        :class="{
          'border-emerald-600 text-emerald-600': rating == 1,
          'border-zinc-300 text-zinc-300 hover:border-zinc-500 hover:text-zinc-500': rating == 0,
        }"
        @click="rating = 1"
      >
        <LucideThumbsUp class="size-8 stroke-1" />
      </button>
      <button
        class="cursor-pointer rounded-full border p-2"
        :class="{
          'border-watermelon-600 text-watermelon-600': rating == 0,
          'border-zinc-300 text-zinc-300 hover:border-zinc-500 hover:text-zinc-500': rating == 1,
        }"
        @click="rating = 0"
      >
        <LucideThumbsDown class="size-8 stroke-1" />
      </button>
    </div>

    <TextInput :label="$t('products.feedback')" type="text" required v-model="feedback" />

    <div class="flex w-full flex-row items-center justify-center gap-2 font-semibold">
      <button
        class="cursor-pointer rounded-xl bg-zinc-200 px-4 py-1 duration-200 hover:bg-zinc-300"
        @click="$emit('cancel')"
      >
        {{ $t("general.cancel") }}
      </button>

      <button
        class="bg-claret-600 hover:bg-claret-700 cursor-pointer rounded-xl px-4 py-1 text-white duration-200 disabled:cursor-not-allowed disabled:opacity-50"
        @click="rate"
        :disabled="loading"
      >
        {{ loading ? $t("general.loading") : $t("products.rate") }}
      </button>
    </div>
  </div>
</template>
