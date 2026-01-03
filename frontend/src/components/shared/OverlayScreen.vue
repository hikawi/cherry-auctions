<script setup lang="ts">
defineProps<{
  shown?: boolean;
}>();
</script>

<template>
  <div
    class="fixed inset-0 z-100 flex size-full items-center justify-center bg-black/50 duration-200"
    :class="{ 'opacity-0': !shown, 'opacity-100': shown }"
    v-if="shown"
  >
    <slot />
  </div>
</template>
