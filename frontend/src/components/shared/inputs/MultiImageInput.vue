<script setup lang="ts">
import { LucideImage, LucideX } from "lucide-vue-next";
import { onUnmounted, ref, useTemplateRef } from "vue";

defineProps<{
  label: string;
  required?: boolean;
}>();

const previewURLs = ref<string[]>([]);
const model = defineModel<File[]>({ required: true });

const hiddenInput = useTemplateRef<HTMLInputElement>("hiddenInput");

onUnmounted(() => {
  for (const url of previewURLs.value) {
    URL.revokeObjectURL(url);
  }
  previewURLs.value = [];
});

function onInputChange(e: Event) {
  const el = e.target as HTMLInputElement;
  if (el.files && el.files[0]) {
    // Then redo the model
    model.value = [...model.value, el.files[0]];
    previewURLs.value = [...previewURLs.value, URL.createObjectURL(el.files[0])];
  }

  el.value = "";
}

function remove(idx: number) {
  model.value = [...model.value.slice(0, idx), ...model.value.slice(idx + 1)];
  previewURLs.value = previewURLs.value.filter((val, id) => {
    if (id == idx) {
      URL.revokeObjectURL(val);
    }
    return id != idx;
  });
}
</script>

<template>
  <label class="flex w-full flex-col gap-4">
    <input ref="hiddenInput" class="hidden" type="file" accept="image/*" @change="onInputChange" />

    <span class="relative w-fit">
      {{ label }}

      <span class="text-claret-600 absolute top-0 -right-2.5">*</span>
    </span>

    <div class="flex flex-row flex-wrap gap-2">
      <div v-for="(url, idx) in previewURLs" :key="url" class="relative">
        <img
          :src="url"
          width="128"
          height="128"
          class="size-32 rounded-2xl border border-black object-cover object-center"
        />

        <button
          class="absolute top-0 right-0 translate-x-1/2 -translate-y-1/2 cursor-pointer rounded-full bg-zinc-200 p-1 hover:bg-zinc-300"
          @click="remove(idx)"
        >
          <LucideX class="size-4 text-black" />
        </button>
      </div>

      <button
        class="flex size-32 cursor-pointer items-center justify-center rounded-2xl bg-zinc-200 hover:bg-zinc-300"
        @click="() => hiddenInput?.click()"
      >
        <LucideImage class="size-8 text-black" />
      </button>
    </div>
  </label>
</template>
