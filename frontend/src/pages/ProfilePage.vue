<script setup lang="ts">
import ChangePasswordSection from "@/components/profile/ChangePasswordSection.vue";
import FavoritesSection from "@/components/profile/FavoritesSection.vue";
import MyRatingsSection from "@/components/profile/MyRatingsSection.vue";
import ProfileSection from "@/components/profile/ProfileSection.vue";
import NavigationBar from "@/components/shared/NavigationBar.vue";
import WhiteContainer from "@/components/shared/WhiteContainer.vue";
import { useHead } from "@unhead/vue";
import { useI18n } from "vue-i18n";

const { t } = useI18n();

useHead({
  title: t("meta.profile"),
});
</script>

<template>
  <WhiteContainer class="justify-start gap-8">
    <NavigationBar />

    <section class="flex w-full max-w-4xl flex-col gap-8">
      <ProfileSection />
    </section>

    <section class="flex w-full max-w-4xl flex-col gap-8">
      <ChangePasswordSection />
    </section>

    <section class="flex w-full max-w-4xl flex-col gap-8">
      <MyRatingsSection />
    </section>

    <section class="flex w-full max-w-4xl flex-col gap-8">
      <FavoritesSection />
    </section>
  </WhiteContainer>
</template>
