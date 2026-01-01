<script setup lang="ts">
import { computed } from "vue";

const props = defineProps<{
  name?: string;
  avatar_url?: string;
  hover?: boolean;
}>();

const imageSrc = computed(() => {
  return (
    props.avatar_url ||
    `https://ui-avatars.com/api/?background=random&name=${encodeURIComponent(props.name || "")}`
  );
});
</script>

<template>
  <img
    :src="imageSrc"
    :alt="name"
    class="aspect-square h-10 w-auto rounded-full"
    :class="{
      'hover:ring-claret-600 duration-200 hover:ring-2': hover,
    }"
  />
</template>
