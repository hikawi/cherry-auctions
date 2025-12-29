<script setup lang="ts">
import { useProfileStore } from "@/stores/profile";

const profile = useProfileStore();
</script>

<template>
  <div
    class="flex w-full flex-col items-center justify-center self-stretch rounded-2xl bg-white p-6"
  >
    <h1 class="text-2xl font-bold">Hello, world</h1>
    <p v-if="profile.profile">
      {{ profile.profile }}
    </p>
    <p v-if="profile.isAdmin">waow, admin!</p>
    <p v-else>are you logged in?</p>
  </div>
</template>
