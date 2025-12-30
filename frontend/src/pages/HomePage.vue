<script setup lang="ts">
import NavigationBar from "@/components/shared/NavigationBar.vue";
import WhiteContainer from "@/components/shared/WhiteContainer.vue";
import { useProfileStore } from "@/stores/profile";

const profile = useProfileStore();
</script>

<template>
  <WhiteContainer class="justify-start gap-4">
    <NavigationBar />

    <h1 class="text-2xl font-bold">Hello, world</h1>
    <p v-if="profile.profile">
      {{ profile.profile }}
    </p>
    <p v-if="profile.isAdmin">waow, admin!</p>
    <p v-else>are you logged in?</p>
  </WhiteContainer>
</template>
