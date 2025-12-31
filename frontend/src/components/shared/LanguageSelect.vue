<script setup lang="ts">
import { computed, onMounted } from "vue";
import { useI18n } from "vue-i18n";

const { locale } = useI18n({ useScope: "global" });
const emoji = computed(() => {
  switch (locale.value) {
    case "ja-JP":
      return "🇯🇵";
    case "en-US":
      return "🇺🇸";
  }
  return "?";
});

onMounted(() => {
  locale.value = localStorage.getItem("locale") || "en-US";
});

function changeLanguage() {
  locale.value = locale.value === "en-US" ? "ja-JP" : "en-US";
  localStorage.setItem("locale", locale.value);
}
</script>

<template>
  <button
    @click="changeLanguage"
    class="bg-claret-600 fixed right-3 cursor-pointer rounded-lg rounded-tl-none p-3 text-white"
  >
    {{ emoji }}
  </button>
</template>
