<script setup lang="ts">
import { LucideCircleX } from "lucide-vue-next";

defineProps<{
  title: string;
  description: string;
}>();

defineEmits<{
  close: [];
}>();
</script>

<template>
  <div
    class="flex w-full max-w-sm flex-col items-center justify-center gap-4 rounded-2xl bg-white p-6 shadow-md"
  >
    <LucideCircleX class="text-claret-600 size-8" />

    <div class="flex flex-col gap-2 text-center text-balance">
      <h2 class="text-xl font-semibold">{{ title }}</h2>
      <p>{{ description }}</p>
    </div>

    <button
      class="bg-claret-600 hover:bg-claret-700 rounded-xl px-4 py-1 font-semibold text-white duration-200"
      @click="$emit('close')"
    >
      {{ $t("general.ok") }}
    </button>
  </div>
</template>
