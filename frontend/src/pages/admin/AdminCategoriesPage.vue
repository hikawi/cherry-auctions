<script setup lang="ts">
import CategoryCard from "@/components/admin/CategoryCard.vue";
import CreateCategoryDialog from "@/components/admin/CreateCategoryDialog.vue";
import LoadingSpinner from "@/components/shared/LoadingSpinner.vue";
import OverlayScreen from "@/components/shared/OverlayScreen.vue";
import PrimaryButton from "@/components/shared/PrimaryButton.vue";
import { endpoints } from "@/consts";
import { useFetch } from "@/hooks/use-fetch";
import type { Category } from "@/types";
import { LucidePlus } from "lucide-vue-next";
import { computed, onMounted, ref } from "vue";
import { useI18n } from "vue-i18n";

const { data, loading, doFetch } = useFetch<Category[]>();
const { t } = useI18n();

const flatCategories = computed<Category[]>(() => {
  if (!data.value) {
    return [];
  }

  const arr: Category[] = [];
  function traverse(cats: Category[]) {
    for (const cat of cats) {
      arr.push({ ...cat, subcategories: [] });
      traverse(cat.subcategories);
    }
  }

  traverse(data.value);
  return arr;
});

onMounted(async () => {
  await doFetch(endpoints.categories.get);
});

const dialogShown = ref();

async function createCategory(name: string, parentId: number | undefined) {
  // TODO: Send a POST Category to backend
  console.log(name, parentId);
}
</script>

<template>
  <h1 class="text-2xl font-bold">{{ t("admin.categories.title") }}</h1>

  <OverlayScreen :shown="dialogShown != undefined">
    <CreateCategoryDialog
      v-if="dialogShown == 'create'"
      @close="dialogShown = undefined"
      @confirm="createCategory"
      :categories="flatCategories"
    />
  </OverlayScreen>

  <LoadingSpinner v-if="loading" />
  <div v-else-if="data" class="flex w-full max-w-4xl flex-col gap-2">
    <button
      class="bg-claret-600 hover:bg-claret-700 flex flex-row items-center-safe justify-center gap-1 self-end rounded-full px-4 py-2 font-semibold text-white"
      @click="dialogShown = 'create'"
    >
      <LucidePlus class="size-4 stroke-white" />
      {{ t("general.create") }}
    </button>

    <CategoryCard v-for="category in data" :key="category.id" :category />
  </div>
  <div v-else class="flex flex-col gap-2">
    <p>{{ t("admin.categories.cant_load") }}</p>
    <PrimaryButton
      :disabled="loading"
      :label="t('general.try_again')"
      @click="doFetch(endpoints.categories.get)"
    />
  </div>
</template>
