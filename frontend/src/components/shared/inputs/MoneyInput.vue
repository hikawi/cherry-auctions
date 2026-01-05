<script setup lang="ts">
import { ref, watch } from "vue";

defineProps<{
  label: string;
  required?: boolean;
  placeholder?: string;
}>();

const emit = defineEmits<{
  change: [val: number];
}>();

const value = ref<string>("");

watch(value, (nc) => {
  const numeric = parseFloat(nc.replace(/,/g, ""));
  const cents = Math.round(numeric * 100);
  emit("change", cents);
});

// Courtesy of AI
function onChange(e: Event) {
  const input = e.target as HTMLInputElement;
  const digits = input.value.replace(/\D/g, "");
  const amount = Number(digits) / 100;
  const formatted = new Intl.NumberFormat("en-US", {
    minimumFractionDigits: 2,
    maximumFractionDigits: 2,
  }).format(amount);

  value.value = formatted;
}
</script>

<template>
  <label class="flex w-full flex-col gap-1">
    <span class="relative w-fit">
      {{ label }}

      <span class="text-claret-600 absolute top-0 -right-2.5" v-if="required">*</span>
    </span>

    <div class="flex w-full flex-row rounded-lg border border-zinc-300">
      <!-- All is in USD, idk how to support other currencies -->
      <div class="w-fit rounded-l-lg border-r border-zinc-300 bg-zinc-200 px-4 py-2">$</div>

      <input
        type="text"
        inputmode="numeric"
        class="w-full px-4 py-2 outline-none"
        :placeholder
        :value
        @input="onChange"
      />
    </div>
  </label>
</template>
