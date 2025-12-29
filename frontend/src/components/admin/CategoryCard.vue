<script setup lang="ts">
import type { Category } from "@/types";
import { LucideChevronDown, LucideChevronUp, LucidePencil, LucideTrash } from "lucide-vue-next";
import { ref } from "vue";
import { useI18n } from "vue-i18n";

const { t } = useI18n();
const props = defineProps<{
  category: Category;
  editCallback: (cat: Category) => void;
  deleteCallback: (cat: Category) => void;
}>();

const open = ref(false);
</script>

<template>
  <div
    class="hover:border-claret-600/50 active:border-claret-600 flex flex-col gap-2 rounded-2xl border-2 border-black/20 px-4 py-2 shadow-md duration-200"
  >
    <div class="flex w-full flex-row items-center justify-between">
      <div class="flex flex-row items-center gap-2">
        <span class="font-semibold">{{ props.category.name }}</span>

        <button
          class="cursor-pointer disabled:cursor-not-allowed disabled:opacity-50"
          @click="() => editCallback(category)"
        >
          <LucidePencil class="size-4 stroke-black" />
        </button>

        <button
          class="cursor-pointer disabled:cursor-not-allowed disabled:opacity-50"
          :disabled="props.category.subcategories.length != 0"
          @click="() => deleteCallback(category)"
        >
          <LucideTrash class="size-4 stroke-black" />
        </button>
      </div>

      <div>
        <LucideChevronUp
          @click="open = !open"
          v-if="open && props.category.subcategories.length > 0"
          class="size-6 stroke-black"
        />
        <LucideChevronDown
          @click="open = !open"
          v-else-if="props.category.subcategories.length > 0"
          class="size-6 stroke-black"
        />
      </div>
    </div>

    <div class="flex flex-col gap-2" v-if="open">
      <p v-if="props.category.subcategories.length == 0">
        {{ t("admin.categories.no_subcategories") }}
      </p>
      <CategoryCard
        v-else
        v-for="cat in [...props.category.subcategories].sort((a, b) => a.id - b.id)"
        :category="cat"
        :key="cat.id"
        :editCallback="() => editCallback(cat)"
        :deleteCallback="() => deleteCallback(cat)"
      />
    </div>
  </div>
</template>
