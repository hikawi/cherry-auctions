<script setup lang="ts">
const props = defineProps<{
  disabled: boolean;
  label: string;
}>();
</script>

<template>
  <button
    :disabled="props.disabled"
    class="bg-claret-600 disabled:bg-claret-700 border-claret-600 enabled:hover:text-claret-600 disabled:border-claret-700 w-full cursor-pointer rounded-xl border-2 p-2 py-3 text-white transition-all duration-200 hover:shadow-md enabled:hover:bg-transparent disabled:cursor-progress disabled:opacity-50"
  >
    {{ props.label }}
  </button>
</template>
