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
    <LucideCheckSquare2 class="text-claret-600 size-6 translate-y-px" v-if="model" />
    <LucideSquare class="size-6 translate-y-px text-black" v-else />

    {{ label }}
  </label>
</template>
