<script setup lang="ts">
import { useI18n } from "vue-i18n";
import { LucideX } from "lucide-vue-next";
import type { Category } from "@/types";

const { t } = useI18n();

const props = defineProps<{
  category: Category;
}>();

const emits = defineEmits<{
  close: [];
  confirm: [id: number];
}>();

function emitConfirm() {
  emits("confirm", props.category.id);
}
</script>

<template>
  <div class="flex w-full max-w-lg flex-col gap-4 rounded-2xl bg-white p-6 shadow-md">
    <div class="flex w-full flex-row items-center justify-between gap-2">
      <p class="text-lg font-semibold">{{ t("admin.categories.delete_category") }}</p>
      <button @click="$emit('close')" class="cursor-pointer rounded-full p-2 hover:bg-black/20">
        <LucideX class="size-4 stroke-black" />
      </button>
    </div>

    <p>
      {{ t("admin.categories.delete_description", { id: category.id, name: category.name }) }}
    </p>

    <button
      class="bg-claret-600 hover:bg-claret-700 mt-2 flex cursor-pointer flex-row items-center-safe justify-center gap-1 self-end rounded-full px-4 py-2 font-semibold text-white"
      @click="emitConfirm"
    >
      {{ t("general.confirm") }}
    </button>
  </div>
</template>
