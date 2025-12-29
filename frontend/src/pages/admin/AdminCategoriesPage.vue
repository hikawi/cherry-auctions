<script setup lang="ts">
import CategoryCard from "@/components/admin/CategoryCard.vue";
import CreateCategoryDialog from "@/components/admin/CreateCategoryDialog.vue";
import DeleteCategoryDialog from "@/components/admin/DeleteCategoryDialog.vue";
import EditCategoryDialog from "@/components/admin/EditCategoryDialog.vue";
import LoadingSpinner from "@/components/shared/LoadingSpinner.vue";
import OverlayScreen from "@/components/shared/OverlayScreen.vue";
import PrimaryButton from "@/components/shared/PrimaryButton.vue";
import { endpoints } from "@/consts";
import { useAuthFetch } from "@/hooks/use-auth-fetch";
import type { Category } from "@/types";
import { LucidePlus } from "lucide-vue-next";
import { computed, onMounted, ref } from "vue";
import { useI18n } from "vue-i18n";

const { authFetch } = useAuthFetch({ json: true });
const { t } = useI18n();

const loading = ref(true);
const categoriesData = ref<Category[]>();
const deletingCategory = ref<Category>();
const editingCategory = ref<Category>();
const dialogShown = ref();

const flatCategories = computed<Category[]>(() => {
  if (!categoriesData.value) {
    return [];
  }

  const arr: Category[] = [];
  function traverse(cats: Category[]) {
    for (const cat of cats) {
      arr.push({ ...cat, subcategories: [] });
      traverse(cat.subcategories);
    }
  }

  traverse(categoriesData.value);
  return arr;
});

onMounted(async () => {
  fetchCategories();
});

async function fetchCategories() {
  loading.value = true;
  try {
    const res = await authFetch(endpoints.categories.get);
    if (res.ok) {
      categoriesData.value = await res.json();
    }
  } finally {
    loading.value = false;
  }
}

async function createCategory(name: string, parentId: number | undefined) {
  loading.value = true;
  try {
    const res = await authFetch(endpoints.categories.post, {
      method: "POST",
      body: JSON.stringify({ name, parent_id: parentId }),
    });
    if (res.status == 201) {
      await fetchCategories();
    }
    dialogShown.value = undefined;
  } finally {
    loading.value = false;
    dialogShown.value = undefined;
  }
}

async function editCategory(id: number, name: string, parentId?: number) {
  loading.value = true;
  try {
    const res = await authFetch(endpoints.categories.edit(id), {
      method: "PUT",
      body: JSON.stringify({ name, parent_id: parentId }),
    });
    if (res.status == 200) {
      await fetchCategories();
    }
    dialogShown.value = undefined;
  } finally {
    loading.value = false;
    dialogShown.value = undefined;
  }
}

async function deleteCategory(id: number) {
  loading.value = true;
  try {
    const res = await authFetch(endpoints.categories.delete(id), {
      method: "DELETE",
    });
    if (res.status == 204) {
      await fetchCategories();
    }
    dialogShown.value = undefined;
  } finally {
    loading.value = false;
    dialogShown.value = undefined;
  }
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
    <DeleteCategoryDialog
      v-if="dialogShown == 'delete' && deletingCategory"
      @close="dialogShown = undefined"
      @confirm="deleteCategory"
      :category="deletingCategory!"
    />
    <EditCategoryDialog
      v-if="dialogShown == 'edit' && editingCategory"
      @close="dialogShown = undefined"
      @confirm="editCategory"
      :category="editingCategory!"
      :categories="flatCategories"
    />
  </OverlayScreen>

  <div v-if="categoriesData" class="flex w-full max-w-4xl flex-col gap-2">
    <button
      class="bg-claret-600 hover:bg-claret-700 flex flex-row items-center-safe justify-center gap-1 self-end rounded-full px-4 py-2 font-semibold text-white"
      @click="dialogShown = 'create'"
    >
      <LucidePlus class="size-4 stroke-white" />
      {{ t("general.create") }}
    </button>

    <CategoryCard
      v-for="category in categoriesData"
      :key="category.id"
      :category
      :editCallback="
        (cat) => {
          editingCategory = cat;
          dialogShown = 'edit';
        }
      "
      :deleteCallback="
        (cat) => {
          deletingCategory = cat;
          dialogShown = 'delete';
        }
      "
    />
  </div>
  <LoadingSpinner v-else-if="loading" />
  <div v-else class="flex flex-col gap-2">
    <p>{{ t("admin.categories.cant_load") }}</p>
    <PrimaryButton :disabled="loading" :label="t('general.try_again')" @click="fetchCategories" />
  </div>
</template>
