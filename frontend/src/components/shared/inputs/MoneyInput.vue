<script setup lang="ts">
import { ref, watch } from "vue";
import { useI18n } from "vue-i18n";

const { n } = useI18n();

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
  let rawValue = input.value.replace(/[^0-9.]/g, ""); // Allow only digits and one dot

  // Prevent multiple decimal points (e.g., 12.34.56 -> 12.3456)
  const parts = rawValue.split(".");
  if (parts.length > 2) {
    rawValue = parts[0] + "." + parts.slice(1).join("");
  }

  const integerPart = parts[0];
  const decimalPart = parts[1];

  let formatted = "";

  if (integerPart) {
    // Format the integer part with commas
    formatted = n(parseInt(integerPart, 10), "decimal");
  } else if (rawValue.startsWith(".")) {
    // Handle case where user types "." first
    formatted = "0";
  }

  // If there is a dot, append it back manually
  if (rawValue.includes(".")) {
    // Limit to 2 decimal places for currency
    formatted += "." + (decimalPart ? decimalPart.substring(0, 2) : "");
  }

  // Cursor Management
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
