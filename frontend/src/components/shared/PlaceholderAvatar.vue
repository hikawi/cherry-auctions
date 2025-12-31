<script setup lang="ts">
import { computed } from "vue";

const props = defineProps<{
  name: string;
}>();

const link = computed(() => {
  return `https://ui-avatars.com/api/?name=${encodeURIComponent(props.name)}`;
});
</script>

<template>
  <img :src="link" class="aspect-square h-10 w-auto rounded-full" />
</template>
