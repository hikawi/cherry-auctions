<script setup lang="ts">
import { LucideCheckSquare2, LucideSquare } from "lucide-vue-next";

defineProps<{
  label: string;
}>();

const model = defineModel<boolean>({ required: true });
</script>

<template>
  <label
    class="flex cursor-pointer flex-row items-center-safe gap-2 rounded-xl"
    @click="model = !model"
  >
    <LucideSquare class="size-5 translate-y-0.5 text-black" v-if="model" />
    <LucideCheckSquare2 class="text-claret-600 size-5 translate-y-0.5" v-else />

    {{ label }}
  </label>
</template>
