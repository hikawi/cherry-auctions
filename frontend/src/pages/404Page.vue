<script setup lang="ts">
import NavigationBar from "@/components/shared/NavigationBar.vue";
import WhiteContainer from "@/components/shared/WhiteContainer.vue";
import { useI18n } from "vue-i18n";

const { t } = useI18n();
</script>

<template>
  <WhiteContainer class="justify-start gap-4">
    <NavigationBar />

    <div class="flex h-full w-full flex-col items-center gap-4 py-16">
      <div class="flex flex-col gap-2 text-center">
        <h1 class="text-watermelon-600 text-4xl font-bold">{{ t("others.404.title") }}</h1>
        <p class="text-center text-balance">{{ t("others.404.description") }}</p>
      </div>
      <a
        href="/"
        class="bg-claret-600 hover:bg-claret-700 w-fit rounded-full px-4 py-2 font-semibold text-white"
        >{{ t("general.back_home") }}</a
      >
    </div>
  </WhiteContainer>
</template>
