<script setup lang="ts">
import { endpoints } from "@/consts";
import { useAuthFetch } from "@/hooks/use-auth-fetch";
import { useProfileStore } from "@/stores/profile";
import type { Product } from "@/types";
import { useInterval } from "@vueuse/core";
import { LucideChevronLeft, LucideChevronRight, LucideHeart } from "lucide-vue-next";
import { ref, watch } from "vue";

const profile = useProfileStore();
const { authFetch } = useAuthFetch();
const now = useInterval(3000);

const props = defineProps<{
  data: Product & { similar_products?: Product[]; categories: { id: number; name: string }[] };
}>();

const emits = defineEmits<{
  toggleFavorite: [];
}>();

const currentImage = ref(0);
const images = ref([
  props.data.thumbnail_url,
  ...(props.data.product_images || []).map((img) => img.url),
]);

watch(now, () => {
  currentImage.value = (currentImage.value + 1) % images.value.length;
});

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
      :src="images[currentImage]"
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

    <div class="flex w-full flex-row items-center justify-between">
      <button
        class="min-w-fit cursor-pointer rounded-full p-2 duration-200 hover:bg-zinc-200"
        :disabled="currentImage == 0"
        @click="currentImage = Math.max(currentImage - 1, 0)"
      >
        <LucideChevronLeft />
      </button>

      <div class="flex w-fit flex-row items-center justify-center gap-2">
        <img
          v-for="(img, idx) in images"
          :key="img"
          :src="img"
          class="relative size-12 cursor-pointer rounded-lg object-cover object-center"
          @click="currentImage = idx"
          :class="{
            'ring-claret-600 ring-2': idx == currentImage,
            'hover:ring-claret-600/50 duration-200 hover:ring-2': idx != currentImage,
          }"
        />
      </div>

      <button
        class="min-w-fit cursor-pointer rounded-full p-2 duration-200 hover:bg-zinc-200"
        :disabled="currentImage == images.length - 1"
        @click="currentImage = Math.min(currentImage + 1, images.length - 1)"
      >
        <LucideChevronRight />
      </button>
    </div>
  </div>
</template>
