<script setup lang="ts">
import { useI18n } from "vue-i18n";

const { locale } = useI18n({ useScope: "global" });

function changeLanguage() {
  locale.value = locale.value === "en-US" ? "ja-JP" : "en-US";
}
</script>

<template>
  <button
    @click="changeLanguage"
    class="bg-claret-600 fixed top-4 right-4 cursor-pointer rounded-full p-4 text-white"
  >
    {{ locale }}
  </button>
</template>
