<script setup lang="ts">
defineProps<{
  label: string;
  name: string;
  choices: {
    label: string;
    value: string;
  }[];
}>();

const model = defineModel({ required: true });
</script>

<template>
  <fieldset class="flex w-full flex-col gap-2 pt-2">
    <legend class="relative w-fit">
      {{ label }}

      <span class="text-claret-600 absolute top-0 -right-2.5">*</span>
    </legend>

    <label
      v-for="choice in choices"
      :key="choice.value"
      @click="model = choice.value"
      class="hover:bg-claret-50 has-checked:bg-claret-50 flex cursor-pointer flex-row items-center-safe gap-2 rounded-xl border border-black px-4 py-2"
    >
      <input
        type="radio"
        :name
        v-model="model"
        :value="choice.value"
        class="checked:border-claret-600 size-4 appearance-none rounded-full border border-black checked:border-4"
      />

      {{ choice.label }}
    </label>
  </fieldset>
</template>
