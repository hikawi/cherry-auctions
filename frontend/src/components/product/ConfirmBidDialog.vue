<script setup lang="ts">
defineProps<{
  name: string;
  nextBidValue: number;
  loading: boolean;
}>();

defineEmits<{
  cancel: [];
  confirm: [];
}>();
</script>

<template>
  <div class="flex w-full max-w-lg flex-col gap-4 rounded-2xl bg-white p-4 shadow-md sm:p-6">
    <h2 class="text-center text-xl font-bold">{{ $t("products.confirm_bid") }}</h2>
    <hr class="h-px w-full rounded-full border-zinc-300" />
    <p class="text-center text-balance">
      {{
        $t("products.confirm_bid_description", {
          name: name,
          price: $n(nextBidValue / 100, "currency"),
        })
      }}
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
        @click="$emit('confirm')"
        :disabled="loading"
      >
        {{ loading ? $t("general.loading") : $t("general.lets_go") }}
      </button>
    </div>
  </div>
</template>
