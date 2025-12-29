<script setup lang="ts">
defineProps<{
  label: string;
  choices: {
    label: string;
    value: string;
  }[];
}>();

const model = defineModel();
</script>

<template>
  <label class="flex w-full flex-col gap-1">
    {{ label }}
    <select
      class="hover:ring-claret-200 focus:ring-claret-600 w-full appearance-none rounded-lg border border-zinc-300 px-4 py-2 duration-200 outline-none placeholder:text-black/50 hover:ring-2 focus:ring-2"
      v-model="model"
      required
    >
      <option
        v-for="choice in choices"
        :value="choice.value"
        :key="choice.value"
        :selected="model == choice.value"
      >
        {{ choice.label }}
      </option>
    </select>
  </label>
</template>
