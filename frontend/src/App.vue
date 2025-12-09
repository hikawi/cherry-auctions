<script setup lang="ts">
import { useI18n } from "vue-i18n";
import LanguageSelect from "./components/shared/LanguageSelect.vue";

const { locale } = useI18n({ useScope: "global" });
</script>

<template>
  <LanguageSelect />
  <main
    class="bg-claret-600 font-fredoka flex h-fit min-h-screen w-full items-stretch justify-stretch p-6"
    :class="{ 'font-mplus': locale == 'ja-JP' }"
  >
    <router-view></router-view>
  </main>
</template>
