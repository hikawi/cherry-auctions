<script setup lang="ts">
import ProfileSection from "@/components/profile/ProfileSection.vue";
import NavigationBar from "@/components/shared/NavigationBar.vue";
import WhiteContainer from "@/components/shared/WhiteContainer.vue";
</script>

<template>
  <WhiteContainer class="justify-start gap-8">
    <NavigationBar />

    <section class="flex w-full max-w-4xl flex-col gap-8">
      <ProfileSection />
    </section>
  </WhiteContainer>
</template>
