<script setup lang="ts">
import { useI18n } from "vue-i18n";

useI18n();
</script>

<template>
  <main
    class="flex items-stretch justify-stretch w-full min-h-screen h-fit bg-claret-600 p-6 font-fredoka"
    :class="{ 'font-mplus': $i18n.locale == 'ja' }"
  >
    <router-view></router-view>
  </main>
</template>
