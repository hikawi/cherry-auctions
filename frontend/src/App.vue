<script setup lang="ts">
import LanguageSelect from "./components/shared/LanguageSelect.vue";
</script>

<template>
  <LanguageSelect />
  <main
    class="bg-claret-600 font-chiron flex h-fit min-h-screen w-full items-stretch justify-stretch p-6"
  >
    <router-view></router-view>
  </main>
</template>
