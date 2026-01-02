<script setup lang="ts">
import { QuillEditor } from "@vueup/vue-quill";

defineProps<{
  label: string;
  required?: boolean;
}>();

const editorOptions = {
  modules: {
    toolbar: [
      [{ header: [1, 2, false] }],
      ["bold", "italic", "underline"],
      [{ list: "ordered" }, { list: "bullet" }],
      ["link"],
      ["clean"],
    ],
  },
  theme: "snow",
};

const model = defineModel<string>();
</script>

<template>
  <label class="flex w-full flex-col gap-1">
    <span class="relative w-fit">
      {{ label }}

      <span class="text-claret-600 absolute top-0 -right-2.5">*</span>
    </span>

    <QuillEditor
      v-model:content="model"
      content-type="html"
      theme="snow"
      :options="editorOptions"
      class="hover:border-claret-600 h-fit max-w-full border-2 border-transparent duration-200"
    />
  </label>
</template>
