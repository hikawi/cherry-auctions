<script setup lang="ts">
import LanguageSelect from "./components/shared/LanguageSelect.vue";
import LoadingSpinner from "./components/shared/LoadingSpinner.vue";
import WhiteContainer from "./components/shared/WhiteContainer.vue";
import { useProfileStore } from "./stores/profile";

const profile = useProfileStore();
</script>

<template>
  <LanguageSelect />

  <main
    class="bg-claret-600 font-chiron flex h-fit min-h-screen w-full items-stretch justify-stretch p-6"
  >
    <WhiteContainer v-if="profile.isFetching">
      <LoadingSpinner />
    </WhiteContainer>
    <router-view v-else></router-view>
  </main>
</template>
