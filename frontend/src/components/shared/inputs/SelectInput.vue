<script setup lang="ts">
import { LucideX } from "lucide-vue-next";
import { computed, ref, useTemplateRef } from "vue";

const props = defineProps<{
  label: string;
  options?: {
    value: string;
    label: string;
  }[];
  required?: boolean;
}>();

const emits = defineEmits<{
  add: [value: string];
  remove: [value: string];
}>();

const typed = ref<string>("");
const selected = ref<
  {
    value: string;
    label: string;
  }[]
>([]);
const shownValues = computed(() => {
  return props.options
    ?.filter((opt) => opt.label.toLowerCase().startsWith(typed.value.toLowerCase()))
    .filter((opt) => !selected.value.map((val) => val.value).includes(opt.value));
});
const input = useTemplateRef<HTMLInputElement>("textInput");
const select = useTemplateRef<HTMLSelectElement>("selectOption");

function removeOption(val: string) {
  selected.value = selected.value.filter((opt) => opt.value != val);
  emits("remove", val);
}

function addOption(val: { value: string; label: string }) {
  selected.value = [...selected.value, val];
  emits("add", val.value);
  if (select.value) {
    select.value.value = "";
  }
}
</script>

<template>
  <label class="flex w-full flex-col gap-1" @click="() => input?.focus()">
    <span class="relative w-fit">
      {{ label }}

      <span class="text-claret-600 absolute top-0 -right-2.5" v-if="required">*</span>
    </span>

    <div class="items-ce flex flex-row flex-wrap gap-2 rounded-lg text-sm">
      <div
        class="bg-claret-100 relative flex flex-row items-center gap-1 rounded-full px-2 py-1"
        v-for="opt in selected"
        :key="opt.value"
      >
        {{ opt.label }}

        <button class="cursor-pointer rounded-full" @click="() => removeOption(opt.value)">
          <LucideX class="size-4 translate-y-px text-black" />
        </button>
      </div>
    </div>

    <select
      class="mt-1 cursor-pointer rounded-lg border border-zinc-300 bg-white px-4 py-2 hover:bg-zinc-200"
      ref="selectOption"
      :disabled="!shownValues || shownValues.length == 0"
    >
      <option selected value="" disabled>{{ $t("auctions.choose_category") }}</option>
      <option
        v-for="val in shownValues"
        :value="val.value"
        :key="val.value"
        @click="() => addOption(val)"
      >
        {{ val.label }}
      </option>
    </select>
  </label>
</template>
