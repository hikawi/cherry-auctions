<script setup lang="ts">
import type { Category } from "@/types";
import { LucideChevronDown, LucideChevronUp } from "lucide-vue-next";
import { ref } from "vue";
import { useI18n } from "vue-i18n";

const { t } = useI18n();
const props = defineProps<{
  category: Category;
}>();

const open = ref(false);
</script>

<template>
  <div
    class="hover:border-claret-600/50 active:border-claret-600 flex flex-col gap-2 rounded-2xl border-2 border-black/20 px-4 py-2 shadow-md duration-200"
  >
    <div class="flex w-full flex-row items-center justify-between">
      <span class="font-semibold">{{ props.category.name }}</span>

      <LucideChevronUp @click="open = !open" v-if="open" class="size-6 stroke-black" />
      <LucideChevronDown @click="open = !open" v-else class="size-6 stroke-black" />
    </div>

    <div class="flex flex-col gap-2" v-if="open">
      <p v-if="props.category.subcategories.length == 0">
        {{ t("admin.categories.no_subcategories") }}
      </p>
      <CategoryCard
        v-else
        v-for="cat in props.category.subcategories"
        :category="cat"
        :key="cat.id"
      />
    </div>
  </div>
</template>
