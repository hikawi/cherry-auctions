<script setup lang="ts">
import { useI18n } from "vue-i18n";
import { LucideX } from "lucide-vue-next";
import { computed, ref } from "vue";
import SimpleTextInput from "../shared/SimpleTextInput.vue";
import SimpleSelectInput from "../shared/SimpleSelectInput.vue";
import type { Category } from "@/types";
import z from "zod";

const props = defineProps<{
  categories: Category[];
  category: Category;
}>();

const { t } = useI18n();

const name = ref(props.category.name);
const error = ref<string>();
const choice = ref<string | undefined>(props.category.parent_id?.toString());

const choices = computed(() =>
  props.categories
    .filter((cat) => cat.id != props.category.id)
    .sort((a, b) => a.id - b.id)
    .map((cat) => ({
      label: cat.name,
      value: cat.id.toString(),
    })),
);

const emits = defineEmits<{
  close: [];
  confirm: [id: number, name: string, parentId?: number];
}>();

function emitConfirm() {
  error.value = undefined;
  const schema = z.object({
    name: z.string().min(2).max(200),
    parentId: z.coerce.number().optional(),
  });

  const body = schema.safeParse({ name: name.value, parentId: choice.value });
  if (body.error) {
    error.value = "admin.categories.error_invalid_name";
    return;
  }

  emits("confirm", props.category.id, body.data.name, body.data.parentId);
}
</script>

<template>
  <div class="flex w-full max-w-lg flex-col gap-4 rounded-2xl bg-white p-6 shadow-md">
    <div class="flex w-full flex-row items-center justify-between gap-2">
      <p class="text-lg font-semibold">{{ t("admin.categories.edit_category") }}</p>
      <button @click="$emit('close')" class="cursor-pointer rounded-full p-2 hover:bg-black/20">
        <LucideX class="size-4 stroke-black" />
      </button>
    </div>

    <form class="flex w-full flex-col gap-2">
      <SimpleTextInput
        :label="t('admin.categories.category_name')"
        type="text"
        v-model="name"
        :placeholder="t('admin.categories.category_name_placeholder')"
      />

      <SimpleSelectInput :label="t('admin.categories.parent_category')" :choices v-model="choice" />

      <p
        v-if="error"
        class="bg-claret-100 border-claret-500 text-claret-700 w-full rounded-xl border-2 px-4 py-2"
      >
        {{ t(error) }}
      </p>
    </form>

    <button
      class="bg-claret-600 hover:bg-claret-700 mt-2 flex cursor-pointer flex-row items-center-safe justify-center gap-1 self-end rounded-full px-4 py-2 font-semibold text-white"
      @click="emitConfirm"
    >
      {{ t("general.confirm") }}
    </button>
  </div>
</template>
