<script setup lang="ts">
import { endpoints } from "@/consts";
import { useAuthFetch } from "@/hooks/use-auth-fetch";
import { useProfileStore } from "@/stores/profile";
import type { Product } from "@/types";
import { LucideHeart } from "lucide-vue-next";

const profile = useProfileStore();
const { authFetch } = useAuthFetch();

const props = defineProps<{
  data: Product & { similar_products?: Product[]; categories: { id: number; name: string }[] };
}>();

const emits = defineEmits<{
  toggleFavorite: [];
}>();

async function toggleFavorite() {
  if (!props.data) {
    return;
  }

  const res = await authFetch(endpoints.products.favorite, {
    method: "POST",
    body: JSON.stringify({ id: props.data.id }),
  });
  if (res.ok) {
    emits("toggleFavorite");
  }
}
</script>

<template>
  <div class="relative flex w-full flex-col items-center justify-start gap-4">
    <img
      :src="data.thumbnail_url"
      :alt="data.name"
      width="400"
      height="400"
      class="aspect-square h-auto w-full rounded-xl object-cover"
    />

    <button
      class="absolute top-0 right-0 z-10 translate-x-1/2 -translate-y-1/2 cursor-pointer rounded-full border border-zinc-300 bg-white p-2 shadow-md hover:border-zinc-500"
      v-if="profile.hasProfile"
      @click="toggleFavorite"
    >
      <LucideHeart
        class="duration-200"
        :class="{
          'text-claret-600 hover:fill-claret-600/50': !data.is_favorite,
          'fill-claret-600 text-claret-600': data.is_favorite,
        }"
      />
    </button>
  </div>
</template>
