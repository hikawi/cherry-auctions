<script setup lang="ts">
import NavigationBar from "@/components/shared/NavigationBar.vue";
import WhiteContainer from "@/components/shared/WhiteContainer.vue";
import { ref } from "vue";

const choice = ref<"profile" | "bids" | "favorites" | "">("profile");
</script>

<template>
  <WhiteContainer class="justify-start gap-8">
    <NavigationBar />

    <section class="flex w-full max-w-4xl flex-col gap-4"></section>
  </WhiteContainer>
</template>
