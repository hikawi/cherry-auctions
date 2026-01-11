<script setup lang="ts">
import { ref } from "vue";
import MoneyInput from "../shared/inputs/MoneyInput.vue";
import z from "zod";
import CheckboxInput from "../shared/inputs/CheckboxInput.vue";

const props = defineProps<{
  name: string;
  nextBidValue: number;
  loading: boolean;
}>();

const emits = defineEmits<{
  cancel: [];
  confirm: [cents: number, auto: boolean];
}>();

const nextBid = ref<number>(props.nextBidValue);
const auto = ref<boolean>(false);
const error = ref("");

function confirm() {
  error.value = "";

  const cents = nextBid.value;
  const number = z.number().min(props.nextBidValue).safeParse(cents);

  if (number.error) {
    error.value = "products.minimum_bid";
    return;
  }

  emits("confirm", cents, auto.value);
}
</script>

<template>
  <div class="flex w-full max-w-lg flex-col gap-4 rounded-2xl bg-white p-4 shadow-md sm:p-6">
    <h2 class="text-center text-xl font-bold">{{ $t("products.confirm_bid") }}</h2>
    <hr class="h-px w-full rounded-full border-zinc-300" />

    <MoneyInput
      :label="$t('products.next_bid')"
      :default="$n(nextBidValue / 100, 'decimal')"
      @change="(val) => (nextBid = val)"
    />
    <CheckboxInput :label="$t('products.automated_bid')" v-model="auto" />

    <p
      class="border-watermelon-600 text-watermelon-600 bg-watermelon-200/50 rounded-xl border-2 px-4 py-2"
      v-if="error"
    >
      {{ $t(error, { bid: $n(nextBidValue / 100, "currency") }) }}
    </p>

    <div class="flex w-full flex-row items-center justify-center gap-2 font-semibold">
      <button
        class="cursor-pointer rounded-xl bg-zinc-200 px-4 py-1 duration-200 hover:bg-zinc-300"
        @click="$emit('cancel')"
      >
        {{ $t("general.cancel") }}
      </button>

      <button
        class="bg-claret-600 hover:bg-claret-700 cursor-pointer rounded-xl px-4 py-1 text-white duration-200 disabled:cursor-not-allowed disabled:opacity-50"
        @click="confirm"
        :disabled="loading"
      >
        {{ loading ? $t("general.loading") : $t("general.lets_go") }}
      </button>
    </div>
  </div>
</template>
