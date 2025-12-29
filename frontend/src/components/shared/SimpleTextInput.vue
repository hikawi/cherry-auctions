<script setup lang="ts">
import type { InputTypeHTMLAttribute } from "vue";

defineProps<{
  label: string;
  type?: InputTypeHTMLAttribute;
  placeholder?: string;
}>();

const model = defineModel();
</script>

<template>
  <label class="flex w-full flex-col gap-1">
    {{ label }}
    <input
      :type
      required
      v-model="model"
      :placeholder
      class="hover:ring-claret-200 focus:ring-claret-600 w-full rounded-lg border border-zinc-300 px-4 py-2 duration-200 outline-none placeholder:text-black/50 hover:ring-2 focus:ring-2"
    />
  </label>
</template>
