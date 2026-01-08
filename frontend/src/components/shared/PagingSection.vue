<script setup lang="ts">
import { LucideChevronLeft, LucideChevronRight } from "lucide-vue-next";

defineProps<{
  page: number;
  maxPages: number;
}>();

defineEmits<{
  setPage: [page: number];
}>();
</script>

<template>
  <!-- Paging section -->
  <div class="flex w-full flex-row items-center justify-between sm:col-span-2 md:col-span-3">
    <button
      class="cursor-pointer rounded-lg border border-zinc-300 p-2 hover:border-zinc-500 disabled:cursor-not-allowed disabled:opacity-50"
      @click="$emit('setPage', Math.max(page - 1, 1))"
      :disabled="page == 1"
    >
      <LucideChevronLeft class="size-4 text-black" />
    </button>

    <span>{{ page }} / {{ maxPages }}</span>

    <button
      class="cursor-pointer rounded-lg border border-zinc-300 p-2 hover:border-zinc-500 disabled:cursor-not-allowed disabled:opacity-50"
      @click="$emit('setPage', Math.min(page + 1, maxPages))"
      :disabled="page == maxPages"
    >
      <LucideChevronRight class="size-4 text-black" />
    </button>
  </div>
</template>
