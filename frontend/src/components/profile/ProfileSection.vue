<script setup lang="ts">
import { useProfileStore } from "@/stores/profile";

const profile = useProfileStore();
</script>

<template>
  <h2 class="text-2xl font-semibold">{{ $t("profile.profile") }}</h2>
  <p>{{ profile.profile }}</p>
</template>
