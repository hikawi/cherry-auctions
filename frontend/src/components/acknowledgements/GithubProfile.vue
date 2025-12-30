<script setup lang="ts">
import { computed } from "vue";

const props = defineProps<{
  username: string;
  name: string;
}>();

const githubAvatarLink = computed(() => {
  return `https://avatars.githubusercontent.com/${props.username}`;
});
const githubLink = computed(() => {
  return `https://github.com/${props.name}`;
});
</script>

<template>
  <div class="flex flex-col items-center justify-center gap-2">
    <img
      :src="githubAvatarLink"
      :alt="name"
      class="border-claret-600 aspect-square size-32 rounded-full border-2 object-cover shadow-md"
    />

    <div class="flex w-full flex-col gap-0 text-center">
      <span class="font-semibold">{{ name }}</span>
      <a class="hover:text-claret-600 text-sm text-black/50 duration-200" :href="githubLink"
        >@{{ username }}</a
      >
    </div>
  </div>
</template>
